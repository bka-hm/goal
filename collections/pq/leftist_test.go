package pq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreation(t *testing.T) {
	var pq PriorityQueue[string, float32] = NewLeftistHeap[string, float32]()
	assert.Equal(t, 0, pq.Size())
}

func TestSingleInsert(t *testing.T) {
	pq := NewLeftistHeap[string, float32]()
	pq.Insert("X", 2.0)
	assert.Equal(t, 1, pq.Size())
}

func TestSingleMinExtract(t *testing.T) {
	pq := NewLeftistHeap[string, float32]()
	pq.Insert("X", 2.0)
	e, k := pq.ExtractMin()
	assert.Equal(t, 0, pq.Size())
	assert.Equal(t, "X", e)
	assert.Equal(t, float32(2.0), k)
}

func TestMultipleInsert(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	pq.Insert("X", 4)
	pq.Insert("Y", 3)
	pq.Insert("Z", 7)
	pq.Insert("G", 5)
	assert.Equal(t, 4, pq.Size())
	e, _ := pq.ExtractMin()
	assert.Equal(t, "Y", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "X", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "G", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "Z", e)
}

func TestSimpleRemoveRoot(t *testing.T) {
	pq := NewLeftistHeap[string, float32]()
	xh := pq.Insert("X", 2.0)
	pq.Insert("Y", 3.0)
	pq.Insert("Z", 7.0)
	pq.Remove(xh)
	assert.Equal(t, 2, pq.Size())
	e, _ := pq.ExtractMin()
	assert.Equal(t, "Y", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "Z", e)
}

func TestRemoveRightChild(t *testing.T) {
	// Inserting keys 1,2,3 in order gives root A(left=B, right=C).
	// Removing C (a right child) covers the fix.right = nil branch in Remove.
	pq := NewLeftistHeap[string, int]()
	pq.Insert("A", 1)
	pq.Insert("B", 2)
	hC := pq.Insert("C", 3)
	pq.Remove(hC)
	assert.Equal(t, 2, pq.Size())
	e, _ := pq.ExtractMin()
	assert.Equal(t, "A", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "B", e)
}

func TestSimpleRemoveInner(t *testing.T) {
	pq := NewLeftistHeap[string, float32]()
	pq.Insert("X", 2.0)
	xh := pq.Insert("Y", 3.0)
	pq.Insert("Z", 7.0)
	pq.Insert("G", 5.0)
	pq.Remove(xh)
	assert.Equal(t, 3, pq.Size())
	e, _ := pq.ExtractMin()
	assert.Equal(t, "X", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "G", e)
}

func Test_DecreaseKey(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	pq.Insert("X", 30)
	pq.Insert("Y", 20)
	pq.Insert("Z", 40)
	h := pq.Insert("U", 60)
	pq.Insert("V", 100)
	pq.Insert("W", 100)
	pq.DecreaseKey(h, 22)
	e, x := pq.ExtractMin()
	assert.Equal(t, "Y", e)
	assert.Equal(t, 20, x)
	e, x = pq.ExtractMin()
	assert.Equal(t, "U", e)
	assert.Equal(t, 22, x)
	assert.Equal(t, 4, pq.Size())
}

func Test_MultiDecreaseKey(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	pq.Insert("X", 30)
	h := pq.Insert("Y", 50)
	pq.DecreaseKey(h, 25)
	pq.DecreaseKey(h, 20)
	assert.Equal(t, 2, pq.Size())
	e, x := pq.ExtractMin()
	assert.Equal(t, 20, x)
	assert.Equal(t, "Y", e)
	e, x = pq.ExtractMin()
	assert.Equal(t, 30, x)
	assert.Equal(t, "X", e)
}

func Test_MultiInsertOnNonEmptyHeap(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	pq.Insert("A", 10)
	pq.Insert("B", 50)
	pq.MultiInsert([]string{"C", "D", "E"}, []int{3, 7, 1})
	assert.Equal(t, 5, pq.Size())
	e, _ := pq.ExtractMin()
	assert.Equal(t, "E", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "C", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "D", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "A", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "B", e)
}

func Test_Factory(t *testing.T) {
	elements := []string{"X", "Y", "Z", "G"}
	keys := []int{4, 3, 7, 5}
	pq, _ := BuildLeftistHeap[string, int](elements, keys)

	assert.Equal(t, 4, pq.Size())
	e, _ := pq.ExtractMin()
	assert.Equal(t, "Y", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "X", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "G", e)
	e, _ = pq.ExtractMin()
	assert.Equal(t, "Z", e)
}

func TestInsertMulti(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	handles := pq.MultiInsert([]string{"X", "Y", "Z", "G"}, []int{4, 3, 7, 5})

	assert.Equal(t, 4, pq.Size())
	assert.Len(t, handles, 4)

	e, k := pq.ExtractMin()
	assert.Equal(t, "Y", e)
	assert.Equal(t, 3, k)

	e, k = pq.ExtractMin()
	assert.Equal(t, "X", e)
	assert.Equal(t, 4, k)
}

func TestNewLeftistHeapFromItems(t *testing.T) {
	pq, _ := BuildLeftistHeap([]string{"X", "Y", "Z"}, []int{4, 3, 7})

	assert.Equal(t, 3, pq.Size())

	e, k := pq.ExtractMin()
	assert.Equal(t, "Y", e)
	assert.Equal(t, 3, k)
}

func TestNewLeftistHeapFromItemsEmpty(t *testing.T) {
	pq, _ := BuildLeftistHeap([]string{}, []int{})

	assert.True(t, pq.IsEmpty())
	assert.Equal(t, 0, pq.Size())
}

func TestMerge(t *testing.T) {
	left := NewLeftistHeap[string, int]()
	left.Insert("X", 4)
	left.Insert("Z", 7)

	right := NewLeftistHeap[string, int]()
	right.Insert("Y", 3)
	right.Insert("G", 5)

	left.Merge(right)

	assert.Equal(t, 4, left.Size())
	assert.True(t, right.IsEmpty())
	assert.Equal(t, 0, right.Size())

	e, k := left.ExtractMin()
	assert.Equal(t, "Y", e)
	assert.Equal(t, 3, k)

	e, k = left.ExtractMin()
	assert.Equal(t, "X", e)
	assert.Equal(t, 4, k)

	e, k = left.ExtractMin()
	assert.Equal(t, "G", e)
	assert.Equal(t, 5, k)

	e, k = left.ExtractMin()
	assert.Equal(t, "Z", e)
	assert.Equal(t, 7, k)
}

func TestMinEmptyPanic(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	assert.Panics(t, func() { pq.Min() })
}

func TestExtractMinEmptyPanic(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	assert.Panics(t, func() { pq.ExtractMin() })
}

func TestKeyInvalidHandlePanic(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	assert.Panics(t, func() { pq.Key(nil) })
}

func TestDecreaseKeyInvalidHandlePanic(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	pq.Insert("X", 4)
	assert.Panics(t, func() { pq.DecreaseKey(nil, 1) })
}

func TestDecreaseKeyHigherKeyPanic(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	h := pq.Insert("X", 4)
	assert.Panics(t, func() { pq.DecreaseKey(h, 10) })
}

func TestRemoveInvalidHandlePanic(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	pq.Insert("X", 4)
	assert.Panics(t, func() { pq.Remove(nil) })
}

func TestMultiInsertMismatchedLengthsPanic(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	assert.Panics(t, func() { pq.MultiInsert([]string{"X", "Y"}, []int{1}) })
}

func TestMin(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	pq.Insert("X", 4)
	pq.Insert("Y", 2)
	pq.Insert("Z", 7)

	e, k := pq.Min()
	assert.Equal(t, "Y", e)
	assert.Equal(t, 2, k)
	assert.Equal(t, 3, pq.Size()) // Min must not remove the element
}

func TestKey(t *testing.T) {
	pq := NewLeftistHeap[string, int]()
	hX := pq.Insert("X", 4)
	hY := pq.Insert("Y", 2)

	assert.Equal(t, 4, pq.Key(hX))
	assert.Equal(t, 2, pq.Key(hY))

	pq.DecreaseKey(hX, 1)
	assert.Equal(t, 1, pq.Key(hX))
}
