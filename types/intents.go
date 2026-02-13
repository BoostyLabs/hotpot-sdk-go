package types

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
	// CombinedStatusDeclinedDueToKytCheck defines the `DeclinedDueToKytCheck` status.
	CombinedStatusDeclinedDueToKytCheck CombinedStatus = "DeclinedDueToKytCheck"
	// CombinedStatusUserDeposited defines the `UserDeposited` status.
	CombinedStatusUserDeposited CombinedStatus = "UserDeposited"
	// CombinedStatusKycRequested defines the `KycRequested` status.
	CombinedStatusKycRequested CombinedStatus = "KycRequested"
	// CombinedStatusFulfilled defines the `Fulfilled` status.
	CombinedStatusFulfilled CombinedStatus = "Fulfilled"
	// CombinedStatusExpired defines the `Expired` status.
	CombinedStatusExpired CombinedStatus = "Expired"
	// CombinedStatusRefundRequested defines the `RefundRequested` status.
	CombinedStatusRefundRequested CombinedStatus = "RefundRequested"
	// CombinedStatusRefunded defines the `Refunded` status.
	CombinedStatusRefunded CombinedStatus = "Refunded"
)
