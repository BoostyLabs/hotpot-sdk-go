package tron

import (
	"fmt"

	"github.com/BoostyLabs/hotpot-sdk-go/crypto/evm"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

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
