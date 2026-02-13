package client

import (
	"context"
	"net/http"
	"strings"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// ListTokenParams represents the parameters used to list tokens.
type ListTokenParams struct {
	// Limit defines the maximum number of tokens to return.
	Limit int64
	// Offset defines the number of tokens to skip.
	Offset int64
	// Query is used to search tokens by symbol or contract address. Example: "USDC".
	Query string
	// NetworkID adds filter tokens by network.
	NetworkID int64
}

// ListTokenResponse represents the response type for listing tokens.
type ListTokenResponse Page[types.Token]

// ListTokens returns the list of all supported tokens with optional filtration.
func (c *Client) ListTokens(ctx context.Context, params ListTokenParams) (ListTokenResponse, error) {
	var resp = ListTokenResponse{}
	var queryParams = make([]string, 0, 4)
	var args = make([]any, 0, 4)

	if params.Limit != 0 {
		queryParams = append(queryParams, "limit=%d")
		args = append(args, params.Limit)
	}

	if params.Offset != 0 {
		queryParams = append(queryParams, "offset=%d")
		args = append(args, params.Offset)
	}

	if params.Query != "" {
		queryParams = append(queryParams, "q=%s")
		args = append(args, params.Query)
	}

	if params.NetworkID != 0 {
		queryParams = append(queryParams, "network_id=%d")
		args = append(args, params.NetworkID)
	}

	var endpoint = c.buildURL("tokens?"+strings.Join(queryParams, "&"), args...)

	return resp, c.doRequest(ctx, http.MethodGet, endpoint, nil, &resp)
}
