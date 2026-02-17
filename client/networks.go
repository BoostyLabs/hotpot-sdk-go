package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// ListNetworkParams represents the parameters used to list networks.
type ListNetworkParams struct {
	Token string
}

func (params *ListNetworkParams) toQueryParams() string {
	q := make(url.Values, 1)

	if params.Token != "" {
		q.Set("token", params.Token)
	}

	return q.Encode()
}

// ListNetworkResponse represents the response type for listing networks.
type ListNetworkResponse = []types.Network

// ListNetworks returns the list of networks, optionally providing the token by which to list networks.
func (c *Client) ListNetworks(ctx context.Context, params ListNetworkParams) (ListNetworkResponse, error) {
	var resp = make(ListNetworkResponse, 0)
	var endpoint = c.buildURL("networks?%s", params.toQueryParams())

	return resp, c.doRequest(ctx, http.MethodGet, endpoint, nil, &resp)
}
