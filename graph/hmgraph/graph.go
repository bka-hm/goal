package hmgraph

import (
	"fmt"
	"slices"
)

// A Graph is a struct of vertices and both arcs and
// edges, unifying the notion of directed and undirected
// graphs. A Graph can have arbitrarily many maps associated
// with vertices, arcs and edges.
type Graph struct {
	// Slice containing all vertices, consistent with vertex indices, but subject to reordering
	vertices   []*Vertex
	vertexMaps []vertexMap // this is the interface, we store pointers!

	// Slice containing all arcs, consistent with arc indices, but subject to reordering
	arcs    []*Arc
	arcMaps []arcMap // this is the interface, we store pointers!

	// Slice containing all edges, consistent with edges indices, but subject to reordering
	edges    []*Edge
	edgeMaps []edgeMap // this is the interface, we store pointers!
}

// NewGraph creates a new, empty graph.
func NewGraph() *Graph {
	g := Graph{ // slice initializations are optional
		vertices:   []*Vertex{},
		vertexMaps: []vertexMap{},
		arcs:       []*Arc{},
		arcMaps:    []arcMap{},
		edges:      []*Edge{},
		edgeMaps:   []edgeMap{},
	}
	return &g
}

// VertexCount returns the number of vertices in O(1)
func (graph *Graph) VertexCount() int {
	return len(graph.vertices)
}

// ArcCount returns the number of arcs in O(1)
func (graph *Graph) ArcCount() int {
	return len(graph.arcs)
}

// EdgeCount returns the number of edges in O(1)
func (graph *Graph) EdgeCount() int {
	return len(graph.edges)
}

// CreateVertex creates a vertex in amortized O(1)
func (graph *Graph) CreateVertex() *Vertex {
	vertex := &Vertex{graph: graph, index: len(graph.vertices)}
	graph.vertices = append(graph.vertices, vertex)

	// notify vertex maps, so they can update to reflect the insertion
	for _, vl := range graph.vertexMaps {
		vl.notifyVertexCreation()
	}
	return vertex
}

// CreateArc creates an arc in amortized O(1)
// caveat: may create multi-arcs if arc already exists
func (vertex *Vertex) CreateArc(target *Vertex) *Arc {
	arc := &Arc{source: vertex, target: target, index: len(vertex.graph.arcs)}
	vertex.graph.arcs = append(vertex.graph.arcs, arc)
	vertex.outArcs = append(vertex.outArcs, arc)
	arc.sourceIndex = len(vertex.outArcs) - 1
	target.inArcs = append(target.inArcs, arc)
	arc.targetIndex = len(target.inArcs) - 1

	// notify arc maps, so they can update to reflect the insertion
	for _, al := range target.graph.arcMaps {
		al.notifyArcCreation()
	}
	return arc
}

// CreateArcs creates several arcs at once in O(len(target)).
func (vertex *Vertex) CreateArcs(target []*Vertex) []*Arc {
	arcs := make([]*Arc, len(target))
	for i, al := range target {
		arcs[i] = vertex.CreateArc(al)
	}
	return arcs
}

// CreateEdges creates several edges at once in O(len(target)).
func (vertex *Vertex) CreateEdges(target []*Vertex) []*Edge {
	edges := make([]*Edge, len(target))
	for i, al := range target {
		edges[i] = vertex.CreateEdge(al)
	}
	return edges
}

// CreateEdge creates an edge in amortized O(1)
// caveat: may create multi-edges if edge already exists
// panics when creating a loop.
func (vertex *Vertex) CreateEdge(opposite *Vertex) *Edge {
	if opposite == nil {
		panic("the input is not a vertex (nil or different value)")
	}

	if vertex.index == opposite.index {
		panic("no loop edges allowed")
	}

	edge := &Edge{vertices: [...]*Vertex{vertex, opposite}, index: len(vertex.graph.edges)}
	vertex.edges = append(vertex.edges, edge)
	edge.vertexIndices[0] = len(vertex.edges) - 1
	opposite.edges = append(opposite.edges, edge)
	edge.vertexIndices[1] = len(opposite.edges) - 1
	vertex.graph.edges = append(vertex.graph.edges, edge)

	// notify edge maps, so they can update to reflect the insertion
	for _, el := range opposite.graph.edgeMaps {
		el.notifyEdgeCreation()
	}
	return edge
}

// ForVertices runs a function on all vertices in O(n) time
func (graph *Graph) ForVertices(visitor VertexVisitor) {
	for _, vertex := range graph.vertices {
		visitor(vertex)
	}
}

// Vertices returns a slice containing all vertices without giving access
// to the underlying structures in O(n).
// Please note that the preferred way of iterating over all vertices
// is the ForVertices function.
func (graph *Graph) Vertices() (vertices []*Vertex) {
	return slices.Clone(graph.vertices)
}

// ForArcs runs a function on all arcs in O(m) time (m=number of arcs)
func (graph *Graph) ForArcs(visitor ArcVisitor) {
	for _, arc := range graph.arcs {
		visitor(arc)
	}
}

// Arcs returns a slice containing all arcs without giving access
// to the underlying structures in O(n).
// Please note that the preferred way of iterating over all arcs
// is the ForArcs function.
func (graph *Graph) Arcs() (arcs []*Arc) {
	return slices.Clone(graph.arcs)
}

// ForEdges calls a visitor on all edges in O(m) time
func (graph *Graph) ForEdges(visitor EdgeVisitor) {
	for _, edge := range graph.edges {
		visitor(edge)
	}
}

// Edges returns a slice containing all edges without giving access
// to the underlying structures in O(n).
// Please note that the preferred way of iterating over all edges
// is the ForEdges function.
func (graph *Graph) Edges() (edges []*Edge) {
	return slices.Clone(graph.edges)
}

// String converts a graph with all attached data in human-readable form
func (graph *Graph) String() string {
	result := "Graph\n{"
	for _, v := range graph.vertices {
		result += fmt.Sprintf("\n  %v", v)
	}
	return result + "\n}"
}

// ContainsEdgesOrArcs returns true if the graph contains any edges or arcs.
func (graph *Graph) ContainsEdgesOrArcs() bool {
	return graph.ArcCount()+graph.EdgeCount() != 0
}

// AnyVertex returns an arbitrary vertex of the graph.
func (graph *Graph) AnyVertex() *Vertex {
	if len(graph.vertices) == 0 {
		panic("graph is empty")
	}
	return graph.vertices[0]
}

// AnyEdge returns an arbitrary edge of the graph.
// Panics if the graph contains no edges.
func (graph *Graph) AnyEdge() *Edge {
	if len(graph.edges) == 0 {
		panic("graph contains no edges")
	}
	return graph.edges[0]
}

// AnyArc returns an arbitrary arc of the graph.
// Panics if the graph contains no arcs.
func (graph *Graph) AnyArc() *Arc {
	if len(graph.arcs) == 0 {
		panic("graph contains no arcs")
	}
	return graph.arcs[0]
}

// AnyVertexWhere returns an arbitrary vertex satisfying predicate, or nil if none does.
func (graph *Graph) AnyVertexWhere(predicate func(*Vertex) bool) *Vertex {
	for _, v := range graph.vertices {
		if predicate(v) {
			return v
		}
	}
	return nil
}

// AnyArcWhere returns an arbitrary arc satisfying predicate, or nil if none does.
func (graph *Graph) AnyArcWhere(predicate func(*Arc) bool) *Arc {
	for _, a := range graph.arcs {
		if predicate(a) {
			return a
		}
	}
	return nil
}

// AnyEdgeWhere returns an arbitrary edge satisfying predicate, or nil if none does.
func (graph *Graph) AnyEdgeWhere(predicate func(*Edge) bool) *Edge {
	for _, e := range graph.edges {
		if predicate(e) {
			return e
		}
	}
	return nil
}

// RemoveEdge removes an edge from the graph. It uses the swap-and-shrink
// approach to work in O(1) but with a need to update incident vertices and
// all edge-maps about the change.
func (graph *Graph) RemoveEdge(edge *Edge) {
	// swap the last edge into the removed edge's position
	index := edge.index
	lastIndex := graph.EdgeCount() - 1
	graph.edges[index] = graph.edges[lastIndex]
	graph.edges[index].index = index
	graph.edges = graph.edges[:lastIndex]

	// remove edge from both endpoints' edge slices
	for i, v := range edge.vertices {
		idx := edge.vertexIndices[i]
		last := v.Degree() - 1
		v.edges[idx] = v.edges[last]
		moved := v.edges[idx]
		if moved.vertices[0] == v {
			moved.vertexIndices[0] = idx
		} else {
			moved.vertexIndices[1] = idx
		}
		v.edges = v.edges[:last]
	}

	// notify edge maps to remove the edge's data
	for _, el := range graph.edgeMaps {
		el.notifyEdgeRemoval(index)
	}
}

// RemoveArc removes an arc from the graph using the swap-and-shrink approach in O(1).
func (graph *Graph) RemoveArc(arc *Arc) {
	index := arc.index
	arcCount := graph.ArcCount()
	graph.arcs[index] = graph.arcs[arcCount-1]
	graph.arcs[index].index = index
	graph.arcs = graph.arcs[:arcCount-1]

	// remove arc from source vertex's arc slice
	source := arc.Source()
	sourceIndex := arc.sourceIndex
	source.outArcs[sourceIndex] = source.outArcs[len(source.outArcs)-1]
	source.outArcs[sourceIndex].sourceIndex = sourceIndex
	source.outArcs = source.outArcs[:len(source.outArcs)-1]

	// remove arc from target vertex's arc slice
	target := arc.Target()
	targetIndex := arc.targetIndex
	target.inArcs[targetIndex] = target.inArcs[len(target.inArcs)-1]
	target.inArcs[targetIndex].targetIndex = targetIndex
	target.inArcs = target.inArcs[:len(target.inArcs)-1]

	// notify arc maps to remove the arc's data
	for _, arcmap := range graph.arcMaps {
		arcmap.notifyArcRemoval(index)
	}
}

// RemoveVertex removes a vertex and all its incident edges and arcs from the graph.
func (graph *Graph) RemoveVertex(vertex *Vertex) {
	vertex.ForEdges(func(edge *Edge) { graph.RemoveEdge(edge) })
	vertex.ForOutArcs(func(arc *Arc) { graph.RemoveArc(arc) })
	vertex.ForInArcs(func(arc *Arc) { graph.RemoveArc(arc) })

	index := vertex.index
	graph.vertices[index] = graph.vertices[len(graph.vertices)-1]
	graph.vertices[index].index = index
	graph.vertices = graph.vertices[:len(graph.vertices)-1]

	// notify vertex maps to remove the vertex's data
	for _, vertexMap := range graph.vertexMaps {
		vertexMap.notifyVertexRemoval(index)
	}
}

// CreateVertices creates a given number of vertices
func (graph *Graph) CreateVertices(n int) []*Vertex {
	result := make([]*Vertex, n)
	for i := 0; i < n; i++ {
		result[i] = graph.CreateVertex()
	}
	return result
}

// MapCount returns the number of active maps.
func (graph *Graph) MapCount() (vertexMapCount int, arcMapCount int, edgeMapCount int) {
	return len(graph.vertexMaps), len(graph.arcMaps), len(graph.edgeMaps)
}
