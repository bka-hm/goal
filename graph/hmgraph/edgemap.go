package hmgraph

import "fmt"

// An edgeMap is notified by a graph when the edge set
// changes to adapt to the edge reordering.
type edgeMap interface {
	notifyEdgeCreation()
	Label() string
	AsString(e *Edge) string
	notifyEdgeRemoval(index int)
}

// An EdgeMap is a generic map type from edges to data.
// It uses a slice which keeps the data in the same order
// as the graph's edge slice, i.e. must change when the edge
// set is changed.
type EdgeMap[T any] struct {
	graph        *Graph
	data         []T
	defaultValue T
	label        string
}

// Label returns the map's label
func (edgeMap *EdgeMap[T]) Label() string {
	return edgeMap.label
}

// CreateEdgeMap creates a new edge map with a given type and default.
func CreateEdgeMap[T any](graph *Graph, label string, defaultValue T) *EdgeMap[T] {
	edgeMap := EdgeMap[T]{graph: graph, label: label, data: make([]T, graph.EdgeCount()), defaultValue: defaultValue}
	for i := 0; i < graph.EdgeCount(); i++ {
		edgeMap.data[i] = defaultValue
	}
	graph.edgeMaps = append(graph.edgeMaps, &edgeMap)
	return &edgeMap
}

// Dispose de-registers the edge map so it no longer tracks graph changes.
func (edgeMap *EdgeMap[T]) Dispose() {
	// search for the edge map, replace it with the last edge map and cut
	// the edge map slice.
	for i, e := range edgeMap.graph.edgeMaps {
		if edgeMap == e {
			edgeMap.graph.edgeMaps[i] = edgeMap.graph.edgeMaps[len(edgeMap.graph.edgeMaps)-1]
			edgeMap.graph.edgeMaps = edgeMap.graph.edgeMaps[0 : len(edgeMap.graph.edgeMaps)-1]
		}
	}
}

// Get returns the edge's value
func (edgeMap *EdgeMap[T]) Get(edge *Edge) T {
	return edgeMap.data[edge.index]
}

// Set changes the edge's value
func (edgeMap *EdgeMap[T]) Set(edge *Edge, data T) {
	edgeMap.data[edge.index] = data
}

// notifyEdgeCreation is called on edge creation to fix internal structure
func (edgeMap *EdgeMap[T]) notifyEdgeCreation() {
	edgeMap.data = append(edgeMap.data, edgeMap.defaultValue)
}

// notifyEdgeRemoval is called on edge removal to fix internal structure
func (edgeMap *EdgeMap[T]) notifyEdgeRemoval(index int) {
	oldLen := len(edgeMap.data)
	edgeMap.data[index] = edgeMap.data[oldLen-1]
	edgeMap.data = edgeMap.data[:oldLen-1]
}

// AsString returns an edge's value as string (non-generic)
func (edgeMap *EdgeMap[T]) AsString(e *Edge) string {
	return fmt.Sprintf("%v", edgeMap.data[e.index])
}

// Graph returns the graph an EdgeMap belongs to
func (edgeMap *EdgeMap[T]) Graph() *Graph {
	return edgeMap.graph
}
