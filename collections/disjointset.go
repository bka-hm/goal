package collections

import (
	"gitlab.lrz.de/hm/goal-core/base"
)

// PriorityQueue is a Priority Min-Queue.
type DisjointSet[E any] interface {

	// MakeSet Adds an element as a trivial subset.
	MakeSet(element E) base.Handle

	// Union merges the two sets containing the given elements
	Union(first base.Handle, second base.Handle)

	// Find returns the unique representative for the containing set
	Find(handle base.Handle) E
}
