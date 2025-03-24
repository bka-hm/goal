package spanning

import (
	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// A MsfAlgo is an algorithm finding a spanning forest for a given undirected graph
// that minimizes the total cost.
// Panics if the graph is empty (no vertices) or contains any arcs.
type MsfAlgo[T base.Number] func(g *hmgraph.Graph, cost *hmgraph.EdgeMap[T]) (mst []*hmgraph.Edge)
