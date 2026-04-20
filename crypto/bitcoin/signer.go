package bitcoin

import (
	"bytes"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
)

// Signer defines an implementation facade over a private key for a specific address type.
type Signer interface {
	// Address returns the address derived from the underlying private key.
	Address() btcutil.Address
	// PublicKey returns the public key for the underlying private key.
	PublicKey() *btcec.PublicKey
	// SignInput signs the required input in the provided psbt, updating it in place.
	SignInput(packet *psbt.Packet, input int) error
}

// resolveKey returns the final *btcec.PrivateKey from pk.
// For mnemonic-based keys it performs BIP32 HD derivation using the provided path.
func resolveKey(pk PrivateKey, network *chaincfg.Params, path []uint32) (*btcec.PrivateKey, error) {
	if !pk.isFromMnemonic() {
		return pk.rawKey, nil
	}

	master, err := hdkeychain.NewMaster(pk.seed, network)
	if err != nil {
		return nil, err
	}

	k := master
	for _, idx := range path {
		k, err = k.Derive(idx)
		if err != nil {
			return nil, err
		}
	}

	return k.ECPrivKey()
}

// writeWitness serializes a witness stack from the given items.
func writeWitness(stackElements ...[]byte) ([]byte, error) {
	var (
		buf          bytes.Buffer
		witnessItems = append([][]byte{}, stackElements...)
	)

	if err := psbt.WriteTxWitness(&buf, witnessItems); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
