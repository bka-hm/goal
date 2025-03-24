package hmgraph

import "fmt"

// An Arc is an ordered pair of vertices. It directly implements
// the Link interface as it has a natural order.
type Arc struct {
	// connected vertices
	source *Vertex
	target *Vertex

	// index in the graph's arcs slice
	index       int
	sourceIndex int
	targetIndex int
}

// An ArcVisitor is any function to apply to a set of arcs
type ArcVisitor func(*Arc)

// Source returns the arc's source
func (arc *Arc) Source() *Vertex {
	return arc.source
}

// Target returns the arc's target
func (arc *Arc) Target() *Vertex {
	return arc.target
}

// GetIndex returns the arc's index (which might change when the graph is changed)
func (arc *Arc) GetIndex() int {
	return arc.index
}

// Formatter method, includes data of all maps registered with a graph
func (arc *Arc) String() string {
	result := fmt.Sprintf("Arc #%d-->#%d",
		arc.source.index, arc.target.index)
	graph := arc.source.graph
	if len(graph.arcMaps) > 0 {
		result += " ("
		for _, em := range graph.arcMaps {
			result += fmt.Sprint(em.GetLabel() + "=" + em.GetAsString(arc) + ";")
		}
		result += ")"
	}
	return result
}

type arcLink struct {
	arc *Arc
}

func (arc *Arc) AsReverseLink() Link {
	return &arcLink{arc: arc}
}

func (arcLink *arcLink) Source() *Vertex {
	return arcLink.arc.target
}

func (arcLink *arcLink) Target() *Vertex {
	return arcLink.arc.source
}
