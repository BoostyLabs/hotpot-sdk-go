package main

import (
	"context"
	"log"

	"github.com/BoostyLabs/hotpot-sdk-go/client"
	"github.com/BoostyLabs/hotpot-sdk-go/examples"
)

func main() {
	ctx := context.Background()

	cfg := examples.LoadConfig()
	apiClient := cfg.InitClient()

	// INFO: To specify filters, set of env variables:
	//  - LIMIT
	//  - OFFSET
	//  - TOKEN_QUERY
	//  - NETWORK_ID

	tokens, err := apiClient.ListTokens(ctx, client.ListTokenParams{
		Limit:     cfg.Limit,
		Offset:    cfg.Offset,
		Query:     cfg.TokenQuery,
		NetworkID: cfg.NetworkID,
	})
	if err != nil {
		log.Fatalf("failed to list tokens: %v", err)
	}

	log.Printf("Page metadata %+v", tokens.Metadata)

	for _, token := range tokens.Data {
		log.Printf(
			"%s (%s) on network %d - Contract: %s",
			token.Name, token.Symbol, token.NetworkID, token.ContractAddress,
		)
	}
}
