---
title: 'Acyclic'
weight: 2
summary: The `graphs/acyclic` package provides topological ordering and DAG detection for directed graphs.
---

The `graphs/acyclic` package provides algorithms that operate on directed acyclic graphs (DAGs).
All functions require a graph from the `hmgraph` package and expect it to contain only arcs (no
undirected edges).

## Topological Ordering

```go
func TopologicalOrder(g *hmgraph.Graph) (sorting []*hmgraph.Vertex, err error)
```

`TopologicalOrder` computes a topological vertex ordering using **Kahn's Algorithm**.
It returns the vertices in an order such that every arc `u → v` has `u` appearing before `v`.

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(4)
vs[0].CreateArc(vs[1])
vs[0].CreateArc(vs[2])
vs[1].CreateArc(vs[3])
vs[2].CreateArc(vs[3])

sorting, err := acyclic.TopologicalOrder(g)
if err != nil {
    log.Fatal(err)
}
for _, v := range sorting {
    fmt.Println(v)
}
```

The function returns an error if either condition holds:

* The graph contains at least one undirected edge — a `*graph_errors.ContainsEdgeError` is returned.
* The arcs form a cycle — a `*graph_errors.ContainsCycleError` is returned, including the cycle itself.

In both error cases all internally created maps are disposed before returning, so the graph is left
in a clean state.

## DAG Check

```go
func IsDag(g *hmgraph.Graph) bool
```

`IsDag` returns `true` if and only if the graph is a valid directed acyclic graph, i.e. it contains
no undirected edges and its arcs are free of cycles.

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(3)
vs[0].CreateArc(vs[1])
vs[1].CreateArc(vs[2])

fmt.Println(acyclic.IsDag(g)) // true

vs[2].CreateArc(vs[0])        // introduces a cycle
fmt.Println(acyclic.IsDag(g)) // false
```

## Error Types

Both error types are defined in the `graph/errors` package.

### `ContainsEdgeError`

Returned by `TopologicalOrder` when the graph contains at least one undirected edge.
The offending edge is accessible via the `Edge()` method.

```go
sorting, err := acyclic.TopologicalOrder(g)
var edgeErr *graph_errors.ContainsEdgeError
if errors.As(err, &edgeErr) {
    fmt.Println("graph contains edge:", edgeErr.Edge())
}
```

### `ContainsCycleError`

Returned by `TopologicalOrder` when a directed cycle is detected.
The `Cycle()` method returns the sequence of arcs forming one cycle in the graph.

```go
sorting, err := acyclic.TopologicalOrder(g)
var cycleErr *graph_errors.ContainsCycleError
if errors.As(err, &cycleErr) {
    fmt.Println("cycle arcs:", cycleErr.Cycle())
}
```