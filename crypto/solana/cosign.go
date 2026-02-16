package solana

import (
	"encoding/hex"

	bin "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
)

// Signer is the abstraction over the solana private key for cosign transaction signing.
type Signer struct{ pk solanago.PrivateKey }

// NewSigner creates a new Signer instance from a private key hex string.
func NewSigner(privateKeyHex string) (*Signer, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	signer := &Signer{privateKeyBytes}

	return signer, signer.pk.Validate()
}

// SignCosignTransaction partially signs a Solana transaction using the associated private key of the Signer.
func (s *Signer) SignCosignTransaction(tx *solanago.Transaction) error {
	_, err := tx.PartialSign(func(pubKey solanago.PublicKey) *solanago.PrivateKey {
		if pubKey == s.pk.PublicKey() {
			return &s.pk
		}

		return nil
	})

	return err
}

// SignDepositTx signs the provided raw cosign transaction, returning serialized transaction.
func SignDepositTx(signer *Signer, cosignTransaction string) (string, error) {
	tx, err := ParseCosignTransaction(cosignTransaction)
	if err != nil {
		return "", err
	}

	if err = signer.SignCosignTransaction(tx); err != nil {
		return "", err
	}

	txBinary, err := tx.MarshalBinary()

	return hex.EncodeToString(txBinary), err
}

// ParseCosignTransaction parses cosign transaction in hex encoding returning [*solanago.Transaction].
func ParseCosignTransaction(cosignTransactionHex string) (*solanago.Transaction, error) {
	cosignTransactionBytes, err := hex.DecodeString(cosignTransactionHex)
	if err != nil {
		return nil, err
	}

	var tx solanago.Transaction
	if err = bin.NewBinDecoder(cosignTransactionBytes).Decode(&tx); err != nil {
		return nil, err
	}

	return &tx, nil
}
