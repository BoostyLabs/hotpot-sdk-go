package evm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// EIP712Domain represents the primary type used for hashing and signing of typed data in the EIP-712 standard.
const EIP712Domain = "EIP712Domain"

// PrimaryType defines the primary type used for EIP-712 encoding in typed data structures.
const PrimaryType string = "PermitWitnessTransferFrom"

// DataTypeExtractor defines a map of field names to their corresponding raw JSON values for typed data message parsing.
type DataTypeExtractor map[string]json.RawMessage

// BuildTypedData constructs a TypedData object from the provided permit2Data with extra values.
func BuildTypedData(quote *types.Quote, permit2Data *types.ApprovalToSignPermit2, deadline int64) (apitypes.TypedData, error) {
	witnessType, err := ExtractWitnessType(permit2Data.AdditionalData.Types)
	if err != nil {
		return apitypes.TypedData{}, err
	}

	witnessMessage, err := ParseWitnessTypedDataMessage(witnessType, &permit2Data.AdditionalData)
	if err != nil {
		return apitypes.TypedData{}, err
	}

	message := apitypes.TypedDataMessage{
		"permitted": apitypes.TypedDataMessage{
			"token":  quote.SourceToken,
			"amount": &quote.SourceAmountLots.Int,
		},
		"spender":  permit2Data.EscrowContractAddress,
		"nonce":    big.NewInt(permit2Data.Nonce),
		"deadline": big.NewInt(deadline),
		"witness":  witnessMessage,
	}

	// INFO: Ignore domain version if any.
	return apitypes.TypedData{
		Types:       InjectDomainTypeDef(permit2Data.AdditionalData.Types),
		PrimaryType: PrimaryType,
		Domain: apitypes.TypedDataDomain{
			Name:              permit2Data.AdditionalData.Domain.Name,
			ChainId:           permit2Data.AdditionalData.Domain.ChainId,
			VerifyingContract: permit2Data.AdditionalData.Domain.VerifyingContract,
		},
		Message: message,
	}, nil
}

// InjectDomainTypeDef adds the EIP712Domain type definition if it doesn't already exist in the provided types definition map.
func InjectDomainTypeDef(typesDef apitypes.Types) apitypes.Types {
	updated := make(apitypes.Types, len(typesDef)+1)
	for k, v := range typesDef {
		updated[k] = v
	}
	if _, ok := typesDef[EIP712Domain]; !ok {
		updated[EIP712Domain] = []apitypes.Type{
			{Name: "name", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		}
	}

	return updated
}

// ExtractWitnessType extracts and returns the witness type definition from the provided type definitions.
func ExtractWitnessType(typesDef apitypes.Types) ([]apitypes.Type, error) {
	primaryType, ok := typesDef[PrimaryType]
	if !ok {
		return nil, fmt.Errorf("primary type %s not found", PrimaryType)
	}

	witnessPos := slices.IndexFunc(primaryType, func(t apitypes.Type) bool { return t.Name == "witness" })
	if witnessPos == -1 {
		return nil, fmt.Errorf("witness not found in primary type %s", PrimaryType)
	}

	witnessType, ok := typesDef[primaryType[witnessPos].Type]
	if !ok {
		return nil, fmt.Errorf("witness type %s not found", primaryType[witnessPos].Type)
	}

	return witnessType, nil
}

// ParseWitnessTypedDataMessage constructs a TypedDataMessage from a type definition and corresponding witness data.
func ParseWitnessTypedDataMessage(typeDef []apitypes.Type, additionalData *types.Permit2AdditionalData) (apitypes.TypedDataMessage, error) {
	extractor := make(DataTypeExtractor, len(typeDef))

	if err := additionalData.ParseWitness(&extractor); err != nil {
		return apitypes.TypedDataMessage{}, err
	}

	return UnpackTypedDataMessage(typeDef, extractor)

}

// UnpackTypedDataMessage parses each field by specified type definition into a TypedDataMessage.
func UnpackTypedDataMessage(typeDef []apitypes.Type, extractor DataTypeExtractor) (apitypes.TypedDataMessage, error) {
	message := make(apitypes.TypedDataMessage, len(typeDef))
	for _, field := range typeDef {
		val, ok := extractor[field.Name]
		if !ok {
			return apitypes.TypedDataMessage{}, fmt.Errorf("missing field %s", field.Name)
		}

		switch field.Type {
		case "address":
			var address common.Address
			if err := json.Unmarshal(val, &address); err != nil {
				return apitypes.TypedDataMessage{}, err
			}

			message[field.Name] = address.Hex()
		case "uint256":
			var number math.HexOrDecimal256
			if err := json.Unmarshal(val, &number); err != nil {
				return apitypes.TypedDataMessage{}, err
			}

			message[field.Name] = (*big.Int)(&number)
		case "bytes32":
			raw := strings.TrimPrefix(strings.Trim(string(val), "\""), "0x")
			if len(raw) != 64 {
				return apitypes.TypedDataMessage{}, fmt.Errorf("invalid bytes32 hex length %d", len(raw))
			}

			b, err := hex.DecodeString(raw)
			if err != nil {
				return apitypes.TypedDataMessage{}, err
			}

			message[field.Name] = [32]byte(b)
		default:
			return nil, fmt.Errorf("unsupported type %s", field.Type)
		}
	}

	return message, nil

}

// Permit2Hash computes the EIP-712 hash of the typed data for Permit2 transactions.
func Permit2Hash(typedData apitypes.TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct(EIP712Domain, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(
		[]byte{0x19, 0x01},
		domainSeparator,
		typedDataHash,
	), nil
}

// GetPermit2Signature signs the provided digest using the provided signing function.
func GetPermit2Signature(digest []byte, sign func([]byte) ([]byte, error)) (sig []byte, err error) {
	sig, err = sign(digest)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}

	// Normalize to 27/28 for EVM contracts if needed
	sig[64] = sig[64] + 27

	return sig, nil
}
