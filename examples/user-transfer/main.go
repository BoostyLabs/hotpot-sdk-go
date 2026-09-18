package main

import (
	"context"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"

	"github.com/BoostyLabs/hotpot-sdk-go/client"
	"github.com/BoostyLabs/hotpot-sdk-go/examples"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// solanaConfig for user-wallet interactions.
type solanaConfig struct {
	RpcUrl     string `env:"SOLANA_RPC_URL,unset" envDefault:"https://api.devnet.solana.com"`
	PrivateKey string `env:"SOLANA_PRIVATE_KEY,required,unset"`
}

func main() {
	ctx := context.Background()

	cfg := examples.LoadConfig()
	apiClient := cfg.InitClient()

	var solanaCfg solanaConfig
	if err := env.Parse(&solanaCfg); err != nil {
		log.Fatalf("failed to parse solana env: %v", err)
	}

	userKey, err := parsePrivateKey(solanaCfg.PrivateKey)
	if err != nil {
		log.Fatalf("failed to parse SOLANA_PRIVATE_KEY: %v", err)
	}
	userAddress := userKey.PublicKey().String()

	slippageBps, err := types.NewIntFromPercent(2.0)
	if err != nil {
		log.Fatalf("failed to parse slippage: %v", err)
	}

	quote, err := apiClient.GetTheBestQuote(ctx, client.GetTheBestQuoteRequest{
		SourceChain: 3,
		SourceToken: "144vaJc1qmqq2fnceeUG2KqmhxVm2zXkuKJGAtyG1UdE",
		DestChain:   1,
		DestToken:   "0x1c7d4b196cb0c7b01d743fbc6116a902379c7238",
		Amount:      2.,
		Slippage:    slippageBps,
	})
	if err != nil {
		log.Fatalf("failed to get a quote: %v", err)
	}

	log.Printf("Quote %s: %s lots in", quote.ID, quote.SourceAmountLots)

	intentResp, err := apiClient.CreateIntent(ctx, client.CreateIntentRequest{
		QuoteID:                quote.ID,
		UserSourceAddress:      userAddress,
		UserDestinationAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		RefundAddress:          userAddress,
	})
	if err != nil {
		log.Fatalf("failed to create intent: %v", err)
	}
	if intentResp.ApprovalMechanism != types.ApprovalToSignTypeTransfer {
		log.Fatalf("expected the %q approval mechanism, got %q",
			types.ApprovalToSignTypeTransfer, intentResp.ApprovalMechanism)
	}

	log.Printf("Intent %s awaits %s lots at %s",
		intentResp.ID, quote.SourceAmountLots, intentResp.Transfer.Address)

	txHash, err := transfer(ctx, solanaCfg.RpcUrl, userKey,
		quote.SourceToken, intentResp.Transfer.Address, quote.SourceAmountLots.Uint64())
	if err != nil {
		log.Fatalf("failed to transfer to the protocol vault: %v", err)
	}

	log.Printf("Transferred in %s", txHash)

	err = apiClient.SubmitDeposit(ctx, client.SubmitDepositParams{
		IntentID: intentResp.ID,
		TxHash:   txHash,
	})
	if err != nil {
		log.Fatalf("failed to submit deposit: %v", err)
	}

	log.Print("Deposit submitted, waiting for the resolver to fulfill")

	status, err := awaitStatus(ctx, apiClient, intentResp.ID)
	if err != nil {
		log.Fatalf("failed to await a terminal status: %v", err)
	}

	log.Printf("Intent %s finished as %s", intentResp.ID, status)
}

func parsePrivateKey(privateKeyHex string) (solanago.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	privateKey := solanago.PrivateKey(privateKeyBytes)

	return privateKey, privateKey.Validate()
}

// transfer sends the exact quoted amount of an SPL token to the address the intent returned, returns the signature.
func transfer(
	ctx context.Context,
	rpcUrl string,
	from solanago.PrivateKey,
	mintAddress, vaultAddress string,
	amount uint64,
) (string, error) {
	mint, err := solanago.PublicKeyFromBase58(mintAddress)
	if err != nil {
		return "", err
	}

	vault, err := solanago.PublicKeyFromBase58(vaultAddress)
	if err != nil {
		return "", err
	}

	source, _, err := solanago.FindAssociatedTokenAddress(from.PublicKey(), mint)
	if err != nil {
		return "", err
	}

	rpcClient := rpc.New(rpcUrl)

	blockhash, err := rpcClient.GetLatestBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return "", err
	}

	tx, err := solanago.NewTransaction(
		[]solanago.Instruction{
			token.NewTransferInstruction(amount, source, vault, from.PublicKey(), nil).Build(),
		},
		blockhash.Value.Blockhash,
		solanago.TransactionPayer(from.PublicKey()),
	)
	if err != nil {
		return "", err
	}

	_, err = tx.Sign(func(key solanago.PublicKey) *solanago.PrivateKey {
		if key.Equals(from.PublicKey()) {
			return &from
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	// Preflight defaults to simulating against finalized state, which does not know a blockhash
	// this recent yet and rejects the transaction as BlockhashNotFound. It has to be told to
	// simulate at the same commitment the blockhash was read at.
	signature, err := rpcClient.SendTransactionWithOpts(ctx, tx, rpc.TransactionOpts{
		PreflightCommitment: rpc.CommitmentConfirmed,
	})
	if err != nil {
		return "", err
	}

	return signature.String(), awaitConfirmation(ctx, rpcClient, signature)
}

// awaitConfirmation polls until the cluster confirms the transaction or reports it failed.
func awaitConfirmation(ctx context.Context, rpcClient *rpc.Client, signature solanago.Signature) error {
	for range 60 {
		statuses, err := rpcClient.GetSignatureStatuses(ctx, true, signature)
		if err != nil {
			return err
		}

		if len(statuses.Value) > 0 && statuses.Value[0] != nil {
			status := statuses.Value[0]

			if status.Err != nil {
				return errors.New("transaction failed on chain")
			}

			if status.ConfirmationStatus == rpc.ConfirmationStatusConfirmed ||
				status.ConfirmationStatus == rpc.ConfirmationStatusFinalized {
				return nil
			}
		}

		time.Sleep(time.Second)
	}

	return errors.New("transaction was not confirmed in time")
}

// awaitStatus polls the intent until it reaches a status it cannot move on from.
func awaitStatus(ctx context.Context, apiClient *client.Client, intentID uuid.UUID) (types.CombinedStatus, error) {
	for range 120 {
		resp, err := apiClient.GetIntentStatus(ctx, intentID)
		if err != nil {
			return "", err
		}

		switch resp.Status {
		case types.CombinedStatusFulfilled, types.CombinedStatusDeclined,
			types.CombinedStatusExpired, types.CombinedStatusLapsed, types.CombinedStatusRefunded:
			return resp.Status, nil
		}

		log.Printf("status: %s", resp.Status)
		time.Sleep(5 * time.Second)
	}

	return "", errors.New("intent did not reach a terminal status in time")
}
