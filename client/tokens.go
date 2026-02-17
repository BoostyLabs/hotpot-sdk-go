package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

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

func (params *ListTokenParams) toQueryParams() string {
	q := make(url.Values, 4)

	if params.Limit != 0 {
		q.Set("limit", strconv.FormatInt(params.Limit, 10))
	}

	if params.Offset != 0 {
		q.Set("offset", strconv.FormatInt(params.Offset, 10))
	}

	if params.Query != "" {
		q.Set("q", params.Query)
	}

	if params.NetworkID != 0 {
		q.Set("network_id", strconv.FormatInt(params.NetworkID, 10))
	}

	return q.Encode()
}

// ListTokenResponse represents the response type for listing tokens.
type ListTokenResponse Page[types.Token]

// ListTokens returns the list of all supported tokens with optional filtration.
func (c *Client) ListTokens(ctx context.Context, params ListTokenParams) (ListTokenResponse, error) {
	var resp = ListTokenResponse{}
	var endpoint = c.buildURL("tokens?%s", params.toQueryParams())

	return resp, c.doRequest(ctx, http.MethodGet, endpoint, nil, &resp)
}
