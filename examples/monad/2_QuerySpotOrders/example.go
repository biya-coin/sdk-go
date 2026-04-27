package main

import (
	"context"
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	exchangev2types "github.com/InjectiveLabs/sdk-go/chain/exchange/types/v2"
	chainclient "github.com/InjectiveLabs/sdk-go/client/chain"
)

func main() {
	// 连接 gRPC 服务
	grpcConn, _ := grpc.Dial("localhost:19900", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer grpcConn.Close()

	// 使用私钥初始化
	privateKeyHex := "88CBEAD91AEE890D27BF06E003ADE3D4E952427E88F88D31D61D3EF5E5D54305"
	senderAddress, _, _ := chainclient.InitCosmosKeyring("", "", "", "", "", privateKeyHex, false)
	bybAddress, _ := sdktypes.Bech32ifyAddressBytes("byb", senderAddress.Bytes())

	// 查询账户的现货订单
	exchangeClient := exchangev2types.NewQueryClient(grpcConn)
	marketId := "0xb322bce686ec25364be50728812e33741da1d82e9c91c2c89b91b91d26b0e9c5" // BYB/USDT

	// 查询账户订单
	req := &exchangev2types.QueryAccountAddressSpotOrdersRequest{
		MarketId:       marketId,
		AccountAddress: bybAddress,
	}
	res, err := exchangeClient.AccountAddressSpotOrders(context.Background(), req)
	if err != nil {
		fmt.Printf("query error: %v\n", err)
		return
	}

	fmt.Printf("account: %s\n", bybAddress)
	fmt.Printf("market: %s\n", marketId)
	fmt.Printf("total orders: %d\n", len(res.Orders))

	for i, order := range res.Orders {
		orderType := "SELL"
		if order.IsBuy {
			orderType = "BUY"
		}
		fmt.Printf("\norder %d:\n", i+1)
		fmt.Printf("  cid: %s\n", order.Cid)
		fmt.Printf("  type: %s\n", orderType)
		fmt.Printf("  price: %s\n", order.Price.String())
		fmt.Printf("  quantity: %s\n", order.Quantity.String())
		fmt.Printf("  fillable: %s\n", order.Fillable.String())
	}
}
