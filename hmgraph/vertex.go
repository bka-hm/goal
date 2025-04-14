package hmgraph

import (
	"fmt"
	"slices"
)

// A Vertex is a struct containing all incident arcs and edges.
type Vertex struct {
	// incident arcs and edges
	inArcs  []*Arc
	outArcs []*Arc
	edges   []*Edge

	// the graph containing the vertex
	graph *Graph

	// index in the graph's vertex array
	index int
}

// A VertexVisitor is a functions for visiting vertices
type VertexVisitor func(v *Vertex)

// OutDegree returns the number of outgoing arcs
func (vertex *Vertex) OutDegree() int {
	return len(vertex.outArcs)
}

// InDegree returns the number of incoming arcs
func (vertex *Vertex) InDegree() int {
	return len(vertex.inArcs)
}

// Degree returns the number of edges
func (vertex *Vertex) Degree() int {
	return len(vertex.edges)
}

// ForOutArcs calls a visitor on a vertex's outgoing arcs in O(OutDegree(vertex)) time
func (vertex *Vertex) ForOutArcs(visitor ArcVisitor) {
	for _, arc := range vertex.outArcs {
		visitor(arc)
	}
}

// GetOutArcs returns a slice containing all outgoing arcs
// without giving access to the underlying structures in O(outdeg(v)).
// Please note that the preferred way of iterating over all outgoing arcs
// is the ForOutArcs function.
func (vertex *Vertex) GetOutArcs() (outArcs []*Arc) {
	return slices.Clone(vertex.outArcs)
}

// ForInArcs calls a visitor on a vertex's incoming arcs in O(InDegree(vertex)) time
func (vertex *Vertex) ForInArcs(visitor ArcVisitor) {
	for _, arc := range vertex.inArcs {
		visitor(arc)
	}
}

// GetInArcs returns a slice containing all ingoing arcs
// without giving access to the underlying structures in O(indeg(v)).
// Please note that the preferred way of iterating over all ingoing arcs
// is the ForInArcs function.
func (vertex *Vertex) GetInArcs() (inArcs []*Arc) {
	return slices.Clone(vertex.inArcs)
}

// ForEdges calls a visitor on a vertex's edges in O(Degree(vertex)) time
func (vertex *Vertex) ForEdges(visitor EdgeVisitor) {
	for _, edge := range vertex.edges {
		visitor(edge)
	}
}

// GetEdges returns a slice containing all incident edges
// without giving access to the underlying structures in O(deg(v)).
// Please note that the preferred way of iterating over all edges
// is the ForEdges function.
func (vertex *Vertex) GetEdges() (edges []*Edge) {
	return slices.Clone(vertex.edges)
}

// GetIndex returns the vertex's index (not stable under graph modifications!)
func (vertex *Vertex) GetIndex() int {
	return vertex.index
}

// GetGraph returns the graph a vertex belongs to.
func (vertex *Vertex) GetGraph() *Graph {
	return vertex.graph
}

// String returns a formatted representation of a vertex, including all data
// from maps and all arcs/edges.
func (vertex *Vertex) String() string {
	result := fmt.Sprintf("Vertex #%d", vertex.index)

	if len(vertex.graph.vertexMaps) > 0 {
		result += " ("
		for _, vm := range vertex.graph.vertexMaps {
			result += fmt.Sprint(vm.GetLabel() + "=" + vm.GetAsString(vertex) + ";")
		}
		result += ")"
	}

	if vertex.OutDegree()+vertex.InDegree()+vertex.Degree() > 0 {
		result += "\n  {"

		for _, edge := range vertex.edges {
			result += fmt.Sprintf("\n    %v", edge)
		}
		for _, arc := range vertex.outArcs {
			result += fmt.Sprintf("\n    %v", arc)
		}
		for _, arc := range vertex.inArcs {
			result += fmt.Sprintf("\n    %v", arc)
		}
		result += "\n  }"
	}
	return result
}
