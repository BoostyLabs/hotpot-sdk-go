package types

import (
	"github.com/google/uuid"
)

// CombinedStatus combines intent statuses with swaps statuses to provide a general overview on swap flow progress.
type CombinedStatus string

const (
	// CombinedStatusInitiated defines the `Initiated` status.
	CombinedStatusInitiated CombinedStatus = "Initiated"
	// CombinedStatusApprovalAdded defines the `ApprovalAdded` status.
	CombinedStatusApprovalAdded CombinedStatus = "ApprovalAdded"
	// CombinedStatusAccepted defines the `Accepted` status.
	CombinedStatusAccepted CombinedStatus = "Accepted"
	// CombinedStatusDeclined defines the `Declined` status.
	CombinedStatusDeclined CombinedStatus = "Declined"
	// CombinedStatusUserDeposited defines the `UserDeposited` status.
	CombinedStatusUserDeposited CombinedStatus = "UserDeposited"
	// CombinedStatusFulfilled defines the `Fulfilled` status.
	CombinedStatusFulfilled CombinedStatus = "Fulfilled"
	// CombinedStatusExpired defines the `Expired` status.
	CombinedStatusExpired CombinedStatus = "Expired"
	// CombinedStatusRefundRequested defines the `RefundRequested` status.
	CombinedStatusRefundRequested CombinedStatus = "RefundRequested"
	// CombinedStatusRefunded defines the `Refunded` status.
	CombinedStatusRefunded CombinedStatus = "Refunded"
)

func (s CombinedStatus) String() string {
	return string(s)
}

// Fee represents an affiliate fee that will be collected and stored for a specific intent.
type Fee struct {
	IntentID          uuid.UUID `json:"intent_id"`
	SubAffiliateID    string    `json:"sub_affiliate_id"`
	FeeBps            *Int      `json:"fee_bps"`
	NetworkID         int64     `json:"network_id"`
	Token             string    `json:"token"`
	FeeAmountLots     *Int      `json:"fee_amount_lots"`
	FeeAmountDecimals int64     `json:"fee_amount_decimals"`
	Status            FeeStatus `json:"status"`
}

// FeeStatus represents settlement status for an affiliate fee.
type FeeStatus string

const (
	// FeeStatusCreated defines that fee is in the Created status until the swap is finished for the user (swap_status < Fulfilled).
	FeeStatusCreated FeeStatus = "Created"
	// FeeStatusToSettle defines that fee moves to ToSettle when the swap is confirmed and completed for the user (swap_status = Fulfilled).
	FeeStatusToSettle FeeStatus = "ToSettle"
	// FeeStatusCancelled defines that fee is set to Cancelled when the swap fails and transitions to Refunding or Refunded.
	FeeStatusCancelled FeeStatus = "Cancelled"
	// FeeStatusSettled defines that fee is set to Settled when the swap is completed and the fee has been settled to the affiliate.
	FeeStatusSettled FeeStatus = "Settled"
)
