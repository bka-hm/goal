package hmgraph

import "fmt"

// an arcMap is a listener to operations changing the set of arcs.
type arcMap interface {
	notifyArcCreation()
	GetLabel() string
	GetAsString(a *Arc) string
	notifyArcRemoval(index int)
}

// A ArcMap is a generic map type from vertices to data. It uses a slice which keeps the data in the same order as the graph's vertex slice, i.e. must change when the vertex set is changed.
type ArcMap[T any] struct {
	graph        *Graph
	data         []T
	defaultValue T
	label        string
}

// GetLabel returns the map's label.
func (arcmap *ArcMap[T]) GetLabel() string {
	return arcmap.label
}

// CreateArcMap creates a new arc map.
func CreateArcMap[T any](graph *Graph, label string, defaultValue T) *ArcMap[T] {
	arcmap := ArcMap[T]{graph: graph, label: label, data: make([]T, graph.ArcCount()), defaultValue: defaultValue}
	for i := 0; i < graph.ArcCount(); i++ {
		arcmap.data[i] = defaultValue
	}
	graph.arcMaps = append(graph.arcMaps, &arcmap)
	return &arcmap
}

// Dispose de-registers an arc map
func (arcmap *ArcMap[T]) Dispose() {
	for i, e := range arcmap.graph.arcMaps {
		if arcmap == e {
			arcmap.graph.arcMaps[i] = arcmap.graph.arcMaps[len(arcmap.graph.arcMaps)-1]
			arcmap.graph.arcMaps = arcmap.graph.arcMaps[0 : len(arcmap.graph.arcMaps)-1]
		}
	}
}

// Get returns the arc's value
func (arcmap *ArcMap[T]) Get(a *Arc) T {
	return arcmap.data[a.index]
}

// Set changes the arc's value
func (arcmap *ArcMap[T]) Set(a *Arc, data T) {
	arcmap.data[a.index] = data
}

// notifyArcCreation is called on arc creation to extend the map.
func (arcmap *ArcMap[T]) notifyArcCreation() {
	arcmap.data = append(arcmap.data, arcmap.defaultValue)
}

// notifyArcRemoval is called on arc removal
func (arcMap *ArcMap[T]) notifyArcRemoval(index int) {
	oldLen := len(arcMap.data)
	arcMap.data[index] = arcMap.data[oldLen-1]
	arcMap.data = arcMap.data[:oldLen-1]
}

// GetAsString returns an arc's value as string (non-generic).
func (arcmap *ArcMap[T]) GetAsString(a *Arc) string {
	return fmt.Sprintf("%v", arcmap.data[a.index])
}

// GetGraph returns the map's graph.
func (arcmap *ArcMap[T]) GetGraph() *Graph {
	return arcmap.graph
}
