package sp

import (
	"errors"

	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// BellmanFordMoore computes shortest distances from any of the given source
// vertices to all vertices in the graph (multi-source shortest paths).
//
// This is a generalization of the single-source shortest-path problem: the
// distance to each vertex is the minimum over all sources. It is NOT a
// solution to the all-pairs shortest-path problem.
//
// The implementation follows the Bellman-Ford-Moore / SPFA approach: sources
// are all enqueued with distance 0 (equivalent to adding a virtual source
// with zero-cost arcs to each real source), and vertices are re-enqueued
// whenever a shorter path is found.
//
// The graph may be mixed (arcs and/or edges). Pass nil for arcCosts if the
// graph has no arcs, and nil for edgeCosts if it has no edges. Arc costs may
// be negative; edge costs must be non-negative. The graph must not contain
// negative-cost cycles.
//
// Returns:
//   - dist:        distance from the nearest source to each vertex
//     (MaxValue[T] for unreachable vertices)
//   - predecessor: the incoming Link (arc or directed edge) on the shortest
//     path, or nil for source vertices and unreachable vertices
//   - err:         non-nil if a negative cycle is detected; dist and
//     predecessor are disposed and nil in that case
func BellmanFordMoore[T base.Number](g *hmgraph.Graph,
	arcCosts *hmgraph.ArcMap[T], edgeCosts *hmgraph.EdgeMap[T],
	sources []*hmgraph.Vertex) (
	dist *hmgraph.VertexMap[T], predecessor *hmgraph.VertexMap[hmgraph.Link], err error) {
	dist = hmgraph.CreateVertexMap(g, "dist", base.MaxValue[T]())
	predecessor = hmgraph.CreateVertexMap[hmgraph.Link](g, "predecessor", nil)
	inQ := hmgraph.CreateVertexMap(g, "inQ", false)
	defer inQ.Dispose()
	count := hmgraph.CreateVertexMap(g, "count", 0)
	defer count.Dispose()
	var q []*hmgraph.Vertex
	n := g.VertexCount()

	for _, source := range sources {
		q = append(q, source)
		inQ.Set(source, true)
		dist.Set(source, 0)
	}

	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		inQ.Set(v, false)
		// A vertex dequeued n or more times implies a negative cycle:
		// without one, every shortest path visits at most n-1 edges.
		removal := count.Get(v) + 1
		if removal >= n {
			dist.Dispose()
			predecessor.Dispose()
			return nil, nil, errors.New("graph contains negative cycle")
		}
		count.Set(v, removal)

		if arcCosts != nil {
			v.ForOutArcs(func(arc *hmgraph.Arc) {
				u := arc.Target()
				cost := dist.Get(v) + arcCosts.Get(arc)
				if cost < dist.Get(u) {
					dist.Set(u, cost)
					predecessor.Set(u, arc)
					if !inQ.Get(u) {
						q = append(q, u)
						inQ.Set(u, true)
					}
				}
			})
		}

		if edgeCosts != nil {
			v.ForEdges(func(edge *hmgraph.Edge) {
				u := edge.Opposite(v)
				cost := dist.Get(v) + edgeCosts.Get(edge)
				if cost < dist.Get(u) {
					dist.Set(u, cost)
					predecessor.Set(u, edge.AsLink(v))
					if !inQ.Get(u) {
						q = append(q, u)
						inQ.Set(u, true)
					}
				}
			})
		}
	}

	return dist, predecessor, nil
}

// BellmanFordMooreSingle is a convenience wrapper around BellmanFordMoore for
// the single-source case.
func BellmanFordMooreSingle[T base.Number](g *hmgraph.Graph,
	arcCosts *hmgraph.ArcMap[T], edgeCosts *hmgraph.EdgeMap[T],
	source *hmgraph.Vertex) (
	dist *hmgraph.VertexMap[T], predecessor *hmgraph.VertexMap[hmgraph.Link], err error) {
	return BellmanFordMoore(g, arcCosts, edgeCosts, []*hmgraph.Vertex{source})
}
