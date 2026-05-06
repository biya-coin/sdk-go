package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	exchangev2types "github.com/InjectiveLabs/sdk-go/chain/exchange/types/v2"
	injectivetypes "github.com/InjectiveLabs/sdk-go/chain/types"
	chainclient "github.com/InjectiveLabs/sdk-go/client/chain"
)

const (
	grpcEndpoint = "localhost:19900"
	socketPath   = "/home/cyyu/monad-workspace/monad-bft/.monad/monad-a/mempool.sock"
	chainID      = "biyachain-1"
	privateKey   = "88CBEAD91AEE890D27BF06E003ADE3D4E952427E88F88D31D61D3EF5E5D54305"
	marketID     = "0xb322bce686ec25364be50728812e33741da1d82e9c91c2c89b91b91d26b0e9c5" // BYB/USDT

)

// sendToMempool 通过 Unix socket 发送交易到 monad-node 的 mempool
func sendToMempool(txBytes []byte) error {
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return fmt.Errorf("socket 文件不存在: %s，请检查容器是否运行", socketPath)
	}
	tmpFile := "/tmp/go_tx_tosend.raw"
	if err := os.WriteFile(tmpFile, txBytes, 0644); err != nil {
		return fmt.Errorf("保存交易文件失败: %w", err)
	}
	cmd := exec.Command("cosmos-txpool-feed", socketPath, tmpFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("发送失败: %w, 输出: %s", err, string(output))
	}
	return nil
}

// queryAccountInfo 查询账户信息：账户编号、序列号和余额
func queryAccountInfo(clientCtx client.Context, grpcConn *grpc.ClientConn, address string) (accountNumber uint64, sequence uint64, balance string) {
	authClient := authtypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)

	authResp, err := authClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: address})
	if err == nil && authResp.Account != nil {
		var account injectivetypes.EthAccount
		if err := clientCtx.Codec.Unmarshal(authResp.Account.Value, &account); err == nil {
			accountNumber = account.BaseAccount.AccountNumber
			sequence = account.BaseAccount.Sequence
		}
	}

	balResp, _ := bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: address, Denom: "byb"})
	if balResp != nil && balResp.Balance != nil {
		balance = balResp.Balance.Amount.String()
	} else {
		balance = "0"
	}
	return
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("用法: go run example.go <方向> <数量> <最差可接受价格>")
		fmt.Println("  方向: buy 或 sell")
		fmt.Println("  最差可接受价格: 买单填价格上限，卖单填价格下限（滑点保护）")
		fmt.Println("示例: go run example.go buy 1.5 0.002   （买入1.5个，最高接受0.002）")
		fmt.Println("示例: go run example.go sell 1.5 0.0008 （卖出1.5个，最低接受0.0008）")
		return
	}

	side := os.Args[1]
	if side != "buy" && side != "sell" {
		fmt.Println("方向必须为 buy 或 sell")
		return
	}

	quantity, err := sdkmath.LegacyNewDecFromStr(os.Args[2])
	if err != nil {
		fmt.Printf("数量格式错误: %v\n", err)
		return
	}

	worstPrice, err := sdkmath.LegacyNewDecFromStr(os.Args[3])
	if err != nil {
		fmt.Printf("最差价格格式错误: %v\n", err)
		return
	}
	if worstPrice.IsNegative() || worstPrice.IsZero() {
		fmt.Println("最差可接受价格必须大于 0")
		return
	}

	grpcConn, err := grpc.Dial(grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("连接 gRPC 失败: %v\n", err)
		return
	}
	defer grpcConn.Close()

	senderAddress, cosmosKeyring, err := chainclient.InitCosmosKeyring("", "", "", "", "", privateKey, false)
	if err != nil {
		fmt.Printf("初始化密钥环失败: %v\n", err)
		return
	}
	bybAddress, _ := sdktypes.Bech32ifyAddressBytes("byb", senderAddress.Bytes())
	clientCtx, err := chainclient.NewClientContext(chainID, senderAddress.String(), cosmosKeyring)
	if err != nil {
		fmt.Printf("创建 client context 失败: %v\n", err)
		return
	}

	accNum, seq, bal := queryAccountInfo(clientCtx, grpcConn, bybAddress)
	fmt.Printf("发送方地址: %s\n", bybAddress)
	fmt.Printf("账户信息: accountNumber=%d, sequence=%d, balance=%s byb\n", accNum, seq, bal)

	orderId := fmt.Sprintf("market_%s_%d", side, time.Now().UnixNano())

	// worstPrice 为最差可接受价格（滑点保护）：
	//   买单 → 成交价不超过此值
	//   卖单 → 成交价不低于此值
	var orderType exchangev2types.OrderType
	if side == "buy" {
		orderType = exchangev2types.OrderType_BUY
	} else {
		orderType = exchangev2types.OrderType_SELL
	}

	fmt.Printf("\n下市价%s单: market=%s..., worstPrice=%s, quantity=%s, cid=%s\n",
		side, marketID[:10], worstPrice, quantity, orderId)

	msg := &exchangev2types.MsgCreateSpotMarketOrder{
		Sender: bybAddress,
		Order: exchangev2types.SpotOrder{
			MarketId:  marketID,
			OrderType: orderType,
			OrderInfo: exchangev2types.OrderInfo{
				SubaccountId: "",
				FeeRecipient: bybAddress,
				Price:        worstPrice,
				Quantity:     quantity,
				Cid:          orderId,
			},
		},
	}

	txFactory := tx.Factory{}.
		WithChainID(chainID).
		WithKeybase(cosmosKeyring).
		WithTxConfig(clientCtx.TxConfig).
		WithAccountNumber(accNum).
		WithSequence(seq).
		WithGas(75000000).
		WithFees("32000000000000byb")

	txBuilder, err := txFactory.BuildUnsignedTx(msg)
	if err != nil {
		fmt.Printf("构建交易失败: %v\n", err)
		return
	}
	if err := tx.Sign(context.Background(), txFactory, clientCtx.GetFromName(), txBuilder, true); err != nil {
		fmt.Printf("签名失败: %v\n", err)
		return
	}
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		fmt.Printf("编码失败: %v\n", err)
		return
	}
	fmt.Printf("交易大小: %d bytes\n", len(txBytes))

	if err := sendToMempool(txBytes); err != nil {
		fmt.Printf("❌ 发送失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 市价%s单已发送到 mempool\n", side)

	// 等待上链
	fmt.Println("等待交易确认...")
	time.Sleep(3 * time.Second)

	// 查询挂单（市价单成交后不会出现在挂单列表中）
	exchangeClient := exchangev2types.NewQueryClient(grpcConn)
	res, err := exchangeClient.AccountAddressSpotOrders(context.Background(), &exchangev2types.QueryAccountAddressSpotOrdersRequest{
		MarketId:       marketID,
		AccountAddress: bybAddress,
	})
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	fmt.Printf("\n查询结果: 共 %d 个挂单（市价单成交后不在列表中）\n", len(res.Orders))
	found := false
	for i, order := range res.Orders {
		orderSide := "SELL"
		if order.IsBuy {
			orderSide = "BUY"
		}
		mark := " "
		if order.Cid == orderId {
			mark = "★"
			found = true
		}
		fmt.Printf("%s[%d] %s price=%s quantity=%s fillable=%s cid=%s\n",
			mark, i+1, orderSide, order.Price, order.Quantity, order.Fillable, order.Cid)
	}
	if !found {
		fmt.Printf("★ 市价%s单 %s 未在挂单列表中 → 已立即成交\n", side, orderId)
	}
}
