package tron

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/BoostyLabs/hotpot-sdk-go/crypto/evm"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// Tip712Mask is a constant representing a TIP-712 mask for chain id.
const Tip712Mask uint64 = 0xffffffff

// BuildTypedData constructs a TypedData object from the provided permit2Data with extra values.
func BuildTypedData(quote *types.Quote, permit2Data *types.ApprovalToSignPermit2, deadline int64) (apitypes.TypedData, error) {
	// Convert addresses to EVM-compatible format.
	var convSlice = []string{quote.SourceToken, permit2Data.EscrowContractAddress, permit2Data.AdditionalData.Domain.VerifyingContract}
	if err := normalizeAddresses(convSlice); err != nil {
		return apitypes.TypedData{}, err
	}

	// Unpack normalized addresses slice.
	tokenAddress, spenderAddress, verifyingContractAddress := convSlice[0], convSlice[1], convSlice[2]

	// Mask chain id.
	chainID := (*big.Int)(permit2Data.AdditionalData.Domain.ChainId).Uint64()
	chainID &= Tip712Mask

	witnessType, err := evm.ExtractWitnessType(permit2Data.AdditionalData.Types)
	if err != nil {
		return apitypes.TypedData{}, err
	}

	witnessMessage, err := evm.ParseWitnessTypedDataMessage(witnessType, &permit2Data.AdditionalData)
	if err != nil {
		return apitypes.TypedData{}, err
	}

	message := apitypes.TypedDataMessage{
		"permitted": apitypes.TypedDataMessage{
			"token":  tokenAddress,
			"amount": &quote.SourceAmountLots.Int,
		},
		"spender":  spenderAddress,
		"nonce":    big.NewInt(permit2Data.Nonce),
		"deadline": big.NewInt(deadline),
		"witness":  witnessMessage,
	}

	// INFO: Ignore domain version if any.
	return apitypes.TypedData{
		Types:       evm.InjectDomainTypeDef(permit2Data.AdditionalData.Types),
		PrimaryType: evm.PrimaryType,
		Domain: apitypes.TypedDataDomain{
			Name:              permit2Data.AdditionalData.Domain.Name,
			ChainId:           (*math.HexOrDecimal256)((&big.Int{}).SetUint64(chainID)),
			VerifyingContract: verifyingContractAddress,
		},
		Message: message,
	}, nil
}

// normalizeAddresses performs conversion to each address in the slice, updating them.
//
// NOTE: Provided addresses might be modified.
func normalizeAddresses(addresses []string) (err error) {
	for i, address := range addresses {
		addresses[i], err = tronAddressToEvmAddress(address)
		if err != nil {
			return err
		}
	}

	return nil
}

// tronAddressToEvmAddress converts a Tron address to an EVM-compatible address format.
//
// NOTE: If the provided address does not have a 'T' prefix, it will be returned as-is (assuming it's already EVM-compatible address).
func tronAddressToEvmAddress(address string) (string, error) {
	if strings.HasPrefix(address, "T") {
		address, err := tronaddress.Base58ToAddress(address)
		if err != nil {
			return "", err
		}

		return common.HexToAddress(address.Hex()[4:]).String(), nil
	}

	return address, nil
}
