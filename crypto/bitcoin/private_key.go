package bitcoin

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"golang.org/x/crypto/pbkdf2"
)

// PrivateKey is a universal private key holder for hex, WIF, and mnemonic inputs.
// It is used for parsing only; signer implementations derive their specific key from it.
type PrivateKey struct {
	rawKey *btcec.PrivateKey
	// seed is the BIP39 seed for mnemonic inputs, enabling per-signer HD path derivation.
	seed []byte
}

// NewPrivateKeyFromHex creates a PrivateKey from a hex-encoded raw private key.
func NewPrivateKeyFromHex(privateKeyHex string) (PrivateKey, error) {
	b, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return PrivateKey{}, err
	}

	key, _ := btcec.PrivKeyFromBytes(b)

	return PrivateKey{rawKey: key}, nil
}

// NewPrivateKeyFromWIF creates a PrivateKey from a WIF-encoded private key.
func NewPrivateKeyFromWIF(wifStr string) (PrivateKey, error) {
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{rawKey: wif.PrivKey}, nil
}

// NewPrivateKeyFromMnemonic creates a PrivateKey from a BIP39 mnemonic phrase.
// passphrase may be empty. Each signer derives its own key via the appropriate BIP32 path.
func NewPrivateKeyFromMnemonic(mnemonic, passphrase string) (PrivateKey, error) {
	if mnemonic == "" {
		return PrivateKey{}, errors.New("mnemonic cannot be empty")
	}

	seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"+passphrase), 2048, 64, sha512.New)

	return PrivateKey{seed: seed}, nil
}

func (pk *PrivateKey) isFromMnemonic() bool {
	return pk.seed != nil
}
