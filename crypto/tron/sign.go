package tron

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/client/transaction"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"

	"github.com/BoostyLabs/hotpot-sdk-go/crypto/evm"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// SignerSender is used for signing and broadcasting Tron transactions.
type SignerSender struct {
	privateKey *ecdsa.PrivateKey
	// Tron hex address corresponding to the private key.
	fromAddress string
	client      *client.GrpcClient
}

// NewSignerSender creates a new SignerSender instance. It parses the key and sets up an rpc client.
func NewSignerSender(privateKeyHex string, clientURI string, clientTimeoutSecs uint) (*SignerSender, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	ethAddress := crypto.PubkeyToAddress(*publicKey)
	fromAddress := "41" + hex.EncodeToString(ethAddress.Bytes())

	clientInstance := client.NewGrpcClientWithTimeout(
		clientURI,
		time.Duration(clientTimeoutSecs)*time.Second,
	)

	if err := clientInstance.Start(); err != nil {
		return nil, err
	}

	return &SignerSender{privateKey, fromAddress, clientInstance}, nil
}

// SignAndSend builds a transaction, signs and broadcasts it.
func (s *SignerSender) SignAndSend(
	contractAddress string,
	dataHex string,
	feeLimit int64,
	callValue int64,
) (*core.Transaction, string, error) {
	data, err := hex.DecodeString(strings.TrimPrefix(dataHex, "0x"))
	if err != nil {
		return nil, "", fmt.Errorf("invalid calldata: %w", err)
	}

	txExt, err := s.client.TriggerContractWithData(
		s.fromAddress,
		contractAddress,
		data,
		feeLimit,
		callValue,
		"",
		0,
	)
	if err != nil {
		return nil, "", fmt.Errorf("build contract call: %w", err)
	}

	if txExt == nil || txExt.Transaction == nil {
		return nil, "", fmt.Errorf("empty transaction returned by TRON node")
	}

	// Sign the transaction.
	signedTx, err := transaction.SignTransactionECDSA(
		txExt.Transaction,
		s.privateKey,
	)
	if err != nil {
		return nil, "", fmt.Errorf("sign transaction: %w", err)
	}

	result, err := s.client.Broadcast(signedTx)
	if err != nil {
		return nil, "", fmt.Errorf("broadcast transaction: %w", err)
	}

	if result == nil || !result.Result {
		var message string
		if result != nil {
			message = string(result.Message)
		}
		return nil, "", fmt.Errorf(
			"transaction broadcast failed: %s",
			message,
		)
	}

	txHash, err := getTransactionHash(signedTx)
	if err != nil {
		return nil, "", err
	}

	return signedTx, txHash, nil
}

// getTransactionHash calculates the hash of a Tron transaction.
func getTransactionHash(tx *core.Transaction) (string, error) {
	// Serialize raw_data using protobuf
	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return "", err
	}

	// Calculate SHA-256
	hash := sha256.Sum256(rawData)

	return hex.EncodeToString(hash[:]), nil
}

// SignPermit2 signs permit2 approval with for provided signer, returning the signature in hex encoding with the '0x' prefix.
func SignPermit2(signer *evm.Signer, quote *types.Quote, permit2Data *types.ApprovalToSignPermit2, deadline int64) (string, error) {
	typedData, err := BuildTypedData(quote, permit2Data, deadline)
	if err != nil {
		return "", err
	}

	sig, err := signer.SignPermit2(typedData)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("0x%x", sig), nil
}
