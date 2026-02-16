# Solana Cosign Utility

This package provides utilities for partially signing Solana transactions, specifically for cosigning operations within the Hotpot SDK.

## Signer

The `Signer` struct is an abstraction over a Solana private key used for cosigning transactions.

### Functions

#### `NewSigner(privateKeyHex string) (*Signer, error)`
Creates a new `Signer` instance from a private key hex string.

#### `SignDepositTx(signer *Signer, cosignTransaction string) (string, error)`
Signs the provided hex-encoded cosign transaction and returns the updated, hex-encoded transaction.

#### `(s *Signer) SignCosignTransaction(tx *solanago.Transaction) error`
A method on `Signer` that partially signs a Solana transaction using the associated private key.

#### `ParseCosignTransaction(cosignTransactionHex string) (*solanago.Transaction, error)`
Parses a hex-encoded Solana transaction into a `*solanago.Transaction` object.

## Usage Example

The following example demonstrates how to use the `Signer` to sign a Solana cosign transaction.

```go
package main

import (
	"fmt"
	"log"

	"github.com/BoostyLabs/hotpot-sdk-go/crypto/solana"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

func main() {
	var (
		// Replace with your hex-encoded private key
		signerPrivateKeyHex string = "your-private-key-hex"

		// This is a part of the response of the `create-intent` API call for the Solana-source network.
		approvalToSign = types.ApprovalToSign{
			ApprovalMechanism: types.ApprovalToSignTypeCosign,
			Cosign: &types.ApprovalToSignCosign{
				Transaction: "your-transaction-hex", // Replace with your hex-encoded transaction.
				Nonce:       0,
			},
		}
	)

	// 1. Initialize the signer
	signer, err := solana.NewSigner(signerPrivateKeyHex)
	if err != nil {
		log.Fatalf("failed to create signer: %v", err)
	}

	// 2. Sign the deposit transaction
	signedTxHex, err := solana.SignDepositTx(signer, approvalToSign.Cosign.Transaction)
	if err != nil {
		log.Fatalf("failed to sign deposit tx: %v", err)
	}

	// 3. Use the signed transaction hex as an approval.
	fmt.Println("Signed Transaction Hex:", signedTxHex)
	_ = types.NewCosignIntentApproval(signedTxHex, "your-address") // Replace with your address.
}
```

## Implementation Details

- **Partial Signing**: The implementation handles partial signing of Solana transactions, allowing multiple parties to sign the same transaction.
- **Hex Encoding**: Transactions are expected to be passed and returned as hex-encoded strings.
- **Dependency**: It relies on `github.com/gagliardetto/solana-go` for Solana primitives and transaction handling.
