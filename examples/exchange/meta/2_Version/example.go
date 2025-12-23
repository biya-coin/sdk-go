package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/biya-sdk/sdk-go/client/common"
	exchangeclient "github.com/biya-sdk/sdk-go/client/exchange"
	metaPB "github.com/biya-sdk/sdk-go/exchange/meta_rpc/pb"
)

func main() {
	// network := common.LoadNetwork("mainnet", "k8s")
	network := common.LoadNetwork("testnet", "lb")
	exchangeClient, err := exchangeclient.NewExchangeClient(network)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	req := metaPB.VersionRequest{}

	res, err := exchangeClient.GetVersion(ctx, &req)
	if err != nil {
		fmt.Println(err)
	}

	str, _ := json.MarshalIndent(res, "", "\t")
	fmt.Print(string(str))
}
