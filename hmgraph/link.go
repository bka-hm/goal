package hmgraph

// A Link is an arc or an edge with a direction, and can hence be a
// part of a (mixed) path
type Link interface {
	Source() *Vertex
	Target() *Vertex
}

func LinkIsArc(link *Link) bool {
	_, result := (*link).(*Arc)
	return result
}

func LinkIsEdge(link *Link) bool {
	_, result := (*link).(*edgeLink)
	return result
}

func LinkIsReverseArc(link *Link) bool {
	_, result := (*link).(*arcLink)
	return result
}

func LinkGetArc(link *Link) *Arc {
	arc, result := (*link).(*Arc)
	if !result {
		panic("link is not an arc")
	}
	return arc
}

func LinkGetEdge(link *Link) *Edge {
	edgeLink, result := (*link).(*edgeLink)
	if !result {
		panic("link is not an edge")
	}
	return edgeLink.edge
}

func ReverseLinkGetArc(link *Link) *Arc {
	arcLink, result := (*link).(*arcLink)
	if !result {
		panic("link is not a reverse arc")
	}
	return arcLink.arc
}
