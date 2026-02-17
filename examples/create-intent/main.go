package main

import (
	"context"
	"log"

	"github.com/BoostyLabs/hotpot-sdk-go/client"
	"github.com/BoostyLabs/hotpot-sdk-go/examples"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

func main() {
	ctx := context.Background()

	apiClient := examples.LoadConfig().InitClient()

	slippageBps, err := types.NewIntFromPercent(2.0)
	if err != nil {
		log.Fatalf("failed to parse slippage: %v", err)
	}

	quote, err := apiClient.GetTheBestQuote(ctx, client.GetTheBestQuoteRequest{
		SourceChain: 1,
		SourceToken: "0xdac17f958d2ee523a2206206994597c13d831ec7", // usdt.
		DestChain:   1,
		DestToken:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", // weth.
		Amount:      2.,
		Slippage:    slippageBps,
		SwapType:    types.SwapTypeStandard,
		DepositType: types.DepositTypeEscrowed,
		// RetailUserID: "optional retail user id, now not set",
	})
	if err != nil {
		log.Fatalf("failed to get a quote: %v", err)
	}

	log.Printf("quote: %+v", quote)

	intentResp, err := apiClient.CreateIntent(ctx, client.CreateIntentRequest{
		QuoteID:                quote.ID,
		UserSourceAddress:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		UserDestinationAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		RefundAddress:          "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		// UserSourcePublicKey: "Required for Bitcoin/Solana source chains, now not set",
	})
	if err != nil {
		log.Fatalf("failed to create intent: %v", err)
	}

	log.Printf("Intent created: %+v", intentResp)
}
