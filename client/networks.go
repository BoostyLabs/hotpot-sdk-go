package client

import (
	"context"
	"net/http"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// ListNetworkParams represents the parameters used to list networks.
type ListNetworkParams struct {
	Token string
}

// ListNetworkResponse represents the response type for listing networks.
type ListNetworkResponse = []types.Network

// ListNetworks returns the list of networks, optionally providing the token by which to list networks.
func (c *Client) ListNetworks(ctx context.Context, params ListNetworkParams) (ListNetworkResponse, error) {
	var resp = make(ListNetworkResponse, 0)

	var endpoint string
	if params.Token == "" {
		endpoint = c.buildURL("networks")
	} else {
		endpoint = c.buildURL("networks?token=%s", params.Token)
	}

	return resp, c.doRequest(ctx, http.MethodGet, endpoint, nil, &resp)
}
