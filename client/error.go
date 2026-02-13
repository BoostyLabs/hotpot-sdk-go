package client

import (
	"fmt"
)

// ApiError represents an error returned by the API.
type ApiError struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Entity  EntityKind `json:"entity"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("api error (code: %d, message: %s, entity: %s)", e.Code, e.Message, e.Entity)
}

// EntityKind represents the kind of entity that caused the error.
type EntityKind string

const (
	// ApprovalEntityKind defines 'Approval' entity kind.
	ApprovalEntityKind EntityKind = "Approval"
	// IntentEntityKind defines 'Intent' entity kind.
	IntentEntityKind EntityKind = "Intent"
	// NetworkEntityKind defines 'Network' entity kind.
	NetworkEntityKind EntityKind = "Network"
	// QuoteEntityKind defines 'Quote' entity kind.
	QuoteEntityKind EntityKind = "Quote"
	// ResolverEntityKind defines 'Resolver' entity kind.
	ResolverEntityKind EntityKind = "Resolver"
	// SupportedTokenEntityKind defines 'SupportedToken' entity kind.
	SupportedTokenEntityKind EntityKind = "SupportedToken"
	// SwapEntityKind defines 'Swap' entity kind.
	SwapEntityKind EntityKind = "Swap"
	// TokenEntityKind defines 'Token' entity kind.
	TokenEntityKind EntityKind = "Token"
	// UserNonceEntityKind defines 'UserNonce' entity kind.
	UserNonceEntityKind EntityKind = "UserNonce"
	// TransactionEntityKind defines 'Transaction' entity kind.
	TransactionEntityKind EntityKind = "Transaction"
	// TransactionReceiptEntityKind defines 'TransactionReceipt' entity kind.
	TransactionReceiptEntityKind EntityKind = "TransactionReceipt"
	// InternalEventEntityKind defines 'InternalEvent' entity kind.
	InternalEventEntityKind EntityKind = "InternalEvent"
	// BlockchainEventEntityKind defines 'BlockchainEvent' entity kind.
	BlockchainEventEntityKind EntityKind = "BlockchainEvent"
	// WatcherRequestEntityKind defines 'WatcherRequest' entity kind.
	WatcherRequestEntityKind EntityKind = "WatcherRequest"
	// WebhooksEntityKind defines 'Webhooks' entity kind.
	WebhooksEntityKind EntityKind = "Webhooks"
	// AffiliateEntityKind defines 'Affiliate' entity kind.
	AffiliateEntityKind EntityKind = "Affiliate"
	// AddressEntityKind defines 'Address' entity kind.
	AddressEntityKind EntityKind = "Address"
	// UserWithdrawTaskEntityKind defines 'UserWithdrawTask' entity kind.
	UserWithdrawTaskEntityKind EntityKind = "UserWithdrawTask"
	// OtherEntityKind defines 'Other' entity kind.
	OtherEntityKind EntityKind = "Other"
	// HtlcEntityKind defines 'htlc' entity kind.
	HtlcEntityKind EntityKind = "htlc"
	// ResolverApiClientEntityKind defines 'ResolverApiClient' entity kind.
	ResolverApiClientEntityKind EntityKind = "ResolverApiClient"
)
