package sp

import (
	"errors"
	"slices"

	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/collections/pq"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// Dijkstra computes single-source shortest paths from source using an SPFA-style
// (Shortest Path Faster Algorithm) variant of Dijkstra's algorithm.
//
// Vertices are inserted into the priority queue lazily (only when first reached)
// and may be re-inserted if a shorter path is found later. This makes the algorithm
// correct for graphs with negative-weight arcs, as long as there are no negative-weight
// cycles (detected via per-vertex extraction counts; if one is found, err is non-nil
// and both dist and predecessor are disposed and nil) and no negative-weight edges
// (undirected edges with negative cost would trivially form negative cycles and are
// rejected upfront with a panic).
//
// If one or more targets are given and all arc costs are non-negative, the algorithm
// stops as soon as the first target is extracted from the priority queue — at that point
// its distance is guaranteed optimal by the standard Dijkstra argument. With negative
// arc costs this early stop is suppressed, since a later negative arc could still yield
// a shorter path to the target.
func Dijkstra[N base.Number](g *hmgraph.Graph,
	arcCosts *hmgraph.ArcMap[N], edgeCosts *hmgraph.EdgeMap[N], source *hmgraph.Vertex,
	targets ...*hmgraph.Vertex) (
	dist *hmgraph.VertexMap[N], predecessor *hmgraph.VertexMap[hmgraph.Link], err error) {
	if g == nil {
		panic("graph is nil")
	}
	if source == nil {
		panic("source is nil")
	}

	// store whether the input has negative arcs
	// must be checked in advance to avoid stopping at a target in the presence of negative weights
	hasNegativeArcCosts := false
	g.ForArcs(func(arc *hmgraph.Arc) {
		if arcCosts.Get(arc) < 0 {
			hasNegativeArcCosts = true
		}
	})
	g.ForEdges(func(edge *hmgraph.Edge) {
		if edgeCosts.Get(edge) < 0 {
			panic("costs of edge cannot be negative")
		}
	})

	// store targets as map for O(1) lookup
	targetSet := hmgraph.CreateVertexMap[bool](g, "targets", false)
	defer targetSet.Dispose()
	for _, t := range targets {
		targetSet.Set(t, true)
	}

	// Predecessor is nil if it was not reached yet (stays nil for unreachable vertices)
	predecessor = hmgraph.CreateVertexMap[hmgraph.Link](g, "predecessors", nil)
	dist = hmgraph.CreateVertexMap[N](g, "distance", base.MaxValue[N]())

	// initialize priority queue (all distances are maxValues)
	priorityQueue := pq.NewLeftistHeap[*hmgraph.Vertex, N]()
	// handles also signal whether a vertex is in the PQ or not
	handles := hmgraph.CreateVertexMap[base.Handle](g, "handles", nil)
	defer handles.Dispose()
	extractCounts := hmgraph.CreateVertexMap(g, "extractions", 0)
	defer extractCounts.Dispose()

	handles.Set(source, priorityQueue.Insert(source, 0))
	dist.Set(source, 0)

	processNeighbor := func(neighbor *hmgraph.Vertex, newDistance N, link hmgraph.Link) {
		if newDistance < dist.Get(neighbor) {
			dist.Set(neighbor, newDistance)
			predecessor.Set(neighbor, link)
			if handles.Get(neighbor) == nil {
				handles.Set(neighbor, priorityQueue.Insert(neighbor, newDistance))
			} else {
				priorityQueue.DecreaseKey(handles.Get(neighbor), newDistance)
			}
		}
	}

	// Dijkstra
	for !priorityQueue.IsEmpty() {
		// extract the minimum
		current, distance := priorityQueue.ExtractMin()
		handles.Set(current, nil)

		// a vertex extracted more than |V| times implies a negative cycle
		extractCount := extractCounts.Get(current) + 1
		if extractCount > g.VertexCount() {
			dist.Dispose()
			predecessor.Dispose()
			return nil, nil, errors.New("graph contains negative cycle")
		}
		extractCounts.Set(current, extractCount)

		// stop if any target has been reached *and* all arc weights are non-negative
		if !hasNegativeArcCosts && targetSet.Get(current) {
			break
		}

		current.ForOutArcs(func(arc *hmgraph.Arc) {
			neighbor := arc.Target()
			cost := arcCosts.Get(arc)
			newCosts := distance + cost
			processNeighbor(neighbor, newCosts, arc)
		})

		current.ForEdges(func(edge *hmgraph.Edge) {
			neighbor := edge.Opposite(current)
			cost := edgeCosts.Get(edge)
			newCosts := distance + cost
			processNeighbor(neighbor, newCosts, edge.AsLink(current))
		})

	}
	return dist, predecessor, nil
}

// DijkstraSinglePairShortestPath returns the shortest path from source to target
// as a slice of Links in forward order (source → target). Returns nil if target
// is unreachable or source == target. Returns an error if a negative cycle is
// detected (see Dijkstra for full constraint details).
func DijkstraSinglePairShortestPath[N base.Number](g *hmgraph.Graph,
	arcCosts *hmgraph.ArcMap[N], edgeCosts *hmgraph.EdgeMap[N], source *hmgraph.Vertex,
	target *hmgraph.Vertex) (path []hmgraph.Link, err error) {
	dist, predecessorMap, err := Dijkstra(g, arcCosts, edgeCosts, source, target)
	if err != nil {
		return nil, err
	}
	defer dist.Dispose()
	defer predecessorMap.Dispose()

	pred := predecessorMap.Get(target)
	for pred != nil {
		path = append(path, pred)
		pred = predecessorMap.Get(pred.Source())
	}
	slices.Reverse(path)
	return path, nil
}
