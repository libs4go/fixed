package fixed

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/libs4go/errors"
)

var errVendor = errors.WithVendor("fixed")

// errors.
var (
	ErrRHS = errors.New("right hand side operator value must has same decimals", errVendor, errors.WithCode(-1))
)

// Number the fixed number present object
type Number struct {
	RawValue *big.Int `json:"raw"`
	Decimals int      `json:"decimals"`
}

// Source raw value source
type Source func(decimals int) (*big.Int, error)

// Int .
func Int(value int64) Source {
	return func(decimals int) (*big.Int, error) {
		val := big.NewInt(value)

		return BigInt(val)(decimals)
	}
}

// BigInt .
func BigInt(value *big.Int) Source {
	return func(decimals int) (*big.Int, error) {

		return new(big.Int).Mul(value, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)), nil
	}
}

// Float .
func Float(value float64) Source {
	return func(decimals int) (*big.Int, error) {
		return BigFloat(big.NewFloat(value))(decimals)
	}
}

// BigFloat .
func BigFloat(value *big.Float) Source {
	return func(decimals int) (*big.Int, error) {

		component := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

		result, _ := new(big.Float).Mul(value, new(big.Float).SetInt(component)).Int(nil)

		return result, nil
	}
}

func hexBytes(value string) ([]byte, error) {
	value = strings.TrimPrefix(value, "0x")

	if len(value)%2 != 0 {
		value = "0" + value
	}

	return hex.DecodeString(value)
}

// HexRawValue .
func HexRawValue(source string) Source {
	return func(decimals int) (*big.Int, error) {
		valueBytes, err := hexBytes(strings.TrimPrefix(source, "-"))

		if err != nil {
			return nil, errors.Wrap(err, "get hex string %s bytes error", source)
		}

		value := new(big.Int).SetBytes(valueBytes)

		println(value.String())

		if strings.HasPrefix(source, "-") {
			value = value.Neg(value)
		}

		return value, nil
	}
}

func (number *Number) String() string {
	buff, _ := json.Marshal(number)

	return string(buff)
}

// HexRawValue get rawvalue's hex string
func (number *Number) HexRawValue() string {
	return fmt.Sprintf("%x", number.RawValue)
}

// New create fixed number object
func New(decimals int, source Source) (*Number, error) {
	rawValue, err := source(decimals)

	if err != nil {
		return nil, err
	}

	return &Number{
		RawValue: rawValue,
		Decimals: decimals,
	}, nil
}

// Float .
func (number *Number) Float() *big.Float {
	component := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(number.Decimals)), nil)

	return new(big.Float).Quo(new(big.Float).SetInt(number.RawValue), new(big.Float).SetInt(component))
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
func (number *Number) Cmp(other *Number) int {
	if number.Decimals != other.Decimals {
		panic(ErrRHS)
	}

	return number.RawValue.Cmp(other.RawValue)
}

// Add x add y and return new fixed number object
func (number *Number) Add(other *Number) *Number {
	if number.Decimals != other.Decimals {
		panic(ErrRHS)
	}

	return &Number{
		Decimals: number.Decimals,
		RawValue: new(big.Int).Add(number.RawValue, other.RawValue),
	}
}

// Sub x sub y and return new fixed number object
func (number *Number) Sub(other *Number) *Number {
	if number.Decimals != other.Decimals {
		panic(ErrRHS)
	}

	return &Number{
		Decimals: number.Decimals,
		RawValue: new(big.Int).Sub(number.RawValue, other.RawValue),
	}
}

// Sign get fixed number sign
func (number *Number) Sign() int {
	return number.RawValue.Sign()
}
