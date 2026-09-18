package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// CreateIntentRequest represents the request payload for creating new intent.
type CreateIntentRequest struct {
	QuoteID                uuid.UUID `json:"quote_id"`
	UserSourcePublicKey    string    `json:"user_source_public_key,omitempty"`
	UserSourceAddress      string    `json:"user_source_address"`
	UserDestinationAddress string    `json:"user_destination_address"`
	RefundAddress          string    `json:"refund_address"`

	// ReferrerID is an optional identifier for further filtration.
	ReferrerID string `json:"referrer_id,omitempty"`
}

// CreateIntentResponse represents the response payload of creating new intent.
type CreateIntentResponse struct {
	ID         uuid.UUID
	Deadline   int64 // In seconds.
	SecretHash string

	types.ApprovalToSign
	DepositTxData *UBCalldata
}

// UBCalldata contains the transaction calldata for user to sign and broadcast.
type UBCalldata struct {
	Psbt    *UBCalldataPsbt
	Cosign  *UBCalldataCosign
	EvmLike *UBCalldataEvm
}

// UBCalldataPsbt contains the Bitcoin transaction calldata for user to sign and broadcast.
type UBCalldataPsbt struct {
	BitcoinPSBT  BitcoinPsbt `json:"bitcoin_psbt"`
	InputsToSign []int       `json:"inputs_to_sign"`
}

// UBCalldataCosign contains the Solana transaction calldata for user to sign and broadcast.
type UBCalldataCosign struct {
	Transaction            string `json:"transaction"`
	RecentBlockhash        string `json:"recent_blockhash"`
	ResolverDepositAddress string `json:"resolver_deposit_address"`
}

// UBCalldataPsbt contains the EVM or Tron transaction calldata for user to sign and broadcast.
type UBCalldataEvm struct {
	Transaction            string `json:"transaction"`
	ResolverDepositAddress string `json:"resolver_deposit_address"`
}

type BitcoinPsbt struct {
	PSBT                   string  `json:"psbt"`
	XOnlyPublicKey         *string `json:"x_only_public_key"`
	RefundControlBlock     *string `json:"refund_control_block"`
	FastRefundControlBlock *string `json:"fast_refund_control_block"`
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (u *UBCalldata) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("decode deposit_tx_data: %w", err)
	}

	if len(raw) != 1 {
		return fmt.Errorf(
			"expected exactly one deposit_tx_data variant, got %d",
			len(raw),
		)
	}

	for variant, payload := range raw {
		switch variant {
		case "Psbt":
			var value UBCalldataPsbt
			if err := json.Unmarshal(payload, &value); err != nil {
				return fmt.Errorf("decode Psbt: %w", err)
			}
			u.Psbt = &value

		case "Cosign":
			var value UBCalldataCosign
			if err := json.Unmarshal(payload, &value); err != nil {
				return fmt.Errorf("decode Cosign: %w", err)
			}
			u.Cosign = &value

		case "EvmLike":
			var value UBCalldataEvm
			if err := json.Unmarshal(payload, &value); err != nil {
				return fmt.Errorf("decode EvmLike: %w", err)
			}
			u.EvmLike = &value

		default:
			return fmt.Errorf("unrecognized deposit_tx_data variant %q", variant)
		}
	}

	return nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (resp *CreateIntentResponse) UnmarshalJSON(data []byte) error {
	type createIntentResponseCodec struct {
		ID                uuid.UUID                 `json:"intent_id"`
		Deadline          int64                     `json:"deadline_secs"`
		SecretHash        string                    `json:"secret_hash"`
		ApprovalMechanism *types.ApprovalToSignType `json:"approval_mechanism"`
		ParamsToSign      json.RawMessage           `json:"params_to_sign"`
		DepositTxData     *UBCalldata               `json:"deposit_tx_data"`
	}

	var raw createIntentResponseCodec

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	resp.ID = raw.ID
	resp.Deadline = raw.Deadline
	resp.SecretHash = raw.SecretHash
	resp.DepositTxData = raw.DepositTxData

	// No approval data.
	if raw.ApprovalMechanism == nil {
		resp.ApprovalToSign = types.ApprovalToSign{}
		return nil
	}

	resp.ApprovalToSign = types.ApprovalToSign{
		ApprovalMechanism: *raw.ApprovalMechanism,
	}

	switch *raw.ApprovalMechanism {
	case types.ApprovalToSignTypePermit2:
		resp.Permit2 = new(types.ApprovalToSignPermit2)
		if err := json.Unmarshal(raw.ParamsToSign, resp.Permit2); err != nil {
			return fmt.Errorf("decode permit2 params: %w", err)
		}

	case types.ApprovalToSignTypeHtlc:
		resp.Htlc = new(types.ApprovalToSignHtlc)
		if err := json.Unmarshal(raw.ParamsToSign, resp.Htlc); err != nil {
			return fmt.Errorf("decode htlc params: %w", err)
		}

	case types.ApprovalToSignTypeCosign:
		resp.Cosign = new(types.ApprovalToSignCosign)
		if err := json.Unmarshal(raw.ParamsToSign, resp.Cosign); err != nil {
			return fmt.Errorf("decode cosign params: %w", err)
		}

	default:
		return fmt.Errorf(
			"unrecognized approval mechanism %q",
			*raw.ApprovalMechanism,
		)
	}

	return nil
}

// CreateIntent creates a new intent using the provided request data and returns the created intent details or an error.
func (c *Client) CreateIntent(ctx context.Context, req CreateIntentRequest) (CreateIntentResponse, error) {
	var resp = CreateIntentResponse{}
	var endpoint = c.buildURL("intents")

	return resp, c.doRequest(ctx, http.MethodPost, endpoint, &req, &resp)
}

// AddIntentApprovalParams represents parameters required to submit an approval for a specific intent.
type AddIntentApprovalParams struct {
	IntentID uuid.UUID
	Approval types.IntentApproval
}

// AddIntentApproval submits approval for the intent, returns an empty body if adding approval was successful.
func (c *Client) AddIntentApproval(ctx context.Context, params AddIntentApprovalParams) error {
	endpoint := c.buildURL("intents/%s/approvals", params.IntentID.String())

	return c.doRequest(ctx, http.MethodPost, endpoint, &params.Approval, nil)
}

// SubmitDepositParams represents parameters required to report a deposit the user made themselves.
type SubmitDepositParams struct {
	IntentID uuid.UUID `json:"-"`
	TxHash   string    `json:"tx_hash"`
}

// SubmitDepositResponse represents the response from the SubmitDeposit API endpoint.
type SubmitDepositResponse struct {
	FulfillmentDeadline int64 `json:"fulfillment_deadline"`
}

// SubmitDeposit reports the transaction hash of a transfer the user broadcast themselves.
func (c *Client) SubmitDeposit(ctx context.Context, params SubmitDepositParams) (SubmitDepositResponse, error) {
	var resp SubmitDepositResponse
	endpoint := c.buildURL("intents/%s/deposit", params.IntentID.String())

	return resp, c.doRequest(ctx, http.MethodPost, endpoint, &params, &resp)
}

// SubmitSwap reports the transaction hash of a single-chain deposit transfer the user broadcast themselves.
func (c *Client) SubmitSwap(ctx context.Context, params SubmitDepositParams) error {
	endpoint := c.buildURL("intents/%s/add-swap", params.IntentID.String())

	return c.doRequest(ctx, http.MethodPost, endpoint, &params, nil)
}

// GetIntentStatusResponse represents the response from the GetIntentStatus API endpoint.
type GetIntentStatusResponse struct {
	Status types.CombinedStatus `json:"status"`
}

// GetIntentStatus returns the status of the intent.
func (c *Client) GetIntentStatus(ctx context.Context, intentID uuid.UUID) (GetIntentStatusResponse, error) {
	var resp GetIntentStatusResponse
	endpoint := c.buildURL("intents/%s/status", intentID.String())

	return resp, c.doRequest(ctx, http.MethodGet, endpoint, nil, &resp)
}
