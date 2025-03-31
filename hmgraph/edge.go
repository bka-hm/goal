package hmgraph

import "fmt"

// An Edge is an unordered pair of vertices.
type Edge struct {
	// connected vertices
	vertices [2]*Vertex

	// index in the graph's edges slice
	index int

	//
	vertexIndices [2]int
}

// An EdgeVisitor is any function for visiting edges
type EdgeVisitor func(*Edge)

// Opposite returns the opposite vertex of a given end vertex
func (edge *Edge) Opposite(vertex *Vertex) *Vertex {
	if edge.vertices[0] == vertex && edge.vertices[1] != vertex {
		return edge.vertices[1]
	}
	if edge.vertices[0] != vertex && edge.vertices[1] == vertex {
		return edge.vertices[0]
	}

	panic("Opposite called with a vertex not incident to the edge")
}

// GetIndex returns the edge's graph-global index (not stable under graph modifications!)
func (edge *Edge) GetIndex() int {
	return edge.index
}

// Vertices returns the vertices
func (edge *Edge) Vertices() [2]*Vertex {
	return [2]*Vertex{edge.vertices[0], edge.vertices[1]}
}

// IsIncident returns whether the edge is incident to a given vertex
func (edge *Edge) IsIncident(vertex *Vertex) bool {
	return vertex == edge.vertices[0] || vertex == edge.vertices[1]
}

// String returns a format string including edge map information
func (edge *Edge) String() string {
	result := fmt.Sprintf("Edge #%d---#%d",
		edge.vertices[0].index, edge.vertices[1].index)

	if len(edge.vertices[0].graph.edgeMaps) > 0 {
		result += " ("
		for _, em := range edge.vertices[0].graph.edgeMaps {
			result += fmt.Sprint(em.GetLabel() + "=" + em.GetAsString(edge) + ";")
		}
		result += ")"
	}

	return result
}

// An edgeLink is the actual implementation of the Link interface for an edge .
type edgeLink struct {
	edge   *Edge
	source *Vertex
}

// AsLink returns a Link for a given edge
func (edge *Edge) AsLink(source *Vertex) Link {
	return &edgeLink{edge: edge, source: source}
}

// Source returns the starting vertex for the edgeLink
func (el *edgeLink) Source() *Vertex {
	return el.source
}

// Target returns the target vertex for the edgeLink
func (el *edgeLink) Target() *Vertex {
	return el.edge.Opposite(el.source)
}
