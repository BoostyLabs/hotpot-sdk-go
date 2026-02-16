package bitcoin

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// Signer is the abstraction over the bitcoin private key for P2TR key-spend path signature.
type Signer struct{ pk *btcec.PrivateKey }

// NewSigner creates a new Signer instance from a private key hex string.
func NewSigner(privateKeyHex string) (*Signer, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	return &Signer{privateKey}, nil
}

// SignPsbtInputKeySpend signs psbt input with a key-spend path, returning updated psbt (finalized for provided input).
func (s *Signer) SignPsbtInputKeySpend(packet *psbt.Packet, signInputIndex int) (*psbt.Packet, error) {
	if len(packet.UnsignedTx.TxIn) <= signInputIndex || len(packet.Inputs) <= signInputIndex {
		return nil, errors.New("invalid input index")
	}

	pInput := packet.Inputs[signInputIndex]
	outsMap := make(map[wire.OutPoint]*wire.TxOut, len(packet.UnsignedTx.TxIn))
	for idx, out := range packet.UnsignedTx.TxIn {
		outsMap[out.PreviousOutPoint] = packet.Inputs[idx].WitnessUtxo
	}

	prevOuts := txscript.NewMultiPrevOutFetcher(outsMap)
	witness, err := txscript.TaprootWitnessSignature(
		packet.UnsignedTx,
		txscript.NewTxSigHashes(packet.UnsignedTx, prevOuts),
		signInputIndex,
		pInput.WitnessUtxo.Value,
		pInput.WitnessUtxo.PkScript,
		pInput.SighashType,
		s.pk,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx (key path): %v", err)
	}

	packet.Inputs[signInputIndex].FinalScriptWitness, err = writeWitness(witness[0])
	if err != nil {
		return nil, err
	}

	return packet, nil
}

// SignDepositTx signs the provided inputs of the raw psbt (p2tr key-spend path), returning updated psbt (finalized for provided inputs).
func SignDepositTx(signer *Signer, psbtB64 string, inputsToSig []int) (string, error) {
	packet, err := ParsePSBT(psbtB64)
	if err != nil {
		return "", err
	}

	for _, idx := range inputsToSig {
		packet, err = signer.SignPsbtInputKeySpend(packet, idx)
		if err != nil {
			return "", err
		}
	}

	return packet.B64Encode()
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

// ParsePSBT parses psbt in base64 encoding returning Packet.
func ParsePSBT(pstbB64 string) (*psbt.Packet, error) {
	return psbt.NewFromRawBytes(strings.NewReader(pstbB64), true)
}
