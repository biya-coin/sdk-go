package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	injectivetypes "github.com/InjectiveLabs/sdk-go/chain/types"
	chainclient "github.com/InjectiveLabs/sdk-go/client/chain"
)

// sendToMempool 通过 Unix socket 发送交易到 monad-node 的 mempool
// 使用 cosmos-txpool-feed 工具确保 RLP 编码兼容性
func sendToMempool(txBytes []byte, socketPath string) error {
	// 保存交易到临时文件
	tmpFile := "/tmp/go_tx_tosend.raw"
	if err := os.WriteFile(tmpFile, txBytes, 0644); err != nil {
		return fmt.Errorf("保存交易文件失败: %w", err)
	}

	// 调用 cosmos-txpool-feed 工具发送交易
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
	}
	if balance == "" {
		balance = "0"
	}
	return
}

func main() {
	// 连接 gRPC 服务查询账户信息
	grpcConn, err := grpc.Dial("localhost:19900", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("连接 gRPC 失败: %v\n", err)
		return
	}
	defer grpcConn.Close()

	// 使用私钥初始化密钥环
	privateKeyHex := "88CBEAD91AEE890D27BF06E003ADE3D4E952427E88F88D31D61D3EF5E5D54305"
	senderAddress, cosmosKeyring, err := chainclient.InitCosmosKeyring("", "", "", "", "", privateKeyHex, false)
	if err != nil {
		fmt.Printf("初始化密钥环失败: %v\n", err)
		return
	}
	bybAddress, _ := sdktypes.Bech32ifyAddressBytes("byb", senderAddress.Bytes())
	toAddress := "byb1l3c5u98xyk6te8pwcvk36ltf8c3f80m3ddzqhn"
	clientCtx, err := chainclient.NewClientContext("biyachain-1", senderAddress.String(), cosmosKeyring)
	if err != nil {
		fmt.Printf("创建 client context 失败: %v\n", err)
		return
	}

	fmt.Printf("发送方地址: %s\n", bybAddress)
	fmt.Printf("接收方地址: %s\n", toAddress)

	// 查询发送前的余额
	senderAccNum, senderSeq, senderBal := queryAccountInfo(clientCtx, grpcConn, bybAddress)
	_, _, targetBal := queryAccountInfo(clientCtx, grpcConn, toAddress)
	fmt.Printf("发送前 - 发送方余额: %s byb, 接收方余额: %s byb\n", senderBal, targetBal)
	fmt.Printf("账户信息: accountNumber=%d, sequence=%d\n", senderAccNum, senderSeq)

	// 构建转账消息
	msg := &banktypes.MsgSend{FromAddress: bybAddress, ToAddress: toAddress, Amount: []sdktypes.Coin{{Denom: "byb", Amount: math.NewInt(10000)}}}
	fmt.Printf("转账金额: %s\n", msg.Amount)

	// 构建、签名并编码交易
	txFactory := tx.Factory{}.WithChainID("biyachain-1").WithKeybase(cosmosKeyring).WithTxConfig(clientCtx.TxConfig).WithAccountNumber(senderAccNum).WithSequence(senderSeq).WithGas(75000000).WithFees("100000000000000000byb")
	txBuilder, err := txFactory.BuildUnsignedTx(msg)
	if err != nil {
		fmt.Printf("构建交易失败: %v\n", err)
		return
	}
	err = tx.Sign(context.Background(), txFactory, clientCtx.GetFromName(), txBuilder, true)
	if err != nil {
		fmt.Printf("签名交易失败: %v\n", err)
		return
	}
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		fmt.Printf("编码交易失败: %v\n", err)
		return
	}
	fmt.Printf("交易大小: %d bytes\n", len(txBytes))

	// 发送交易到 mempool
	socketPath := "/home/cyyu/monad-workspace/monad-bft/.monad/monad-a/mempool.sock"
	fmt.Printf("发送交易到: %s\n", socketPath)

	// 检查 socket 文件是否存在
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		fmt.Printf("❌ 错误: socket 文件不存在: %s\n", socketPath)
		fmt.Printf("请检查容器是否正在运行: docker ps | grep monad-a\n")
		return
	}

	if err := sendToMempool(txBytes, socketPath); err != nil {
		fmt.Printf("❌ 发送交易失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 交易已发送到 mempool\n")

	// 等待交易被确认
	fmt.Printf("等待交易确认...\n")
	time.Sleep(3 * time.Second)

	// 查询发送后的余额
	_, _, senderBalAfter := queryAccountInfo(clientCtx, grpcConn, bybAddress)
	_, _, targetBalAfter := queryAccountInfo(clientCtx, grpcConn, toAddress)
	fmt.Printf("发送后 - 发送方余额: %s byb, 接收方余额: %s byb\n", senderBalAfter, targetBalAfter)
}
