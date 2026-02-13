package types

// Token describes a token approved for use within the protocol.
type Token struct {
	NetworkID           int64   `json:"network_id"`
	Name                string  `json:"name"`
	ContractAddress     string  `json:"contract_address"`
	Symbol              string  `json:"symbol"`
	IconURL             string  `json:"icon_url"`
	WrappedTokenAddress *string `json:"wrapped_token_address,omitempty"`
}
