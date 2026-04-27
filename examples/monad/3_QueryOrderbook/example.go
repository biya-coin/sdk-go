package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	exchangev2types "github.com/InjectiveLabs/sdk-go/chain/exchange/types/v2"
)

func main() {
	// 连接 gRPC 服务
	grpcConn, _ := grpc.Dial("localhost:19900", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer grpcConn.Close()

	exchangeClient := exchangev2types.NewQueryClient(grpcConn)
	marketId := "0xb322bce686ec25364be50728812e33741da1d82e9c91c2c89b91b91d26b0e9c5" // BYB/USDT

	// 查询市场 L3 订单簿（包含完整订单详情）
	req := &exchangev2types.QueryFullSpotOrderbookRequest{
		MarketId: marketId,
	}
	res, err := exchangeClient.L3SpotOrderBook(context.Background(), req)
	if err != nil {
		fmt.Printf("query error: %v\n", err)
		return
	}

	fmt.Printf("market: BYB/USDT\n")
	fmt.Printf("buy orders: %d\n", len(res.Bids))
	fmt.Printf("sell orders: %d\n\n", len(res.Asks))

	if len(res.Bids) > 0 {
		fmt.Println("=== Buy Orders ===")
		for i, order := range res.Bids {
			fmt.Printf("%d. price=%s, quantity=%s, hash=%s\n", 
				i+1, order.Price.String(), order.Quantity.String(), order.OrderHash)
		}
		fmt.Println()
	}

	if len(res.Asks) > 0 {
		fmt.Println("=== Sell Orders ===")
		for i, order := range res.Asks {
			fmt.Printf("%d. price=%s, quantity=%s, hash=%s\n", 
				i+1, order.Price.String(), order.Quantity.String(), order.OrderHash)
		}
	}
}
