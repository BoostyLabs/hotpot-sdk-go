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
	// Limit defines the maximum number of tokens to return.
	Limit int64
	// Offset defines the number of tokens to skip.
	Offset int64
	// Active specifies whether to return only active swaps.
	Active bool
	// Wallets specify addresses filter.
	Wallets []string
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

	q.Set("active", strconv.FormatBool(params.Active))

	for _, wallet := range params.Wallets {
		q.Add("wallet", wallet)
	}

	if params.RetailID != "" {
		q.Set("retail_id", params.RetailID)
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
