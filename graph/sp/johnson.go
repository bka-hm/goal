package sp

import (
	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// Johnson computes Johnson's reweighting for the given directed graph.
//
// It runs multi-source Bellman-Ford-Moore with all vertices as sources (equivalent
// to a virtual zero-cost source connected to every vertex), producing height values
// h(v) such that the transformed arc cost w'(u,v) = w(u,v) + h(u) − h(v) is
// non-negative for every arc. The transformed costs can then be used with Dijkstra.
//
// Only directed graphs (no undirected edges) are supported; a graph with edges panics.
// Returns an error if the graph contains a negative cycle.
func Johnson[N base.Number](g *hmgraph.Graph,
	arcCosts *hmgraph.ArcMap[N]) (heights *hmgraph.VertexMap[N], transformedArcCosts *hmgraph.ArcMap[N], err error) {

	if g.EdgeCount() > 0 {
		panic("Johnson-Transformation not implemented for graphs containing edges.")
	}

	var pred *hmgraph.VertexMap[hmgraph.Link]
	heights, pred, err = BellmanFordMoore(g, arcCosts, nil, g.Vertices())
	if err != nil {
		return nil, nil, err
	}
	pred.Dispose()

	transformedArcCosts = hmgraph.CreateArcMap[N](g, "transformedArcCosts", 0)
	g.ForArcs(func(a *hmgraph.Arc) {
		transformedArcCosts.Set(a, arcCosts.Get(a)+heights.Get(a.Source())-heights.Get(a.Target()))
	})

	return
}

// JohnsonInverse recovers original shortest-path distances from distances computed
// on the Johnson-transformed graph. Given transformed distances d'(source, v) and
// height values h from Johnson, the original distance is:
//
//	d(source, v) = d'(source, v) − h(source) + h(v)
//
// The returned map is owned by the caller and must be disposed when no longer needed.
func JohnsonInverse[N base.Number](g *hmgraph.Graph,
	source *hmgraph.Vertex, transformedDist *hmgraph.VertexMap[N], heights *hmgraph.VertexMap[N]) (dist *hmgraph.VertexMap[N]) {

	sourceHeight := heights.Get(source)
	dist = hmgraph.CreateVertexMap[N](g, "distance", base.MaxValue[N]())
	g.ForVertices(func(v *hmgraph.Vertex) {
		dist.Set(v, transformedDist.Get(v)-sourceHeight+heights.Get(v))
	})

	return
}
