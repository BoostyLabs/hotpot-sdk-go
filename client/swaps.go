package client

import (
	"context"
	"net/http"
	"strings"

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

// ListSwapHistoryResponse represents the response payload for retrieving a swap history.
type ListSwapHistoryResponse Page[types.Swap]

// ListSwapHistory retrieves a swap by its intent ID.
func (c *Client) ListSwapHistory(ctx context.Context, params ListSwapHistoryParams) (ListSwapHistoryResponse, error) {
	var resp = ListSwapHistoryResponse{}
	var queryParams = make([]string, 0, 4+len(params.Wallets))
	var args = make([]any, 0, 4+len(params.Wallets))

	if params.Limit != 0 {
		queryParams = append(queryParams, "limit=%d")
		args = append(args, params.Limit)
	}

	if params.Offset != 0 {
		queryParams = append(queryParams, "offset=%d")
		args = append(args, params.Offset)
	}

	queryParams = append(queryParams, "active=%t")
	args = append(args, params.Active)

	for _, wallet := range params.Wallets {
		queryParams = append(queryParams, "wallet=%s")
		args = append(args, wallet)
	}

	if params.RetailID != "" {
		queryParams = append(queryParams, "retail_id=%s")
		args = append(args, params.RetailID)
	}

	var endpoint = c.buildURL("swaps/history?"+strings.Join(queryParams, "&"), args...)

	return resp, c.doRequest(ctx, http.MethodGet, endpoint, nil, &resp)
}
