package hmgraph

import (
	"fmt"
)

// A vertexMap is notified by a graph when the vertex set changes to adapt to the vertex reordering.
type vertexMap interface {
	notifyVertexCreation()
	AsString(v *Vertex) string
	Label() string
	notifyVertexRemoval(index int)
}

// A VertexMap is a generic map type from vertices to data. It uses a slice which keeps the data in the same order as the graph's vertex slice, i.e. must change when the vertex set is changed.
type VertexMap[T any] struct {
	graph        *Graph
	label        string
	data         []T
	defaultValue T
}

// Label returns the map's label
func (vmap *VertexMap[T]) Label() string {
	return vmap.label
}

// CreateVertexMap creates a new vertex map
func CreateVertexMap[T any](graph *Graph, label string, defaultValue T) *VertexMap[T] {
	vmap := VertexMap[T]{graph: graph, data: make([]T, graph.VertexCount()), label: label, defaultValue: defaultValue}
	for i := 0; i < graph.VertexCount(); i++ {
		vmap.data[i] = defaultValue
	}
	graph.vertexMaps = append(graph.vertexMaps, &vmap)
	return &vmap
}

// Dispose de-registers the vertex map so it no longer tracks graph changes.
func (vmap *VertexMap[T]) Dispose() {
	for i, e := range vmap.graph.vertexMaps {
		if vmap == e {
			vmap.graph.vertexMaps[i] = vmap.graph.vertexMaps[len(vmap.graph.vertexMaps)-1]
			vmap.graph.vertexMaps = vmap.graph.vertexMaps[0 : len(vmap.graph.vertexMaps)-1]
		}
	}
}

// Get returns the vertex's value.
func (vmap *VertexMap[T]) Get(v *Vertex) T {
	return vmap.data[v.index]
}

// Set changes the vertex's value.
func (vmap *VertexMap[T]) Set(v *Vertex, data T) {
	vmap.data[v.index] = data
}

// notifyVertexCreation fixes the internal representation after vertex creation.
func (vmap *VertexMap[T]) notifyVertexCreation() {
	vmap.data = append(vmap.data, vmap.defaultValue)
}

// notifyVertexRemoval fixes the internal representation after vertex removal.
func (vmap *VertexMap[T]) notifyVertexRemoval(index int) {
	vmap.data[index] = vmap.data[len(vmap.data)-1]
	vmap.data = vmap.data[:len(vmap.data)-1]
}

// AsString returns a vertex's value as string (non-generic)
func (vmap *VertexMap[T]) AsString(v *Vertex) string {
	return fmt.Sprintf("%v", vmap.data[v.index])
}

// Graph returns the graph the vertex belongs to.
func (vmap *VertexMap[T]) Graph() *Graph {
	return vmap.graph
}
