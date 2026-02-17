package client

import (
	"context"
	"net/http"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

// GetTheBestQuoteRequest represents parameters to receive the best quote.
type GetTheBestQuoteRequest struct {
	SourceChain  uint64            `json:"source_chain"`
	SourceToken  string            `json:"source_token"`
	DestChain    uint64            `json:"dest_chain"`
	DestToken    string            `json:"dest_token"`
	Amount       float64           `json:"amount"`
	Slippage     *types.Int        `json:"slippage_bps"`
	SwapType     types.SwapType    `json:"swap_type"`
	RetailUserID string            `json:"retail_user_id,omitempty"`
	DepositType  types.DepositType `json:"deposit_type"`
}

// GetTheBestQuoteResponse represents the response api type for the quote request.
type GetTheBestQuoteResponse = types.Quote

// GetTheBestQuote returns the best quote for the provided parameters.
func (c *Client) GetTheBestQuote(ctx context.Context, req GetTheBestQuoteRequest) (types.Quote, error) {
	var resp types.Quote
	endpoint := c.buildURL("quotes/best")

	return resp, c.doRequest(ctx, http.MethodPost, endpoint, &req, &resp)
}
