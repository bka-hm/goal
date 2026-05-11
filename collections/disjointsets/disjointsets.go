// Package disjointsets provides a Disjoint-Set (union-find) data structure.
package disjointsets

import (
	"gitlab.lrz.de/hm/goal/base"
)

// DisjointSets is an interface for Disjoint-Sets types (providing union/find operations).
type DisjointSets[E any] interface {

	// MakeSet adds an item as a trivial subset.
	MakeSet(element E) base.Handle

	// Union merges the two sets identified by the given handles.
	Union(first, second base.Handle)

	// Find returns the unique representative for the containing set
	Find(handle base.Handle) E
}
