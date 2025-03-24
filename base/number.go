package base

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Number is a type constraint for generic functions that use arithmetic ops.
type Number interface {
	constraints.Float | constraints.Integer
}

// MaxValue returns the maximum value of the given Number type N.
func MaxValue[N Number]() N {
	var result N
	switch p := any(&result).(type) {
	case *uint:
		*p = math.MaxUint
	case *uint8:
		*p = math.MaxUint8
	case *uint16:
		*p = math.MaxUint16
	case *uint32:
		*p = math.MaxUint32
	case *uint64:
		*p = math.MaxUint64
	case *int:
		*p = math.MaxInt
	case *int8:
		*p = math.MaxInt8
	case *int16:
		*p = math.MaxInt16
	case *int32:
		*p = math.MaxInt32
	case *int64:
		*p = math.MaxInt64
	case *float32:
		*p = math.MaxFloat32
	case *float64:
		*p = math.MaxFloat64
	default:
		panic("unsupported type")
	}
	return result
}
