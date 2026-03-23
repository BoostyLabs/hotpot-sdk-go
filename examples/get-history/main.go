package main

import (
	"context"
	"log"
	"time"

	"github.com/BoostyLabs/hotpot-sdk-go/client"
	"github.com/BoostyLabs/hotpot-sdk-go/examples"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

const OneDayInMs int64 = 86_400_000

func main() {
	ctx := context.Background()

	cfg := examples.LoadConfig()
	apiClient := cfg.InitClient()

	// INFO: To specify wallet addresses and retailer ID for swap history, set one of env variables:
	//  - WALLET_ADDRESSES: export WALLET_ADDRESSES=0x1234567890123456789012345678901234567890,bc1pg02klrmyzfkeftcn4j3v2dyly5xh9mpcf5dunxhjst25w7ayu9uq6t2ja0
	//  - RETAIL_ID: export RETAIL_ID=1234567890
	//
	// NOTE: Only one of these variables should be set.
	//
	// LIMIT, OFFSET, NETWORK_ID, TOKEN can be set via env variables.

	log.Printf("Getting swap history by addresses: %v", cfg.WalletAddresses)
	log.Printf("or by retail id: %s", cfg.RetailID)

	history, err := apiClient.ListSwapHistory(ctx, client.ListSwapHistoryParams{
		Wallets:  cfg.WalletAddresses,
		RetailID: cfg.RetailID,
		Limit:    cfg.Limit,
		Offset:   cfg.Offset,
		Status:   types.UiStatusSubmitted,
		Network:  cfg.NetworkID,
		From:     time.Now().UnixMilli() - OneDayInMs,
		To:       time.Now().UnixMilli() + OneDayInMs,
		Token:    cfg.Token,
	})
	if err != nil {
		log.Fatalf("failed to get swap history: %v", err)
	}

	log.Printf("Swap history: %+v", history)
}
