package msf

import (
	"cmp"

	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/collections/disjointsets"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// BoruvkaMsf computes a minimum spanning forest using Borůvka's algorithm.
//
// Each round scans all edges to find the cheapest edge leaving each component,
// then merges components via union-find. Because each round at least halves
// the number of components, the loop runs at most O(log V) times.
//
// Complexity: O(E · α(V) · log V), where α is the inverse Ackermann function
// from the union-find operations. Unlike a contraction-based implementation,
// this version rescans all original edges every round (rather than a shrinking
// contracted graph), so the constant factor is higher in practice.
func BoruvkaMsf[T cmp.Ordered](g *hmgraph.Graph, cost *hmgraph.EdgeMap[T]) (msf []*hmgraph.Edge) {
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

	// one round
	// find all min edges
	finished := false
	for !finished {
		finished = true
		// minEdges maps each component representative to the cheapest edge leaving that component
		minEdges := hmgraph.CreateVertexMap[*hmgraph.Edge](g, "minEdge", nil)
		g.ForEdges(func(edge *hmgraph.Edge) {
			edgeCost := cost.Get(edge)
			for _, v := range edge.Vertices() {
				compU := uf.Find(handles.Get(edge.Opposite(v)))
				compV := uf.Find(handles.Get(v))
				minEdge := minEdges.Get(compV)
				if compU != compV && (minEdge == nil || cost.Get(minEdge) > edgeCost) {
					minEdges.Set(compV, edge)
				}
			}
		})
		// contract them
		g.ForVertices(func(v *hmgraph.Vertex) {
			e := minEdges.Get(v)
			if e != nil {
				handle1 := handles.Get(e.Vertices()[0])
				handle2 := handles.Get(e.Vertices()[1])
				if uf.Find(handle1) != uf.Find(handle2) {
					uf.Union(handle1, handle2)
					msf = append(msf, e)
					finished = false
				}
			}
		})
		minEdges.Dispose()
	}
	return msf
}
