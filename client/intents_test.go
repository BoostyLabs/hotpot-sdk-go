package client_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/BoostyLabs/hotpot-sdk-go/client"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

func TestCreateIntentResponse_Unmarshal_RgbLock(t *testing.T) {
	t.Parallel()

	intentID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	raw := `{
		"intent_id": "550e8400-e29b-41d4-a716-446655440000",
		"deadline_secs": 1762999139,
		"secret_hash": "abc123",
		"approval_mechanism": "rgblock",
		"params_to_sign": {
			"psbt": "cHNidP8BAGQAAAAA",
			"witness_invoice": "wvout:bcrt:rgb1...",
			"lock_address": "bcrt1qlockaddress",
			"inputs": [0, 1]
		}
	}`

	var resp client.CreateIntentResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.ID != intentID {
		t.Fatalf("intent id: got %v", resp.ID)
	}
	if resp.ApprovalMechanism != types.ApprovalToSignTypeRgbLock {
		t.Fatalf("approval mechanism: got %q", resp.ApprovalMechanism)
	}
	if resp.RgbLock == nil {
		t.Fatal("expected rgb lock params")
	}
	if resp.RgbLock.Psbt != "cHNidP8BAGQAAAAA" {
		t.Fatalf("psbt: got %q", resp.RgbLock.Psbt)
	}
	if resp.RgbLock.WitnessInvoice != "wvout:bcrt:rgb1..." {
		t.Fatalf("witness_invoice: got %q", resp.RgbLock.WitnessInvoice)
	}
	if resp.RgbLock.LockAddress != "bcrt1qlockaddress" {
		t.Fatalf("lock_address: got %q", resp.RgbLock.LockAddress)
	}
}
