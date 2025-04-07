package base

import (
	"math"
	"reflect"

	"golang.org/x/exp/constraints"
)

// use Number interface for generic functions that use arithmetic ops
type Number interface {
	constraints.Float | constraints.Integer
}

// Returns the max Value of the given Number-Type N
func MaxValue[N Number]() (max N) {
	var num N
	switch reflect.TypeOf(num).Name() {
	case "uint":
		var tmp uint = math.MaxUint
		max = N(tmp)
	case "uint8":
		var tmp uint8 = math.MaxUint8
		max = N(tmp)
	case "uint16":
		var tmp uint16 = math.MaxUint16
		max = N(tmp)
	case "uint32":
		var tmp uint32 = math.MaxUint32
		max = N(tmp)
	case "uint64":
		var tmp uint64 = math.MaxUint64
		max = N(tmp)
	case "int":
		var tmp int = math.MaxInt
		max = N(tmp)
	case "int8":
		var tmp int8 = math.MaxInt8
		max = N(tmp)
	case "int16":
		var tmp int16 = math.MaxInt16
		max = N(tmp)
	case "int32":
		var tmp int32 = math.MaxInt32
		max = N(tmp)
	case "int64":
		var tmp int64 = math.MaxInt64
		max = N(tmp)
	case "float32":
		var tmp float32 = math.MaxFloat32
		max = N(tmp)
	case "float64":
		var tmp float64 = math.MaxFloat64
		max = N(tmp)
	}
	return
}
