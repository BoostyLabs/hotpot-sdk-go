package bitcoin

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// p2trPath is BIP86: m/86'/0'/0'/0/0 (Bitcoin mainnet, account 0, first external address).
var p2trPath = []uint32{
	86 + hdkeychain.HardenedKeyStart,
	0 + hdkeychain.HardenedKeyStart,
	0 + hdkeychain.HardenedKeyStart,
	0,
	0,
}

// P2TRSigner signs P2TR (Taproot) inputs using the key-spend path.
type P2TRSigner struct {
	pk  *btcec.PrivateKey
	net *chaincfg.Params
}

// NewP2TRSigner creates a P2TRSigner from a PrivateKey.
// For mnemonic-based keys it derives using BIP86 path m/86'/0'/0'/0/0.
func NewP2TRSigner(key PrivateKey, network *chaincfg.Params) (*P2TRSigner, error) {
	pk, err := resolveKey(key, network, p2trPath)
	if err != nil {
		return nil, err
	}

	return &P2TRSigner{pk: pk, net: network}, nil
}

// PublicKey returns the public key for the underlying private key.
func (s *P2TRSigner) PublicKey() *btcec.PublicKey {
	return s.pk.PubKey()
}

// Address returns the P2TR (Taproot) address for this key.
func (s *P2TRSigner) Address() btcutil.Address {
	tapKey := txscript.ComputeTaprootKeyNoScript(s.pk.PubKey())
	addr, _ := btcutil.NewAddressTaproot(schnorr.SerializePubKey(tapKey), s.net)
	return addr
}

// SignInput signs a P2TR key-spend input in the PSBT, updating it in place.
func (s *P2TRSigner) SignInput(packet *psbt.Packet, input int) error {
	if len(packet.UnsignedTx.TxIn) <= input || len(packet.Inputs) <= input {
		return errors.New("invalid input index")
	}

	pInput := packet.Inputs[input]
	outsMap := make(map[wire.OutPoint]*wire.TxOut, len(packet.UnsignedTx.TxIn))
	for idx, in := range packet.UnsignedTx.TxIn {
		outsMap[in.PreviousOutPoint] = packet.Inputs[idx].WitnessUtxo
	}

	prevOuts := txscript.NewMultiPrevOutFetcher(outsMap)
	witness, err := txscript.TaprootWitnessSignature(
		packet.UnsignedTx,
		txscript.NewTxSigHashes(packet.UnsignedTx, prevOuts),
		input,
		pInput.WitnessUtxo.Value,
		pInput.WitnessUtxo.PkScript,
		pInput.SighashType,
		s.pk,
	)
	if err != nil {
		return fmt.Errorf("failed to sign p2tr input: %w", err)
	}

	packet.Inputs[input].FinalScriptWitness, err = writeWitness(witness[0])
	return err
}
