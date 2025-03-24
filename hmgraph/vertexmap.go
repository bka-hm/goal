package hmgraph

import (
	"fmt"
)

// A VertexListener is notified by a graph when the vertex set changes to adapt to the vertex reordering.
type vertexMap interface {
	notifyVertexCreation()
	GetAsString(v *Vertex) string
	GetLabel() string
	notifyVertexRemoval(index int)
}

// A VertexMap is a generic map type from vertices to data. It uses a slice which keeps the data in the same order as the graph's vertex slice, i.e. must change when the vertex set is changed.
type VertexMap[T any] struct {
	graph        *Graph
	label        string
	data         []T
	defaultValue T
}

// Returns the map's label
func (vmap *VertexMap[T]) GetLabel() string {
	return vmap.label
}

// CreateVertexMap creates a new arc map
func CreateVertexMap[T any](graph *Graph, label string, defaultValue T) *VertexMap[T] {
	vmap := VertexMap[T]{graph: graph, data: make([]T, graph.VertexCount()), label: label, defaultValue: defaultValue}
	for i := 0; i < graph.VertexCount(); i++ {
		vmap.data[i] = defaultValue
	}
	graph.vertexMaps = append(graph.vertexMaps, &vmap)
	return &vmap
}

// Dispose de-registers a vertex map
func (vmap *VertexMap[T]) Dispose() {
	for i, e := range vmap.graph.vertexMaps {
		if vmap == e {
			vmap.graph.vertexMaps[i] = vmap.graph.vertexMaps[len(vmap.graph.vertexMaps)-1]
			vmap.graph.vertexMaps = vmap.graph.vertexMaps[0 : len(vmap.graph.vertexMaps)-1]
		}
	}
}

// Returns the vertex's value
func (vmap *VertexMap[T]) Get(v *Vertex) T {
	return vmap.data[v.index]
}

// Changes the vertex's value
func (vmap *VertexMap[T]) Set(v *Vertex, data T) {
	vmap.data[v.index] = data
}

// Called on vertex creation to extend the ma
func (vmap *VertexMap[T]) notifyVertexCreation() {
	vmap.data = append(vmap.data, vmap.defaultValue)
}

// Called on vertex creation to extend the ma
func (vmap *VertexMap[T]) notifyVertexRemoval(index int) {
	vmap.data[index] = vmap.data[len(vmap.data)-1]
	vmap.data = vmap.data[:len(vmap.data)-1]
}

// Returns a vertex's value as string (non-generic)
func (vmap *VertexMap[T]) GetAsString(v *Vertex) string {
	return fmt.Sprintf("%v", vmap.data[v.index])
}

func (vmap *VertexMap[T]) GetGraph() *Graph {
	return vmap.graph
}
