package msf

import (
	"cmp"

	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/collections/pq"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// PrimMsf computes a minimum spanning forest using the Prim-Jarník algorithm.
//
// Vertices are inserted into the priority queue lazily: a vertex is only
// enqueued when a connecting edge is first discovered. The outer loop over
// all vertices restarts the process for each connected component, so
// disconnected graphs are handled without needing a sentinel "infinity" key.
//
// Complexity: O((V + E) log V) using a leftist heap with DecreaseKey.
func PrimMsf[T cmp.Ordered](g *hmgraph.Graph, cost *hmgraph.EdgeMap[T]) (msf []*hmgraph.Edge) {
	if g == nil || g.VertexCount() == 0 {
		panic("graph is empty")
	}
	if g.ArcCount() > 0 {
		panic("graph is not undirected")
	}

	q := pq.NewLeftistHeap[*hmgraph.Vertex, T]()
	handles := hmgraph.CreateVertexMap[base.Handle](g, "handles", nil)
	defer handles.Dispose()
	inEdges := hmgraph.CreateVertexMap[*hmgraph.Edge](g, "inEdges", nil)
	defer inEdges.Dispose()
	visited := hmgraph.CreateVertexMap[bool](g, "visited", false)
	defer visited.Dispose()

	relax := func(edge *hmgraph.Edge, v *hmgraph.Vertex) {
		if visited.Get(v) {
			return
		}
		edgeCost := cost.Get(edge)
		h := handles.Get(v)
		if h == nil {
			handles.Set(v, q.Insert(v, edgeCost))
			inEdges.Set(v, edge)
		} else if edgeCost < q.Key(h) {
			inEdges.Set(v, edge)
			q.DecreaseKey(h, edgeCost)
		}
	}

	g.ForVertices(func(start *hmgraph.Vertex) {
		if visited.Get(start) {
			return
		}
		visited.Set(start, true)
		start.ForEdges(func(edge *hmgraph.Edge) {
			relax(edge, edge.Opposite(start))
		})
		for !q.IsEmpty() {
			u, _ := q.ExtractMin()
			visited.Set(u, true)
			handles.Set(u, nil)
			msf = append(msf, inEdges.Get(u))
			u.ForEdges(func(edge *hmgraph.Edge) {
				relax(edge, edge.Opposite(u))
			})
		}
	})
	return msf
}
