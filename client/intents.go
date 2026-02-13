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
}

// CreateIntentResponse represents the response payload of creating new intent.
type CreateIntentResponse struct {
	ID         uuid.UUID
	Deadline   int64 // In seconds.
	SecretHash string

	types.ApprovalToSign
}

func newCreateIntentResponseFromCodec(raw createIntentResponseCodec) (CreateIntentResponse, error) {
	resp := CreateIntentResponse{
		ID:         raw.ID,
		Deadline:   raw.Deadline,
		SecretHash: raw.SecretHash,
		ApprovalToSign: types.ApprovalToSign{
			ApprovalMechanism: raw.ApprovalMechanism,
		},
	}

	switch raw.ApprovalMechanism {
	case types.ApprovalToSignTypePermit2:
		resp.Permit2 = new(types.ApprovalToSignPermit2)
		return resp, json.Unmarshal(raw.ParamsToSign, resp.Permit2)
	case types.ApprovalToSignTypeHtlc:
		resp.Htlc = new(types.ApprovalToSignHtlc)
		return resp, json.Unmarshal(raw.ParamsToSign, resp.Htlc)
	case types.ApprovalToSignTypeCosign:
		resp.Cosign = new(types.ApprovalToSignCosign)
		return resp, json.Unmarshal(raw.ParamsToSign, resp.Cosign)
	default:
		return resp, fmt.Errorf("unrecognized approval mechanism %v", resp.ApprovalMechanism)
	}
}

type createIntentResponseCodec struct {
	ID                uuid.UUID                `json:"intent_id"`
	Deadline          int64                    `json:"deadline_secs"`
	SecretHash        string                   `json:"secret_hash"`
	ApprovalMechanism types.ApprovalToSignType `json:"approval_mechanism"`
	ParamsToSign      json.RawMessage          `json:"params_to_sign"`
}

func (c *Client) CreateIntent(ctx context.Context, req CreateIntentRequest) (CreateIntentResponse, error) {
	var resp = createIntentResponseCodec{}
	var endpoint = c.buildURL("intents")

	if err := c.doRequest(ctx, http.MethodPost, endpoint, &req, &resp); err != nil {
		return CreateIntentResponse{}, err
	}

	return newCreateIntentResponseFromCodec(resp)
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
