// Package disjointset provides a Disjoint-Set (union-find) data structure.
package disjointset

import (
	"gitlab.lrz.de/hm/goal/base"
)

// DisjointSet is a Disjoint-Set data structure (providing union/find operations).
type DisjointSets[E any] interface {

	// MakeSet adds an element as a trivial subset.
	MakeSet(element E) base.Handle

	// Union merges the two sets identified by the given handles.
	Union(first base.Handle, second base.Handle)

	// Find returns the unique representative for the containing set
	Find(handle base.Handle) E
}
