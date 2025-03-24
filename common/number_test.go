package common

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestMaxValue(t *testing.T) {
	assert.True(t, MaxValue[uint]() == math.MaxUint)
	assert.True(t, MaxValue[uint8]() == math.MaxUint8)
	assert.True(t, MaxValue[uint16]() == math.MaxUint16)
	assert.True(t, MaxValue[uint32]() == math.MaxUint32)
	assert.True(t, MaxValue[uint64]() == math.MaxUint64)

	assert.True(t, MaxValue[int]() == math.MaxInt)
	assert.True(t, MaxValue[int8]() == math.MaxInt8)
	assert.True(t, MaxValue[int16]() == math.MaxInt16)
	assert.True(t, MaxValue[int32]() == math.MaxInt32)
	assert.True(t, MaxValue[int64]() == math.MaxInt64)

	assert.True(t, MaxValue[float32]() == math.MaxFloat32)
	assert.True(t, MaxValue[float64]() == math.MaxFloat64)
}
