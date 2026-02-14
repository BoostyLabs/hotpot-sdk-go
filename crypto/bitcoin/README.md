# Bitcoin Signing Utility

This package provides utilities for signing Bitcoin Partially Signed Bitcoin Transactions (PSBT) specifically for P2TR (Taproot) key-spend paths.

## Signer

The `Signer` struct is an abstraction over a Bitcoin private key used for P2TR key-spend path signatures.

### Functions

#### `NewSigner(privateKeyHex string) (*Signer, error)`
Creates a new `Signer` instance from a private key hex string.

#### `SignDepositTx(signer *Signer, psbtB64 string, inputsToSig []int) (string, error)`
Signs the specified inputs of a base64-encoded PSBT using the P2TR key-spend path and returns the updated, base64-encoded PSBT.

#### `SignPsbtInputKeySpend(packet *psbt.Packet, signInputIndex int) (*psbt.Packet, error)`
A method on `Signer` that signs a specific PSBT input with a key-spend path and returns the updated `psbt.Packet`.

## Usage Example

The following example demonstrates how to use the `Signer` to sign a PSBT.

```go
package main

import (
	"fmt"
	"log"

	"github.com/BoostyLabs/hotpot-sdk-go/crypto/bitcoin"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

func main() {
	var (
		// Replace with your hex-encoded private key
		signerPrivateKeyHex string = "your-private-key-hex"

		// This is a part of the response of the `create-intent` API call for the Bitcoin-source network.
		approvalToSign = types.ApprovalToSign{
			ApprovalMechanism: types.ApprovalToSignTypeHtlc,
			Htlc: &types.ApprovalToSignHtlc{
				Psbt:   "your-psbt-base64", // Replace with your base64-encoded PSBT.
				Inputs: []int{0, 1, 2},     // Indices of the inputs you need to sign.
			},
		}
	)

	// 1. Initialize the signer
	signer, err := bitcoin.NewSigner(signerPrivateKeyHex)
	if err != nil {
		log.Fatalf("failed to create signer: %v", err)
	}

	// 2. Sign the deposit transaction
	signedPsbtBase64, err := bitcoin.SignDepositTx(signer, approvalToSign.Htlc.Psbt, approvalToSign.Htlc.Inputs)
	if err != nil {
		log.Fatalf("failed to sign deposit tx: %v", err)
	}

	// 3. Use the signed PSBT as an approval.
	fmt.Println("Signed PSBT Base64:", signedPsbtBase64)
	_ = types.NewHtlcIntentApproval(signedPsbtBase64)
}
```

## Implementation Details

- **P2TR Key-Path**: The implementation specifically handles Taproot key-spend paths.
- **Witness Serialization**: The utility includes internal helpers for serializing witness stacks correctly for PSBTs.
- **Dependency**: It relies on `github.com/btcsuite/btcd` for Bitcoin primitives and PSBT handling.
