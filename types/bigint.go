package types

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
)

// PercentToBps defines the conversion factor from percent to base points.
const PercentToBps float64 = 100.0

// Int is a wrapper type over integer values with high precision and custom JSON marshaling.
type Int struct{ big.Int }

// NewInt creates a new Int from an int64.
func NewInt(i int64) *Int { return &Int{*big.NewInt(i)} }

// NewIntFromBigInt creates a new Int from a big.Int.
func NewIntFromBigInt(i *big.Int) *Int { return &Int{*i} }

// NewIntFromPercent creates a new Int from a percentage value.
//
// NOTE: percent must be within the range [0.0, 100.0].
func NewIntFromPercent(percent float64) (*Int, error) {
	if percent < 0. || percent > 100. {
		return nil, fmt.Errorf("slippage percent must be between 0.0 and 100.0, got %f", percent)
	}

	return NewInt(int64(math.Round(percent * PercentToBps))), nil
}

func (i *Int) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('"')
	buf.WriteString(i.String())
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

func (i *Int) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	return i.Int.UnmarshalJSON(b)
}
