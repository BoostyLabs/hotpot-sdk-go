package types

import (
	"github.com/google/uuid"
)

// Quote represents a quote data for a specific cross-chain token pair.
type Quote struct {
	ID                         uuid.UUID               `json:"id"`
	SourceChain                uint64                  `json:"source_chain"`
	SourceToken                string                  `json:"source_token"`
	DestChain                  uint64                  `json:"dest_chain"`
	DestToken                  string                  `json:"dest_token"`
	IntermediateToken          *string                 `json:"intermediate_token,omitempty"`
	IntermediateTokenAmountMin *Int                    `json:"intermediate_token_amount_min,omitempty"`
	IntermediateTokenAmountMax *Int                    `json:"intermediate_token_amount_max,omitempty"`
	IntermediateTokenDecimals  *int64                  `json:"intermediate_token_decimals,omitempty"`
	SourceAmountLots           *Int                    `json:"source_amount_lots"`
	SourceAmountDecimals       int64                   `json:"source_amount_decimals"`
	MinDestAmountLots          *Int                    `json:"min_dest_amount_lots"`
	MaxDestAmountLots          *Int                    `json:"max_dest_amount_lots"`
	DestAmountDecimals         int64                   `json:"dest_amount_decimals"`
	SlippageBps                *Int                    `json:"slippage_bps"`
	Expiry                     int64                   `json:"expiry"`
	SwapType                   SwapType                `json:"swap_type"`
	DepositType                DepositType             `json:"deposit_type"`
	AffiliateFees              map[string]EstimatedFee `json:"affiliate_fees"`
}

// SwapType defines what swap type (optimized single-chain or standard cross-chain) to use.
type SwapType string

const (
	// SwapTypeStandard defines a cross-chain swap.
	SwapTypeStandard SwapType = "standard"
	// SwapTypeOptimized defines a single-chain swap.
	SwapTypeOptimized SwapType = "optimized"
)

// DepositType defines what deposit type (escrowed or direct) to use.
type DepositType string

const (
	// DepositTypeEscrowed defines a deposit that uses a proxy contract to escrow funds.
	DepositTypeEscrowed DepositType = "escrowed"
	// DepositTypeDirect defines a deposit that transfers funds directly to resolver.
	DepositTypeDirect DepositType = "direct"
)

type EstimatedFee struct {
	FeeBps            *Int   `json:"fee_bps"`
	NetworkID         uint64 `json:"network_id"`
	Token             string `json:"token"`
	FeeAmountLots     *Int   `json:"fee_amount_lots"`
	FeeAmountDecimals int64  `json:"fee_amount_decimals"`
}
