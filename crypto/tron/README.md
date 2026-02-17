# Tron Signing Utility

This package provides utilities for signing Permit2 approvals on the Tron network (TIP-712), adapting Hotpot SDK's EVM Permit2 flow to Tron specifics.

## Signer

Tron signing reuses the EVM `Signer` from the `crypto/evm` package. You can create an EVM signer and use it with Tron-specific typed data built by this package.

### Functions

#### `BuildTypedData(quote *types.Quote, permit2Data *types.ApprovalToSignPermit2, deadline int64) (apitypes.TypedData, error)`
Builds TIP-712/EIP-712 compatible typed data for Tron Permit2 signing. It normalizes Tron Base58 addresses to EVM-compatible hex addresses and applies chain ID masking per TIP-712 rules.

#### `SignPermit2(signer *evm.Signer, quote *types.Quote, permit2Data *types.ApprovalToSignPermit2, deadline int64) (string, error)`
Signs the Permit2 approval for the provided signer and returns the signature in hex encoding with the '0x' prefix.

> Note: The method `(s *evm.Signer) SignPermit2(typesData apitypes.TypedData) ([]byte, error)` is implemented in the `crypto/evm` package and is used internally after Tron-specific typed data is built.

## Usage Example

The following example demonstrates how to use the EVM `Signer` together with the Tron signing utility to sign a Permit2 approval on Tron.

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/BoostyLabs/hotpot-sdk-go/crypto/evm"
    "github.com/BoostyLabs/hotpot-sdk-go/crypto/tron"
    "github.com/BoostyLabs/hotpot-sdk-go/types"
)

func main() {
    var (
        // Replace with your hex-encoded EVM private key
        signerPrivateKeyHex string = "your-private-key-hex"

        // These are parts of the response from Hotpot SDK APIs.
        quote = &types.Quote{
            // ... populate quote details
        }
        permit2Data = &types.ApprovalToSignPermit2{
            // ... populate permit2 data from the API response (may contain Tron Base58 addresses)
        }
        deadline = time.Now().Add(time.Hour).Unix()
    )

    // 1. Initialize the EVM signer
    signer, err := evm.NewSigner(signerPrivateKeyHex)
    if err != nil {
        log.Fatalf("failed to create signer: %v", err)
    }

    // 2. Sign the Permit2 approval for Tron (TIP-712)
    signedPermit2Hex, err := tron.SignPermit2(signer, quote, permit2Data, deadline)
    if err != nil {
        log.Fatalf("failed to sign permit2: %v", err)
    }

    // 3. Use the signed Permit2 signature as an approval.
    fmt.Println("Signed Permit2 Hex:", signedPermit2Hex)
    _ = types.NewPermit2IntentApproval(signedPermit2Hex)
}
```

## Implementation Details

- **TIP-712 Chain ID Masking**: The package applies `Tip712Mask` to the chain ID (`0xffffffff`) to comply with Tron TIP-712 requirements.
- **Address Normalization**: Tron Base58 addresses (prefixed with `T`) are converted to EVM-compatible hex addresses before hashing/signing. Non-`T`-prefixed values are treated as already-hex addresses.
- **Typed Data**: Builds the `PermitWitnessTransferFrom` typed data message compatible with Permit2, reusing helper logic from the EVM package.
- **Dependencies**:
  - `github.com/fbsobreira/gotron-sdk` for Tron address conversions.
  - `github.com/ethereum/go-ethereum` for common types, math, and EIP-712 typed data structures.
  - Reuses `crypto/evm` helpers for witness parsing and signing.
