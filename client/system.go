package client

import (
	"context"
	"net/http"
)

// Live checks if the system is live and ready to process requests.
func (c *Client) Live(ctx context.Context) error {
	endpoint := c.buildURL("system/live")

	return c.doRequest(ctx, http.MethodGet, endpoint, nil, nil)
}
