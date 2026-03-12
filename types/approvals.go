package types

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// ApprovalToSign defines approval mechanisms used by different blockchains to transfer ownership of funds to resolver.
//
// `approvalMechanism` holds the type of approval mechanism required for signing operations, and one of
// the `permit2|htlc|cosign` fields will be not nil (corresponding).
type ApprovalToSign struct {
	ApprovalMechanism ApprovalToSignType
	Permit2           *ApprovalToSignPermit2
	Htlc              *ApprovalToSignHtlc
	Cosign            *ApprovalToSignCosign
}

// ApprovalToSignType represents the type of approval mechanism required for signing operations.
type ApprovalToSignType string

const (
	// ApprovalToSignTypePermit2 defines `permit2` approval types.
	ApprovalToSignTypePermit2 ApprovalToSignType = "permit2"
	// ApprovalToSignTypeHtlc defines `htlc` approval types.
	ApprovalToSignTypeHtlc ApprovalToSignType = "htlc"
	// ApprovalToSignTypeCosign defines `cosign` approval types.
	ApprovalToSignTypeCosign ApprovalToSignType = "cosign"
)

// ApprovalToSignPermit2 represents parameters of a permit2 approval used to authorize a resolver as a spender.
type ApprovalToSignPermit2 struct {
	EscrowContractAddress  string `json:"escrow_contract_address"`
	Permit2ContractAddress string `json:"permit2_contract_address"`
	ResolverDepositAddress string `json:"resolver_deposit_address"`
	Nonce                  int64  `json:"nonce"`

	AdditionalData Permit2AdditionalData `json:"additional_data"`
}

// Permit2AdditionalData holds additional permit2 data to simplify integration.
type Permit2AdditionalData struct {
	Domain      apitypes.TypedDataDomain `json:"domain"` // Eip712Domain.
	Types       apitypes.Types           `json:"types"`  // Eip712Types.
	Witness     json.RawMessage          `json:"witness"`
	WitnessType string                   `json:"witness_type_string"`
	WitnessHash string                   `json:"witness_hash"`
}

// ParseWitness unmarshalls raw witness data into the provided struct.
//
// NOTE: `val` must be a pointer to a struct that implements `json.Unmarshaler`.
func (p *Permit2AdditionalData) ParseWitness(val any) error {
	return json.Unmarshal(p.Witness, val)
}

// ApprovalToSignHtlc represents parameters of a htlc approval used for Bitcoin network.
type ApprovalToSignHtlc struct {
	Psbt   string `json:"psbt"`
	Inputs []int  `json:"inputs"`
}

// ApprovalToSignCosign represents parameters of a cosign approval used for Solana network.
type ApprovalToSignCosign struct {
	Transaction string `json:"transaction"`
	Nonce       int64  `json:"nonce"`
}

// IntentApproval is a container for the intent approval.
type IntentApproval struct {
	approvalMechanism ApprovalToSignType
	permit2           *string
	psbt              *string
	cosign            *CosignApproval
}

// NewPermit2IntentApproval creates a new intent approval for permit2.
func NewPermit2IntentApproval(permit2Signature string) IntentApproval {
	return IntentApproval{
		approvalMechanism: ApprovalToSignTypePermit2,
		permit2:           &permit2Signature,
	}
}

// NewHtlcIntentApproval creates a new intent approval for htlc.
func NewHtlcIntentApproval(psbt string) IntentApproval {
	return IntentApproval{
		approvalMechanism: ApprovalToSignTypeHtlc,
		psbt:              &psbt,
	}
}

// NewCosignIntentApproval creates a new intent approval for cosign.
func NewCosignIntentApproval(transaction string, userAddress string) IntentApproval {
	return IntentApproval{
		approvalMechanism: ApprovalToSignTypeCosign,
		cosign: &CosignApproval{
			Transaction: transaction,
			UserAddress: userAddress,
		},
	}
}

// MarshalJSON implements json.Marshaler interface.
func (a *IntentApproval) MarshalJSON() ([]byte, error) {
	var v = addIntentApprovalRequestCodec{Type: string(a.approvalMechanism)}
	switch {
	case a.approvalMechanism == ApprovalToSignTypePermit2 && a.permit2 != nil:
		v.SignedData = *a.permit2
	case a.approvalMechanism == ApprovalToSignTypeHtlc && a.psbt != nil:
		v.Type = "psbt"
		v.SignedData = *a.psbt
	case a.approvalMechanism == ApprovalToSignTypeCosign && a.cosign != nil:
		v.SignedData = *a.cosign
	default:
		return nil, fmt.Errorf("unrecognized approval mechanism %v or missing signed data", a.approvalMechanism)
	}

	return json.Marshal(&v)
}

type addIntentApprovalRequestCodec struct {
	Type       string `json:"type"`
	SignedData any    `json:"signed_data"`
}

// CosignApproval represents parameters of the cosign approval used for Solana network.
type CosignApproval struct {
	Transaction string `json:"transaction"`
	UserAddress string `json:"user_address"`
}
