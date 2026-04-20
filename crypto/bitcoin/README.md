# Bitcoin Signing Utility

This package provides utilities for signing Bitcoin PSBTs (Partially Signed Bitcoin Transactions) supporting both P2TR (Taproot) and P2WPKH (native SegWit) address types.

## Key types

### `PrivateKey`

A universal private key holder. Used for parsing only — signer implementations derive their address-type-specific key from it.

| Constructor | Input |
|---|---|
| `NewPrivateKeyFromHex(hex string)` | Raw 32-byte private key, hex-encoded |
| `NewPrivateKeyFromWIF(wif string)` | WIF-encoded private key |
| `NewPrivateKeyFromMnemonic(mnemonic, passphrase string)` | BIP39 mnemonic phrase |

### `Signer` interface

Both signer implementations satisfy this interface:

```go
type Signer interface {
    Address()   btcutil.Address
    PublicKey() *btcec.PublicKey
    SignInput(packet *psbt.Packet, input int) error
}
```

### `P2TRSigner`

Signs P2TR (Taproot) inputs via the key-spend path. For mnemonic-based keys, derives using **BIP86** path `m/86'/0'/0'/0/0`.

```go
signer, err := bitcoin.NewP2TRSigner(key, network)
```

### `P2WPKHSigner`

Signs P2WPKH (native SegWit) inputs. For mnemonic-based keys, derives using **BIP84** path `m/84'/0'/0'/0/0`.

```go
signer, err := bitcoin.NewP2WPKHSigner(key, network)
```

## Functions

#### `SignDepositTx(signer Signer, psbtB64 string, inputsToSig []int) (string, error)`

Signs the specified inputs of a base64-encoded PSBT and returns the updated PSBT in base64.

#### `ParsePSBT(psbtB64 string) (*psbt.Packet, error)`

Parses a base64-encoded PSBT into a `psbt.Packet`.

## Usage examples

### From hex private key

```go
key, err := bitcoin.NewPrivateKeyFromHex("your-private-key-hex")
if err != nil {
    log.Fatal(err)
}

signer, err := bitcoin.NewP2TRSigner(key, &chaincfg.MainNetParams)
if err != nil {
    log.Fatal(err)
}

fmt.Println(signer.Address()) // bc1p...

signedPsbt, err := bitcoin.SignDepositTx(signer, psbtBase64, []int{0, 1})
```

### From WIF

```go
key, err := bitcoin.NewPrivateKeyFromWIF("your-wif-encoded-key")
if err != nil {
    log.Fatal(err)
}

signer, err := bitcoin.NewP2WPKHSigner(key, &chaincfg.MainNetParams)
```

### From mnemonic

```go
key, err := bitcoin.NewPrivateKeyFromMnemonic("word1 word2 ... word12", "")
if err != nil {
    log.Fatal(err)
}

// Each signer derives its own key via the appropriate BIP32 path.
p2trSigner,   err := bitcoin.NewP2TRSigner(key, &chaincfg.MainNetParams)   // BIP86
p2wpkhSigner, err := bitcoin.NewP2WPKHSigner(key, &chaincfg.MainNetParams) // BIP84
```

### Signing a deposit transaction

```go
package main

import (
	"fmt"
	"log"

	"github.com/BoostyLabs/hotpot-sdk-go/crypto/bitcoin"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
	"github.com/btcsuite/btcd/chaincfg"
)

func main() {
	var (
		// Replace with your hex-encoded private key (or use NewPrivateKeyFromWIF / NewPrivateKeyFromMnemonic).
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

	// 1. Parse the private key.
	key, err := bitcoin.NewPrivateKeyFromHex(signerPrivateKeyHex)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}

	// 2. Initialize the signer for the desired address type.
	signer, err := bitcoin.NewP2TRSigner(key, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("failed to create signer: %v", err)
	}

	// 3. Sign the deposit transaction.
	signedPsbtBase64, err := bitcoin.SignDepositTx(signer, approvalToSign.Htlc.Psbt, approvalToSign.Htlc.Inputs)
	if err != nil {
		log.Fatalf("failed to sign deposit tx: %v", err)
	}

	// 4. Use the signed PSBT as an approval.
	fmt.Println("Signed PSBT Base64:", signedPsbtBase64)
	_ = types.NewHtlcIntentApproval(signedPsbtBase64)
}
```
