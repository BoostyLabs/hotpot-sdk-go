package evm

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// Signer is the abstraction over the EVM private key for permit2 approval signing.
type Signer struct{ pk *ecdsa.PrivateKey }

// NewSigner creates a new Signer instance from a private key hex string.
func NewSigner(privateKeyHex string) (*Signer, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	return &Signer{privateKey}, nil
}

// SignPermit2 signs permit2 approval with a provided private key, returning the signature.
func (s *Signer) SignPermit2(typesData apitypes.TypedData) ([]byte, error) {
	digest, err := Permit2Hash(typesData)
	if err != nil {
		return nil, err
	}

	return GetPermit2Signature(digest, func(digest []byte) ([]byte, error) { return crypto.Sign(digest, s.pk) })
}

// SignPermit2 signs permit2 approval with for provided signer, returning the signature in hex encoding with the '0x' prefix.
func SignPermit2(signer *Signer, quote *types.Quote, permit2Data *types.ApprovalToSignPermit2, deadline int64) (string, error) {
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
