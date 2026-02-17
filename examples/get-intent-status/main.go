package main

import (
	"context"
	"log"

	"github.com/BoostyLabs/hotpot-sdk-go/examples"
)

func main() {
	ctx := context.Background()

	cfg := examples.LoadConfig()
	apiClient := cfg.InitClient()

	// INFO: Intent id must be specified via env variable::
	//  - INTENT_ID: export INTENT_ID=83c0966e-1111-2222-3333-9a15d04e722c

	log.Printf("Getting intent status id: %v", cfg.IntentID)

	status, err := apiClient.GetIntentStatus(ctx, cfg.IntentID)
	if err != nil {
		log.Fatalf("failed to get intent status: %v", err)
	}

	log.Printf("Intent status: %+v", status)
}
