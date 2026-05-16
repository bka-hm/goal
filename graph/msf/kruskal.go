package msf

import (
	"cmp"
	"slices"

	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/collections/disjointsets"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// KruskalMsf computes a minimum spanning forest using Kruskal's algorithm.
//
// Edges are sorted by cost in ascending order. Each edge is added to the
// forest if and only if its two endpoints belong to different components,
// tracked with a union-find structure.
//
// Complexity: O(E log E) dominated by the sort; union-find operations are
// nearly O(E) with path compression and union by rank.
func KruskalMsf[T cmp.Ordered](g *hmgraph.Graph, cost *hmgraph.EdgeMap[T]) (msf []*hmgraph.Edge) {
	if g == nil || g.VertexCount() == 0 {
		panic("graph is empty")
	}
	if g.ArcCount() > 0 {
		panic("graph is not undirected")
	}

	uf := disjointsets.NewUnionFind[*hmgraph.Vertex]()
	handles := hmgraph.CreateVertexMap[base.Handle](g, "handles", nil)
	defer handles.Dispose()

	g.ForVertices(func(v *hmgraph.Vertex) {
		handles.Set(v, uf.MakeSet(v))
	})

	edges := g.Edges()
	slices.SortFunc(edges, func(a, b *hmgraph.Edge) int {
		return cmp.Compare(cost.Get(a), cost.Get(b))
	})

	for _, edge := range edges {
		ev := edge.Vertices()
		h0, h1 := handles.Get(ev[0]), handles.Get(ev[1])
		if uf.Find(h0) != uf.Find(h1) {
			msf = append(msf, edge)
			uf.Union(h0, h1)
		}
	}

	return msf
}
