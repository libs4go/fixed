package fixed

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNeg(t *testing.T) {

	// println(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-2)), nil).String())

	val, err := New(8, Int(1))

	require.NoError(t, err)

	require.Equal(t, val.RawValue, big.NewInt(100000000))

	val, err = New(8, HexRawValue(val.HexRawValue()))

	require.NoError(t, err)

	require.Equal(t, val.RawValue, big.NewInt(100000000))

}

func TestRaw(t *testing.T) {
	number := &Number{
		RawValue: big.NewInt(1),
		Decimals: 10,
	}

	println(fmt.Sprintf("%.10f", number.Float()))
}
