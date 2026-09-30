package types

import "github.com/google/uuid"

// RgbClaimOffer is the JSON package for an RGB HTLC claim (PSBT + preimage + RGB seals).
// Claim PSBTs are built by the resolver (not served by RFQ); this type is the shared wire
// format for claim offer files / tooling.
type RgbClaimOffer struct {
	Psbt              string  `json:"psbt"`
	Secret            string  `json:"secret"`
	WitnessInvoice    string  `json:"witness_invoice,omitempty"`
	ResolverTapScript string  `json:"resolver_tap_script,omitempty"`
	LockTxid          string  `json:"lock_txid,omitempty"`
	LockVout          *uint32 `json:"lock_vout,omitempty"`
	AssetID           string  `json:"asset_id,omitempty"`
	AssetAmount       *uint64 `json:"asset_amount,omitempty"`
}

// RgbRefundOffer mirrors RFQ `RgbRefundOfferResponse` (POST /v1/intents/{id}/rgb-refund).
type RgbRefundOffer struct {
	IntentID        uuid.UUID `json:"intent_id"`
	Psbt            string    `json:"psbt"`
	WitnessInvoice  string    `json:"witness_invoice"`
	RefundTapScript string    `json:"refund_tap_script,omitempty"`
	LockTxid        string    `json:"lock_txid"`
	LockVout        uint32    `json:"lock_vout"`
	AssetID         string    `json:"asset_id"`
	AssetAmount     uint64    `json:"asset_amount"`
	// DeadlineUnix is intent.deadline (CLTV); nLockTime on the PSBT is typically deadline+1.
	DeadlineUnix int64 `json:"deadline_unix"`
}
