# EVM Signing Utility

This package provides utilities for signing EVM Permit2 approvals, specifically for interacting with the Permit2 contract in the Hotpot SDK.

## Signer

The `Signer` struct is an abstraction over an EVM private key used for Permit2 approval signing.

### Functions

#### `NewSigner(privateKeyHex string) (*Signer, error)`
Creates a new `Signer` instance from a private key hex string.

#### `SignPermit2(signer *Signer, quote *types.Quote, permit2Data *types.ApprovalToSignPermit2, deadline int64) (string, error)`
Signs the Permit2 approval for the provided signer and returns the signature in hex encoding with the '0x' prefix.

#### `(s *Signer) SignPermit2(typesData apitypes.TypedData) ([]byte, error)`
A method on `Signer` that signs the Permit2 approval with the provided typed data and returns the raw signature bytes.

## Usage Example

The following example demonstrates how to use the `Signer` to sign a Permit2 approval.

```go
package main

import (
	"fmt"
	"log"

	"github.com/BoostyLabs/hotpot-sdk-go/client"
	"github.com/BoostyLabs/hotpot-sdk-go/crypto/evm"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

func main() {
	var (
		// Replace with your hex-encoded private key
		signerPrivateKeyHex string = "your-private-key-hex"

		// This is the response of the `best-quote` API call.
		quote = &types.Quote{
			// ... populate quote details
		}
		createIntentResponse = client.CreateIntentResponse{
			// ... populate real data from the API response.
		}
	)

	// 1. Initialize the signer
	signer, err := evm.NewSigner(signerPrivateKeyHex)
	if err != nil {
		log.Fatalf("failed to create signer: %v", err)
	}

	// 2. Sign the Permit2 approval
	signedPermit2Hex, err := evm.SignPermit2(signer, quote, createIntentResponse.Permit2, createIntentResponse.Deadline)
	if err != nil {
		log.Fatalf("failed to sign permit2: %v", err)
	}

	// 3. Use the signed Permit2 signature as an approval.
	fmt.Println("Signed Permit2 Hex:", signedPermit2Hex)
	_ = types.NewPermit2IntentApproval(signedPermit2Hex)
}
```

## Implementation Details

- **Permit2/EIP-712**: The implementation specifically handles EIP-712 typed data hashing and signing for the Permit2 protocol.
- **Typed Data**: The utility includes internal helpers for building and hashing EIP-712 typed data.
- **Dependency**: It relies on `github.com/ethereum/go-ethereum` for EVM primitives, cryptographic functions, and EIP-712 support.
