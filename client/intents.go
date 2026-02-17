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

func (resp *CreateIntentResponse) UnmarshalJSON(data []byte) error {
	type createIntentResponseCodec struct {
		ID                uuid.UUID                `json:"intent_id"`
		Deadline          int64                    `json:"deadline_secs"`
		SecretHash        string                   `json:"secret_hash"`
		ApprovalMechanism types.ApprovalToSignType `json:"approval_mechanism"`
		ParamsToSign      json.RawMessage          `json:"params_to_sign"`
	}

	var raw createIntentResponseCodec
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	resp.ID = raw.ID
	resp.Deadline = raw.Deadline
	resp.SecretHash = raw.SecretHash
	resp.ApprovalToSign = types.ApprovalToSign{ApprovalMechanism: raw.ApprovalMechanism}

	switch raw.ApprovalMechanism {
	case types.ApprovalToSignTypePermit2:
		resp.Permit2 = new(types.ApprovalToSignPermit2)
		return json.Unmarshal(raw.ParamsToSign, resp.Permit2)
	case types.ApprovalToSignTypeHtlc:
		resp.Htlc = new(types.ApprovalToSignHtlc)
		return json.Unmarshal(raw.ParamsToSign, resp.Htlc)
	case types.ApprovalToSignTypeCosign:
		resp.Cosign = new(types.ApprovalToSignCosign)
		return json.Unmarshal(raw.ParamsToSign, resp.Cosign)
	default:
		return fmt.Errorf("unrecognized approval mechanism %v", resp.ApprovalMechanism)
	}
}

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
