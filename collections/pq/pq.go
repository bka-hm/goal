// Package pq provides a priority queue (min-queue) interface.
package pq

import (
	"cmp"

	"gitlab.lrz.de/hm/goal/base"
)

// PriorityQueue is a Priority Min-Queue. It stores Elements of a generic type E with a priority of type K.
// To align with literature, the priority is called "key" throughout
type PriorityQueue[E any, K cmp.Ordered] interface {

	// Insert adds an element with the given key and returns a handle to it.
	// Must accept non-unique keys.
	Insert(element E, key K) base.Handle

	// DecreaseKey decreases the key of the element identified by handle.
	// Panics if the handle is invalid or the new key is greater than the current key.
	DecreaseKey(handle base.Handle, key K)

	// Remove removes an element identified by a handle.
	Remove(handle base.Handle)

	// ExtractMin removes and returns the element with the smallest key.
	// Panics if the queue is empty.
	ExtractMin() (element E, key K)

	// Min returns the element with the smallest key without removing it.
	// Panics if the queue is empty.
	Min() (element E, key K)

	// Key returns the current key of the element identified by handle.
	Key(handle base.Handle) (key K)

	// IsEmpty returns whether the queue is empty
	IsEmpty() bool

	// Size returns the number of elements in the queue.
	Size() int
}
