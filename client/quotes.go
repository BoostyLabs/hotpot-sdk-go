package client

import (
	"context"
	"net/http"

	types2 "github.com/BoostyLabs/hotpot-sdk-go/types"
)

// GetTheBestQuoteRequest represents parameters to receive the best quote.
type GetTheBestQuoteRequest struct {
	SourceChain  uint64             `json:"source_chain"`
	SourceToken  string             `json:"source_token"`
	DestChain    uint64             `json:"dest_chain"`
	DestToken    string             `json:"dest_token"`
	Amount       float64            `json:"amount"`
	Slippage     *types2.Int        `json:"slippage_bps"`
	SwapType     types2.SwapType    `json:"swap_type"`
	RetailUserID string             `json:"retail_user_id,omitempty"`
	DepositType  types2.DepositType `json:"deposit_type"`
}

// GetTheBestQuoteResponse represents the response api type for the quote request.
type GetTheBestQuoteResponse = types2.Quote

// GetTheBestQuote returns the best quote for the provided parameters.
func (c *Client) GetTheBestQuote(ctx context.Context, req GetTheBestQuoteRequest) (types2.Quote, error) {
	var resp types2.Quote
	endpoint := c.buildURL("quotes/best")

	return resp, c.doRequest(ctx, http.MethodPost, endpoint, &req, &resp)
}
