package types_test

import (
	"encoding/json"
	"testing"

	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

func TestApprovalToSignRgbLock_Unmarshal(t *testing.T) {
	t.Parallel()

	raw := `{
		"psbt": "cHNidP8BAGQAAAAA",
		"witness_invoice": "wvout:bcrt:rgb1...",
		"lock_address": "bcrt1qlockaddress",
		"inputs": [0, 1]
	}`

	var rgbLock types.ApprovalToSignRgbLock
	if err := json.Unmarshal([]byte(raw), &rgbLock); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if rgbLock.Psbt != "cHNidP8BAGQAAAAA" {
		t.Fatalf("psbt: got %q", rgbLock.Psbt)
	}
	if rgbLock.WitnessInvoice != "wvout:bcrt:rgb1..." {
		t.Fatalf("witness_invoice: got %q", rgbLock.WitnessInvoice)
	}
	if rgbLock.LockAddress != "bcrt1qlockaddress" {
		t.Fatalf("lock_address: got %q", rgbLock.LockAddress)
	}
	if len(rgbLock.Inputs) != 2 || rgbLock.Inputs[0] != 0 || rgbLock.Inputs[1] != 1 {
		t.Fatalf("inputs: got %v", rgbLock.Inputs)
	}
}

func TestNewRgbLockIntentApproval_MarshalJSON(t *testing.T) {
	t.Parallel()

	approval := types.NewRgbLockIntentApproval("cHNidP8BAGQAAAAA")

	data, err := json.Marshal(&approval)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal encoded approval: %v", err)
	}

	if decoded["type"] != "psbt" {
		t.Fatalf("type: got %v", decoded["type"])
	}
	if decoded["signed_data"] != "cHNidP8BAGQAAAAA" {
		t.Fatalf("signed_data: got %v", decoded["signed_data"])
	}
}
