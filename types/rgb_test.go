package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestRgbRefundOffer_Unmarshal_DeadlineUnix(t *testing.T) {
	intentID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	raw := `{
		"intent_id": "00000000-0000-0000-0000-000000000001",
		"psbt": "cHNidP8=",
		"witness_invoice": "wvout:bcrt:rgb1...",
		"refund_tap_script": "51",
		"lock_txid": "aa",
		"lock_vout": 1,
		"asset_id": "rgb:example",
		"asset_amount": 100,
		"deadline_unix": 1784927052
	}`
	var offer RgbRefundOffer
	if err := json.Unmarshal([]byte(raw), &offer); err != nil {
		t.Fatal(err)
	}
	if offer.IntentID != intentID {
		t.Fatalf("intent_id: got %v", offer.IntentID)
	}
	if offer.DeadlineUnix != 1784927052 {
		t.Fatalf("deadline_unix: got %d", offer.DeadlineUnix)
	}
	if offer.AssetAmount != 100 || offer.LockVout != 1 {
		t.Fatalf("asset/vout: amount=%d vout=%d", offer.AssetAmount, offer.LockVout)
	}
}

func TestRgbClaimOffer_Unmarshal(t *testing.T) {
	vout := uint32(0)
	amount := uint64(100)
	raw := `{
		"psbt": "cHNidP8=",
		"secret": "ab",
		"witness_invoice": "wvout:bcrt:rgb1...",
		"lock_txid": "aa",
		"lock_vout": 0,
		"asset_id": "rgb:example",
		"asset_amount": 100
	}`
	var offer RgbClaimOffer
	if err := json.Unmarshal([]byte(raw), &offer); err != nil {
		t.Fatal(err)
	}
	if offer.Secret != "ab" || offer.AssetID != "rgb:example" {
		t.Fatalf("unexpected offer: %+v", offer)
	}
	if offer.LockVout == nil || *offer.LockVout != vout {
		t.Fatalf("lock_vout: got %v", offer.LockVout)
	}
	if offer.AssetAmount == nil || *offer.AssetAmount != amount {
		t.Fatalf("asset_amount: got %v", offer.AssetAmount)
	}
}
