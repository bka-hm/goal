package hmgraph

import (
	"fmt"
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
	g := Graph{
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
	arc.index = len(vertex.graph.arcs) - 1
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
	edge.index = len(vertex.graph.edges) - 1

	// notify arc maps, so they can update to reflect the insertion
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

// GetVertices returns a slice containing all vertices without giving access
// to the underlying structures in O(n).
// Please note that the preferred way of iterating over all vertices
// is the ForVertices function.
func (graph *Graph) GetVertices() (vertices []*Vertex) {
	graph.ForVertices(func(vertex *Vertex) {
		vertices = append(vertices, vertex)
	})
	return
}

// ForArcs runs a function on all arcs in O(m) time (m=number of arcs)
func (graph *Graph) ForArcs(visitor ArcVisitor) {
	for _, arc := range graph.arcs {
		visitor(arc)
	}
}

// ForEdges calls a visitor on all edges in O(m) time
func (graph *Graph) ForEdges(visitor EdgeVisitor) {
	for _, edge := range graph.edges {
		visitor(edge)
	}
}

// String converts a graph with all attached data in human-readable form
func (graph *Graph) String() string {
	result := "Graph\n{"
	for _, v := range graph.vertices {
		result += fmt.Sprintf("\n  %v", v)
	}
	return result + "\n}"
}

// ContainsEdgesOrArcs tests if the graph contains any edges or arc. Returns false if there aren't any edges or arcs.
func (graph *Graph) ContainsEdgesOrArcs() bool {
	return graph.ArcCount()+graph.EdgeCount() != 0
}

// GetFirstVertex returns the first vertex of the graph
func (graph *Graph) GetAnyVertex() *Vertex {
	if len(graph.vertices) == 0 {
		panic("graph is empty")
	}
	return graph.vertices[0]
}

// ContainsLoops returns whether the graph contains a loop in O(m), not optimal
func (graph *Graph) ContainsLoops() bool {
	containsLoops := false
	graph.ForArcs(func(arc *Arc) {
		if arc.Target() == arc.Source() {
			containsLoops = true
		}
	})
	return containsLoops
}

func (graph *Graph) RemoveEdge(edge *Edge) {
	index := edge.index
	edgeCount := graph.EdgeCount()
	graph.edges[index] = graph.edges[edgeCount-1]
	graph.edges[index].index = index
	graph.edges = graph.edges[:edgeCount-1]

	// remove edge from first end's edge slice
	a := edge.vertices[0]
	aIndex := edge.vertexIndices[0]
	a.edges[aIndex] = a.edges[len(a.edges)-1]
	a.edges[aIndex].vertexIndices[0] = aIndex
	a.edges = a.edges[:len(a.edges)-1]

	// remove edge from second end's edge slice
	b := edge.vertices[1]
	bIndex := edge.vertexIndices[1]
	b.edges[bIndex] = b.edges[len(b.edges)-1]
	b.edges[bIndex].vertexIndices[1] = bIndex
	b.edges = b.edges[:len(b.edges)-1]

	// notify edge maps to remove the edge's data
	for _, el := range graph.edgeMaps {
		el.notifyEdgeRemoval(index)
	}
}

func (graph *Graph) RemoveArc(arc *Arc) {
	index := arc.index
	arcCount := graph.ArcCount()
	graph.arcs[index] = graph.arcs[arcCount-1]
	graph.arcs[index].index = index
	graph.arcs = graph.arcs[:arcCount-1]

	// remove edge from first end's edge slice
	source := arc.Source()
	sourceIndex := arc.sourceIndex
	source.outArcs[sourceIndex] = source.outArcs[len(source.outArcs)-1]
	source.outArcs[sourceIndex].sourceIndex = sourceIndex
	source.outArcs = source.outArcs[:len(source.outArcs)-1]

	// remove edge from target arc slice
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

func (graph *Graph) RemoveVertex(vertex *Vertex) {
	vertex.ForEdges(func(edge *Edge) { graph.RemoveEdge(edge) })
	vertex.ForOutArcs(func(arc *Arc) { graph.RemoveArc(arc) })
	vertex.ForInArcs(func(arc *Arc) { graph.RemoveArc(arc) })

	index := vertex.index
	graph.vertices[index] = graph.vertices[len(graph.vertices)-1]
	graph.vertices[index].index = index
	graph.vertices = graph.vertices[:len(graph.vertices)-1]

	// notify verte maps to remove the edge's data
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
