package bitcoin

import (
	"strings"

	"github.com/btcsuite/btcd/btcutil/psbt"
)

// SignDepositTx signs the provided inputs of the raw psbt, returning the updated psbt in base64.
func SignDepositTx(signer Signer, psbtB64 string, inputsToSig []int) (string, error) {
	packet, err := ParsePSBT(psbtB64)
	if err != nil {
		return "", err
	}

	for _, idx := range inputsToSig {
		if err = signer.SignInput(packet, idx); err != nil {
			return "", err
		}
	}

	return packet.B64Encode()
}

// ParsePSBT parses a base64-encoded PSBT.
func ParsePSBT(psbtB64 string) (*psbt.Packet, error) {
	return psbt.NewFromRawBytes(strings.NewReader(psbtB64), true)
}
