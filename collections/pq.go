package pq

import (
	"gitlab.lrz.de/hm/goal-core/base"
	"golang.org/x/exp/constraints"
)

// PriorityQueue Implements a Priority Min-Queue.
type PriorityQueue[K constraints.Ordered, V any] interface {

	// Insert an element with value V and the key K and return the Handle.
	// Note: Must accept a non-unique keys
	Insert(key K, value V) base.Handle

	// DecreaseKey decreases the key of the specified handle, may error
	// if the handle is not valid for the specific implementation
	// or if the current key is lower than the value to change to.
	DecreaseKey(handle base.Handle, key K) error

	// Remove removes an element given as handle.
	Remove(handle base.Handle) base.Handle

	// ExtractMin returns the Value of the minimal element if it exists.
	ExtractMin() (K, V)

	// IsEmpty shows if the queue is empty
	IsEmpty() bool

	// Size returns the number of elements in the pq
	Size() uint
}
