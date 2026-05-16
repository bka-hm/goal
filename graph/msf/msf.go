package msf

import (
	"cmp"

	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// MsfAlgo is the common signature for minimum spanning forest algorithms.
//
// Given a connected undirected graph it returns a minimum spanning tree;
// given a disconnected graph it returns a minimum spanning forest (one tree
// per connected component). The returned slice contains exactly n-k edges,
// where n is the number of vertices and k the number of connected components.
//
// All implementations panic if g is nil, contains no vertices, or contains
// any directed arcs.
type MsfAlgo[T cmp.Ordered] func(g *hmgraph.Graph, cost *hmgraph.EdgeMap[T]) (mst []*hmgraph.Edge)
