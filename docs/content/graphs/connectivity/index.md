---
title: 'Connectivity'
weight: 2
summary: The `graph/connectivity` package finds strongly and weakly connected components using Tarjan's algorithm.
---

The `graph/connectivity` package finds connected components in mixed graphs (containing
arcs, edges, or both). It supports two notions of connectivity, selected via `ComponentType`:

* **Strongly connected** — every vertex is reachable from every other vertex following arc
  directions. Undirected edges count in both directions.
* **Weakly connected** — every vertex is reachable from every other vertex when arc directions
  are ignored (i.e. arcs are treated as undirected).

## Connected Components

```go
func ConnectedComponents(g *hmgraph.Graph, componentType ComponentType) [][]*hmgraph.Vertex
```

`ConnectedComponents` runs **Tarjan's SCC algorithm** in O(n + m) time and returns all
components as a slice of vertex slices, ordered topologically (sources first).

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(4)
vs[0].CreateArc(vs[1])
vs[1].CreateArc(vs[0]) // {0,1} form one SCC
vs[2].CreateArc(vs[3]) // {2} and {3} are separate SCCs

components := connectivity.ConnectedComponents(g, connectivity.StronglyConnected)
// len(components) == 3
```

To find weakly connected components, pass `WeaklyConnected` instead:

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(3)
vs[0].CreateArc(vs[1])
vs[1].CreateArc(vs[2]) // all reachable ignoring direction

components := connectivity.ConnectedComponents(g, connectivity.WeaklyConnected)
// len(components) == 1
```

## Connectivity Checks

```go
func IsStronglyConnected(g *hmgraph.Graph) bool
func IsWeaklyConnected(g *hmgraph.Graph) bool
```

Convenience wrappers that return `true` when the graph has exactly one component of the
respective type.

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(3)
vs[0].CreateArc(vs[1])
vs[1].CreateArc(vs[2])
vs[2].CreateArc(vs[0])

fmt.Println(connectivity.IsStronglyConnected(g)) // true  — cycle covers all vertices
fmt.Println(connectivity.IsWeaklyConnected(g))   // true

g2 := hmgraph.NewGraph()
ws := g2.CreateVertices(3)
ws[0].CreateArc(ws[1])
ws[1].CreateArc(ws[2])

fmt.Println(connectivity.IsStronglyConnected(g2)) // false — no path back from ws[2]
fmt.Println(connectivity.IsWeaklyConnected(g2))   // true  — connected ignoring direction
```