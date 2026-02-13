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

	// INFO: To specify token-filter, set the env variable:
	//  - TOKEN: export TOKEN=USDC

	networks, err := apiClient.ListNetworks(ctx, client.ListNetworkParams{Token: cfg.Token})
	if err != nil {
		log.Fatalf("failed to list networks: %v", err)
	}

	log.Printf("Found %d supported networks", len(networks))
	log.Printf("Networks: %+v", networks)
}
