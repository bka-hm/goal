package adt

import (
	"gitlab.lrz.de/hm/goal-core/base"
	"golang.org/x/exp/constraints"
)

// PriorityQueue is a Priority Min-Queue.
type PriorityQueue[E any, K constraints.Ordered] interface {

	// Insert an element with value V and the key K and return the Handle.
	// Note: Must accept a non-unique keys
	Insert(element E, key K) base.Handle

	// DecreaseKey decreases the key of the specified handle, may error
	// if the handle is not valid for the specific implementation
	// or if the current key is lower than the value to change to.
	DecreaseKey(handle base.Handle, key K)

	// Remove removes an element given as handle.
	Remove(handle base.Handle)

	// ExtractMin returns the Value of the minimal element if it exists and
	// removes it from the queue.
	ExtractMin() (element E, key K)

	// GetMin returns the minimal element if it exists.
	GetMin() (element E, key K)

	// GetKey returns the current key of an element.
	GetKey(handle base.Handle) (key K)

	// IsEmpty shows if the queue is empty
	IsEmpty() bool

	// Size returns the number of elements in the pq
	Size() uint
}
