package hmgraph

// A Link is an arc or an edge with a direction, and can hence be a
// part of a (mixed) path
type Link interface {
	Source() *Vertex
	Target() *Vertex
}

// LinkIsArc reports whether a link is a forward arc.
func LinkIsArc(link Link) bool {
	_, result := link.(*Arc)
	return result
}

// LinkIsEdge reports whether a link is an edge (with a chosen direction).
func LinkIsEdge(link Link) bool {
	_, result := link.(*edgeLink)
	return result
}

// LinkIsReverseArc reports whether a link is a reverse arc.
func LinkIsReverseArc(link Link) bool {
	_, result := link.(*reverseArcLink)
	return result
}

// LinkGetArc returns the underlying arc. Panics if the link is not a forward arc.
func LinkGetArc(link Link) *Arc {
	arc, result := link.(*Arc)
	if !result {
		panic("link is not an arc")
	}
	return arc
}

// LinkGetEdge returns the underlying edge. Panics if the link is not an edge.
func LinkGetEdge(link Link) *Edge {
	edgeLink, result := link.(*edgeLink)
	if !result {
		panic("link is not an edge")
	}
	return edgeLink.edge
}

// ReverseLinkGetArc returns the underlying arc of a reverse link. Panics if the link is not a reverse arc.
func ReverseLinkGetArc(link Link) *Arc {
	reverseArcLink, result := link.(*reverseArcLink)
	if !result {
		panic("link is not a reverse arc")
	}
	return reverseArcLink.arc
}
