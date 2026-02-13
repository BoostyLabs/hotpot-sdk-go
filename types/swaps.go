package types

import (
	"time"

	"github.com/google/uuid"
)

// Swap represents the swap entity data.
type Swap struct {
	IntentID       uuid.UUID          `json:"intent_id"`
	Status         CombinedStatus     `json:"status"`
	KycLink        *string            `json:"kyc_link,omitempty"`
	Metadata       *SwapMetadata      `json:"swap_metadata,omitempty"`
	AdditionalInfo SwapAdditionalInfo `json:"additional_info"`
}

// SwapMetadata holds swap extra data.
type SwapMetadata struct {
	ID                uuid.UUID  `json:"id"`
	ProxyAddress      string     `json:"proxy_address"`
	UserDepositTx     *string    `json:"user_deposit_tx,omitempty"`
	FulfillTx         *string    `json:"fulfill_tx,omitempty"`
	SwapTx            *string    `json:"swap_tx,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UserDepositedAt   *time.Time `json:"user_deposited_at,omitempty"`
	KycRequestedAt    *time.Time `json:"kyc_requested_at,omitempty"`
	FulfilledAt       *time.Time `json:"fulfilled_at,omitempty"`
	SwappedAt         *time.Time `json:"swapped_at,omitempty"`
	RefundedAt        *time.Time `json:"refunded_at,omitempty"`
	RefundRequestedAt *time.Time `json:"refund_requested_at,omitempty"`
}

// SwapAdditionalInfo holds quote-related swap data.
type SwapAdditionalInfo struct {
	SourceAmountDecimals int64  `json:"source_amount_decimals"`
	SourceAmountLots     *Int   `json:"source_amount_lots"`
	SourceChain          int64  `json:"source_chain"`
	SourceToken          string `json:"source_token"`
	DestAmountDecimals   int64  `json:"dest_amount_decimals"`
	DestChain            int64  `json:"dest_chain"`
	DestToken            string `json:"dest_token"`
	MinDestAmountLots    *Int   `json:"min_dest_amount_lots"`
	MaxDestAmountLots    *Int   `json:"max_dest_amount_lots"`
}
