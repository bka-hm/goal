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

// An ArcVisitor is a function to apply to an arcs
type ArcVisitor func(*Arc)

// Source returns the arc's source
func (arc *Arc) Source() *Vertex {
	return arc.source
}

// Target returns the arc's target
func (arc *Arc) Target() *Vertex {
	return arc.target
}

// Index returns the arc's index (which might change when the graph is changed)
func (arc *Arc) Index() int {
	return arc.index
}

// String returns a formatted representation of the arc, including data of all maps registered with a graph
func (arc *Arc) String() string {
	result := fmt.Sprintf("Arc #%d-->#%d",
		arc.source.index, arc.target.index)
	graph := arc.source.graph
	if len(graph.arcMaps) > 0 {
		result += " ("
		for _, em := range graph.arcMaps {
			result += fmt.Sprint(em.Label() + "=" + em.AsString(arc) + ";")
		}
		result += ")"
	}
	return result
}

// reverseArcLink is a Link representing an arc used from target to source.
type reverseArcLink struct {
	arc *Arc
}

// AsReverseLink returns a Link that traverses the arc in reverse direction.
func (arc *Arc) AsReverseLink() Link {
	return &reverseArcLink{arc: arc}
}

// Source returns the target of the underlying arc (reversed direction).
func (reverseArcLink *reverseArcLink) Source() *Vertex {
	return reverseArcLink.arc.target
}

// Target returns the source of the underlying arc (reversed direction).
func (reverseArcLink *reverseArcLink) Target() *Vertex {
	return reverseArcLink.arc.source
}
