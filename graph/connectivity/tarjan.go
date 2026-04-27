// Package connectivity provides algorithms for finding connected components in graphs.
// It supports both strongly connected components (SCCs) via Tarjan's algorithm and
// weakly connected components, for graphs containing directed arcs, undirected edges, or both.
package connectivity

import (
	"slices"

	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// ComponentType selects whether ConnectedComponents computes strongly or weakly connected components.
type ComponentType int

const (
	// StronglyConnected groups vertices where every vertex is reachable from every other
	// following arc directions. Undirected edges count in both directions.
	StronglyConnected ComponentType = iota
	// WeaklyConnected groups vertices that are connected when all arcs are treated as undirected.
	WeaklyConnected
)

// ConnectedComponents computes strongly connected components with Tarjan's SCC algorithm in O(n+m) time.
// Returns a map containing a component number for each vertex, i.e. if the graph has 4 strongly connected components,
// nodes are assigned values 0 to 3.
func ConnectedComponents(g *hmgraph.Graph, componentType ComponentType) (result [][]*hmgraph.Vertex) {
	k := 0
	S := make([]*hmgraph.Vertex, 0)

	// Maps for index, low and an "in-stack" marker
	indexMap := hmgraph.CreateVertexMap(g, "index", -1)
	defer indexMap.Dispose()
	lowMap := hmgraph.CreateVertexMap(g, "low", -1)
	defer lowMap.Dispose()
	inStackMap := hmgraph.CreateVertexMap(g, "inStack", false)
	defer inStackMap.Dispose()

	// Recursive implementation using closures
	var sccRecursion func(v *hmgraph.Vertex)
	sccRecursion = func(v *hmgraph.Vertex) {
		indexMap.Set(v, k)
		lowMap.Set(v, k)
		k++
		S = append(S, v)
		inStackMap.Set(v, true)

		var vertexVisitor func(w *hmgraph.Vertex)
		vertexVisitor = func(w *hmgraph.Vertex) {
			if indexMap.Get(w) == -1 {
				sccRecursion(w)
				if lowMap.Get(w) < lowMap.Get(v) {
					lowMap.Set(v, lowMap.Get(w))
				}
			} else if inStackMap.Get(w) {
				if indexMap.Get(w) < lowMap.Get(v) {
					lowMap.Set(v, indexMap.Get(w))
				}
			}
		}

		v.ForOutArcs(func(a *hmgraph.Arc) {
			vertexVisitor(a.Target())
		})

		v.ForEdges(func(e *hmgraph.Edge) {
			vertexVisitor(e.Opposite(v))
		})

		if componentType == WeaklyConnected {
			v.ForInArcs(func(a *hmgraph.Arc) {
				vertexVisitor(a.Source())
			})
		}

		if lowMap.Get(v) == indexMap.Get(v) {
			var component []*hmgraph.Vertex

			// while v is in the stack
			for inStackMap.Get(v) {
				w := S[len(S)-1]
				S = S[:len(S)-1]
				component = append(component, w)
				inStackMap.Set(w, false)
			}
			result = append(result, component)
		}
	}

	// "repeated start"
	g.ForVertices(func(v *hmgraph.Vertex) {
		if indexMap.Get(v) != -1 {
			return
		}
		sccRecursion(v)
	})
	slices.Reverse(result)
	return
}

// IsStronglyConnected reports whether g has exactly one strongly connected component,
// i.e. every vertex is reachable from every other vertex following arc directions.
func IsStronglyConnected(g *hmgraph.Graph) bool {
	return len(ConnectedComponents(g, StronglyConnected)) == 1
}

// IsWeaklyConnected reports whether g has exactly one weakly connected component,
// i.e. every vertex is reachable from every other vertex when arc directions are ignored.
func IsWeaklyConnected(g *hmgraph.Graph) bool {
	return len(ConnectedComponents(g, WeaklyConnected)) == 1
}
