package client

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// PrepareRgbRefundRequest is the body for POST /v1/intents/{id}/rgb-refund.
type PrepareRgbRefundRequest struct {
	RefunderAddress string `json:"refunder_address"`
}

// PrepareRgbRefund returns an unsigned RGB timeout-refund offer for the intent.
// Requires a resolver API key (RFQ serves this on the resolver intents API).
func (c *Client) PrepareRgbRefund(ctx context.Context, intentID uuid.UUID, req PrepareRgbRefundRequest) (types.RgbRefundOffer, error) {
	var resp types.RgbRefundOffer
	endpoint := c.buildURL("intents/%s/rgb-refund", intentID.String())
	return resp, c.doRequest(ctx, http.MethodPost, endpoint, &req, &resp)
}
