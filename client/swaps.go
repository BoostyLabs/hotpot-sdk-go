package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// GetSwapByIntentIDResponse represents the response payload for retrieving a swap by its intent ID.
type GetSwapByIntentIDResponse = types.Swap

// GetSwapByIntentID retrieves a swap by its intent ID.
func (c *Client) GetSwapByIntentID(ctx context.Context, intentID uuid.UUID) (GetSwapByIntentIDResponse, error) {
	var resp GetSwapByIntentIDResponse
	endpoint := c.buildURL("swaps/intents/%s", intentID.String())

	return resp, c.doRequest(ctx, http.MethodGet, endpoint, nil, &resp)
}

// ListSwapHistoryParams represents parameters for listing swap history.
type ListSwapHistoryParams struct {
	// Limit defines the maximum number of swaps to return.
	Limit int64
	// Offset defines the number of swaps to skip.
	Offset int64
	// Network defines the swap network filter (optional).
	Network int64
	// From defines the swap `from` timestamp filter (optional).
	From int64
	// To defines the swap `to` timestamp filter (optional).
	To int64
	// Token defines the swap token filter (optional).
	Token string
	// Wallets specify addresses filter.
	Wallets []string
	// Statuses specify statusses filter.
	Statuses []types.CombinedStatus
	// RetailID specifies the retail ID filter.
	RetailID string
}

func (params *ListSwapHistoryParams) toQueryParams() string {
	q := make(url.Values, 5)

	if params.Limit != 0 {
		q.Set("limit", strconv.FormatInt(params.Limit, 10))
	}

	if params.Offset != 0 {
		q.Set("offset", strconv.FormatInt(params.Offset, 10))
	}

	if params.Network != 0 {
		q.Set("network", strconv.FormatInt(params.Network, 10))
	}

	if params.From != 0 {
		q.Set("from", strconv.FormatInt(params.From, 10))
	}

	if params.To != 0 {
		q.Set("to", strconv.FormatInt(params.To, 10))
	}

	if params.Token != "" {
		q.Set("token", params.Token)
	}

	for _, wallet := range params.Wallets {
		q.Add("wallet", wallet)
	}

	if params.RetailID != "" {
		q.Set("retail_id", params.RetailID)
	}

	for _, status := range params.Statuses {
		q.Add("status", status.String())
	}

	return q.Encode()
}

// ListSwapHistoryResponse represents the response payload for retrieving a swap history.
type ListSwapHistoryResponse Page[types.Swap]

// ListSwapHistory retrieves a swap by its intent ID.
func (c *Client) ListSwapHistory(ctx context.Context, params ListSwapHistoryParams) (ListSwapHistoryResponse, error) {
	var resp = ListSwapHistoryResponse{}
	var endpoint = c.buildURL("swaps/history?%s", params.toQueryParams())

	return resp, c.doRequest(ctx, http.MethodGet, endpoint, nil, &resp)
}
