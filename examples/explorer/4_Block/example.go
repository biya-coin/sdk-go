package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/biya-sdk/sdk-go/client/common"
	explorerclient "github.com/biya-sdk/sdk-go/client/explorer"
)

func main() {
	network := common.LoadNetwork("testnet", "lb")
	explorerClient, err := explorerclient.NewExplorerClient(network)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	blockHeight := "5825046"
	res, err := explorerClient.GetBlock(ctx, blockHeight)
	if err != nil {
		fmt.Println(err)
	}

	str, _ := json.MarshalIndent(res, "", "\t")
	fmt.Print(string(str))
}
