package client_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/BoostyLabs/hotpot-sdk-go/client"
	"github.com/BoostyLabs/hotpot-sdk-go/types"
)

func TestCreateIntentResponse_Unmarshal_RgbUserBroadcast(t *testing.T) {
	t.Parallel()

	intentID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	raw := `{
		"intent_id": "550e8400-e29b-41d4-a716-446655440000",
		"deadline_secs": 1762999139,
		"secret_hash": "abc123",
		"approval_mechanism": "broadcast",
		"deposit_tx_data": {
			"RgbLock": {
				"witness_invoice": "wvout:bcrt:rgb1..."
			}
		}
	}`

	var resp client.CreateIntentResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.ID != intentID {
		t.Fatalf("intent id: got %v", resp.ID)
	}
	if resp.ApprovalMechanism != types.ApprovalToSignTypeUserBroadcast {
		t.Fatalf("approval mechanism: got %q", resp.ApprovalMechanism)
	}
	if resp.DepositTxData == nil || resp.DepositTxData.RgbLock == nil {
		t.Fatal("expected RGB lock deposit calldata")
	}
	if resp.DepositTxData.RgbLock.WitnessInvoice != "wvout:bcrt:rgb1..." {
		t.Fatalf("witness_invoice: got %q", resp.DepositTxData.RgbLock.WitnessInvoice)
	}
}
