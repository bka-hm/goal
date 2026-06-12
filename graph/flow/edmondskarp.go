// Package flow provides maximum-flow algorithms for directed graphs.
package flow

import (
	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
	"golang.org/x/exp/constraints"
)

// EdmondsKarp computes a maximum flow from source to target using repeated
// breadth-first search for augmenting paths in the residual graph (the
// Edmonds-Karp variant of Ford-Fulkerson).
//
// If scaling is true, capacity scaling is applied: augmenting paths are
// searched in decreasing "phases", each restricted to residual arcs whose
// remaining capacity is at least a threshold C. C starts at the largest arc
// capacity in the graph and is halved whenever no augmenting path is found
// at the current threshold, until it drops below 1. This bounds the number
// of phases to O(log U), where U is the maximum capacity, giving an overall
// complexity of O(E^2 log U) instead of the O(V E^2) of plain Edmonds-Karp.
// If scaling is false, C is fixed at 1, which degenerates to plain
// Edmonds-Karp (every positive-residual arc is eligible in every phase).
//
// g must contain only arcs (no undirected edges); EdmondsKarp panics
// otherwise. Capacities and the resulting flow are non-negative integers.
//
// Besides the flow, EdmondsKarp also returns the source side of a minimum
// s-t cut: the set of vertices still reachable from source in the residual
// graph once no further augmenting path exists. By the max-flow min-cut
// theorem, the arcs leaving this set form a minimum cut whose capacity
// equals the value of the maximum flow.
func EdmondsKarp[T constraints.Integer](g *hmgraph.Graph, kapa *hmgraph.ArcMap[T],
	source *hmgraph.Vertex, target *hmgraph.Vertex, scaling bool) (flow *hmgraph.ArcMap[T], cut []*hmgraph.Vertex) {

	if g.EdgeCount() > 0 {
		panic("graph contains edges")
	}

	// initialize primal solution
	flow = hmgraph.CreateArcMap[T](g, "flow", 0)

	// pred records, for every vertex reached during the current BFS, the arc
	// it was reached over (so the augmenting path can be reconstructed).
	pred := hmgraph.CreateVertexMap[*hmgraph.Arc](g, "predecessor", nil)
	defer pred.Dispose()
	// raise doubles as both the "visited" marker (0 means unvisited) and the
	// bottleneck residual capacity of the path discovered so far from source
	// to that vertex.
	raise := hmgraph.CreateVertexMap[T](g, "raise", 0)
	defer raise.Dispose()

	// initialize C (maximum capacity); with scaling disabled, C stays at 1
	// so every positive-residual arc qualifies (plain Edmonds-Karp).
	C := T(1)
	if scaling {
		g.ForArcs(func(arc *hmgraph.Arc) {
			C = max(C, kapa.Get(arc))
		})
	}

	toVisit := []*hmgraph.Vertex{}

	// Each iteration of this loop is one capacity-scaling phase: it
	// repeatedly augments along C-restricted shortest paths until none is
	// left, then halves C. The loop ends once no augmenting path exists
	// even for C == 1, i.e. the true residual graph is exhausted.
	for C >= 1 {
		g.ForVertices(func(v *hmgraph.Vertex) {
			raise.Set(v, 0)
		})
		raise.Set(source, base.MaxValue[T]())
		toVisit = []*hmgraph.Vertex{source}

		// find a path from source to target
		for len(toVisit) > 0 && raise.Get(target) == 0 {
			// we use a queue (BFS) in Edmonds-Karp.
			v := toVisit[0]
			toVisit = toVisit[1:]

			// forward residual arcs: usable capacity is kapa - flow.
			v.ForOutArcs(func(arc *hmgraph.Arc) {
				opp := arc.Target()
				if raise.Get(opp) == 0 {
					residualKapa := kapa.Get(arc) - flow.Get(arc)
					if residualKapa >= C {
						toVisit = append(toVisit, opp)
						raise.Set(opp, min(raise.Get(v), residualKapa))
						pred.Set(opp, arc)
					}
				}
			})

			// backward residual arcs: walking an arc against its direction
			// lets us cancel existing flow, with usable capacity equal to
			// the flow currently on it.
			v.ForInArcs(func(arc *hmgraph.Arc) {
				opp := arc.Source()
				if raise.Get(opp) == 0 {
					residualKapa := flow.Get(arc)
					if residualKapa >= C {
						toVisit = append(toVisit, opp)
						raise.Set(opp, min(raise.Get(v), residualKapa))
						pred.Set(opp, arc)
					}
				}
			})
		}

		increment := raise.Get(target)
		if increment > 0 {
			// an augmenting path was found; push `increment` units of flow
			// along it by walking predecessors back from target to source,
			// adding flow on forward arcs and cancelling it on backward arcs.
			for current := target; current != source; {
				arc := pred.Get(current)
				if current == arc.Target() {
					flow.Set(arc, flow.Get(arc)+increment)
					current = arc.Source()
				} else {
					flow.Set(arc, flow.Get(arc)-increment)
					current = arc.Target()
				}
			}
		} else {
			// no augmenting path at the current threshold: move to the next
			// (finer) scaling phase.
			C /= 2
		}
	}

	// the final BFS (the one that found no further augmenting path) leaves
	// raise > 0 exactly on the vertices reachable from source in the
	// residual graph — the source side of a minimum cut.
	g.ForVertices(func(v *hmgraph.Vertex) {
		if raise.Get(v) > 0 {
			cut = append(cut, v)
		}
	})

	return
}

// MultiTerminalMaxFlow computes a maximum flow on a graph with several
// sources and sinks instead of a single source/target pair. Each vertex v
// is annotated in support with a supply/demand value:
//
//   - support.Get(v) > 0: v is a source supplying that many flow units.
//   - support.Get(v) < 0: v is a sink demanding that many flow units.
//   - support.Get(v) == 0: v is a transit vertex (neither source nor sink).
//
// Internally this is reduced to the single-source/single-target case: a
// super-source and a super-target are added temporarily, connected to every
// source/sink by an arc whose capacity equals its supply/demand, and
// EdmondsKarp is run between them. The super-vertices (and the arcs to/from
// them) are removed again before returning, so g and kapa are left as they
// were except for the computed flow on the original arcs.
//
// The returned cut is the source side of a minimum cut among the original
// vertices, with the super-source and super-target filtered out.
func MultiTerminalMaxFlow[T constraints.Integer](g *hmgraph.Graph, kapa *hmgraph.ArcMap[T],
	support *hmgraph.VertexMap[T], scaling bool) (*hmgraph.ArcMap[T], []*hmgraph.Vertex) {

	vs := g.CreateVertices(2) // create Super-Source and Super-Target
	source, target := vs[0], vs[1]
	defer g.RemoveVertex(source)
	defer g.RemoveVertex(target)
	support.Set(source, 0)
	support.Set(target, 0)

	// connect the super-source to every supplying vertex, and every
	// demanding vertex to the super-target, with the supply/demand as
	// capacity.
	g.ForVertices(func(v *hmgraph.Vertex) {
		newKapa := support.Get(v)
		if newKapa > 0 {
			arc := source.CreateArc(v)
			kapa.Set(arc, newKapa)
		} else if newKapa < 0 {
			arc := v.CreateArc(target)
			kapa.Set(arc, -newKapa)
		}
	})

	flow, cut := EdmondsKarp(g, kapa, source, target, scaling)
	// remove source and target from cut if present
	filteredCut := []*hmgraph.Vertex{}
	for _, v := range cut {
		if v != source && v != target {
			filteredCut = append(filteredCut, v)
		}
	}
	return flow, filteredCut
}
