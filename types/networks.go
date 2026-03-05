package types

// Network represents the network entity data.
type Network struct {
	ID                    int64       `json:"id"`
	Name                  string      `json:"name"`
	Type                  NetworkType `json:"type"`
	IconURL               string      `json:"icon_url"`
	SupportsOptimizedSwap bool        `json:"supports_optimized_swap"`
	SupportsCustomTokens  bool        `json:"supports_custom_tokens"`
}

// NetworkType enumerates supported blockchain runtime environments.
type NetworkType string

const (
	// NetworkTypeEVM defines EVM-compatible network types.
	NetworkTypeEVM NetworkType = "EVM"
	// NetworkTypeTron defines Tron network type.
	NetworkTypeTron NetworkType = "TRON"
	// NetworkTypeBitcoin defines Bitcoin network type.
	NetworkTypeBitcoin NetworkType = "BITCOIN"
	// NetworkTypeSolana defines Solana network type.
	NetworkTypeSolana NetworkType = "SOLANA"
)
