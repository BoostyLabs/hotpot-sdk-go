package bitcoin

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// p2wpkhPath is BIP84: m/84'/0'/0'/0/0 (Bitcoin mainnet, account 0, first external address).
var p2wpkhPath = []uint32{
	84 + hdkeychain.HardenedKeyStart,
	0 + hdkeychain.HardenedKeyStart,
	0 + hdkeychain.HardenedKeyStart,
	0,
	0,
}

// P2WPKHSigner signs P2WPKH (native SegWit) inputs.
type P2WPKHSigner struct {
	pk  *btcec.PrivateKey
	net *chaincfg.Params
}

// NewP2WPKHSigner creates a P2WPKHSigner from a PrivateKey.
// For mnemonic-based keys it derives using BIP84 path m/84'/0'/0'/0/0.
func NewP2WPKHSigner(key PrivateKey, network *chaincfg.Params) (*P2WPKHSigner, error) {
	pk, err := resolveKey(key, network, p2wpkhPath)
	if err != nil {
		return nil, err
	}

	return &P2WPKHSigner{pk: pk, net: network}, nil
}

// PublicKey returns the public key for the underlying private key.
func (s *P2WPKHSigner) PublicKey() *btcec.PublicKey {
	return s.pk.PubKey()
}

// Address returns the P2WPKH (native SegWit) address for this key.
func (s *P2WPKHSigner) Address() btcutil.Address {
	pubKeyHash := btcutil.Hash160(s.pk.PubKey().SerializeCompressed())
	addr, _ := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, s.net)
	return addr
}

// SignInput signs a P2WPKH input in the PSBT, updating it in place.
func (s *P2WPKHSigner) SignInput(packet *psbt.Packet, input int) error {
	if len(packet.UnsignedTx.TxIn) <= input || len(packet.Inputs) <= input {
		return errors.New("invalid input index")
	}

	pInput := packet.Inputs[input]
	outsMap := make(map[wire.OutPoint]*wire.TxOut, len(packet.UnsignedTx.TxIn))
	for idx, in := range packet.UnsignedTx.TxIn {
		outsMap[in.PreviousOutPoint] = packet.Inputs[idx].WitnessUtxo
	}

	prevOuts := txscript.NewMultiPrevOutFetcher(outsMap)
	sigHashes := txscript.NewTxSigHashes(packet.UnsignedTx, prevOuts)

	// BIP143: scriptCode for P2WPKH is the P2PKH script over the same pubkey hash.
	pubKeyHash := btcutil.Hash160(s.pk.PubKey().SerializeCompressed())
	p2pkhAddr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, s.net)
	if err != nil {
		return fmt.Errorf("failed to build p2pkh address: %w", err)
	}

	subscript, err := txscript.PayToAddrScript(p2pkhAddr)
	if err != nil {
		return fmt.Errorf("failed to build subscript: %w", err)
	}

	sig, err := txscript.RawTxInWitnessSignature(
		packet.UnsignedTx,
		sigHashes,
		input,
		pInput.WitnessUtxo.Value,
		subscript,
		pInput.SighashType,
		s.pk,
	)
	if err != nil {
		return fmt.Errorf("failed to sign p2wpkh input: %w", err)
	}

	packet.Inputs[input].FinalScriptWitness, err = writeWitness(sig, s.pk.PubKey().SerializeCompressed())
	return err
}
