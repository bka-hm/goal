package acyclic

import (
	"slices"

	graph_errors "gitlab.lrz.de/hm/goal/graph/errors"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// TopologicalOrder computes any topological vertex order with Kahn's Algorithm
// It will return an error if either the graph contains any (undirected) edge or if the arcs form a cycle.
func TopologicalOrder(g *hmgraph.Graph) (sorting []*hmgraph.Vertex, err error) {
	if g.EdgeCount() > 0 {
		return nil, graph_errors.NewContainsEdgeError(g.AnyEdge())
	}

	// instead of using a queue to store elements that can be processed, we store them in the result list
	// and store how many of them have been processed so far.
	processed := 0

	//number of inbound neighbors not yet processed
	blockers := hmgraph.CreateVertexMap(g, "blockers", 0)
	defer blockers.Dispose()
	g.ForVertices(func(v *hmgraph.Vertex) {
		blockers.Set(v, v.InDegree())
		if v.InDegree() == 0 {
			sorting = append(sorting, v)
		}
	})

	// process next unblocked vertex
	for processed < len(sorting) {
		vertex := sorting[processed]
		vertex.ForOutArcs(func(arc *hmgraph.Arc) {
			target := arc.Target()
			blockingCount := blockers.Get(target) - 1
			blockers.Set(target, blockingCount)
			if blockingCount == 0 {
				sorting = append(sorting, target)
			}
		})
		processed++
	}

	if processed < g.VertexCount() {
		blockingArc := hmgraph.CreateVertexMap[*hmgraph.Arc](g, "blockingArc", nil)
		defer blockingArc.Dispose()
		var vertex *hmgraph.Vertex

		// select one blocking in-arc per blocked vertex and a single blocked start node
		g.ForArcs(func(a *hmgraph.Arc) {
			target := a.Target()
			source := a.Source()
			if blockers.Get(source) > 0 && blockers.Get(target) > 0 {
				blockingArc.Set(target, a)
				vertex = target
			}
		})

		// build the path from the start until we hit the first vertex for second time.
		path := make([]*hmgraph.Arc, 0)
		for blockingArc.Get(vertex) != nil {
			arc := blockingArc.Get(vertex)
			path = append(path, arc)
			blockingArc.Set(vertex, nil)
			vertex = arc.Source()
		}

		// remove the arcs not belonging to the cycle and reverse it
		for path[0].Target() != vertex {
			path = path[1:]
		}
		slices.Reverse(path)

		return nil, graph_errors.NewContainsCycleError(path)
	}

	return sorting, nil
}

func IsDag(g *hmgraph.Graph) bool {
	_, err := TopologicalOrder(g)
	return err == nil
}
