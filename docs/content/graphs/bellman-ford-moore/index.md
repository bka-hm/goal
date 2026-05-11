---
title: 'Bellman-Ford-Moore'
weight: 4
summary: The `graph/sp` package implements multi-source shortest paths via the Bellman-Ford-Moore algorithm, supporting mixed graphs and negative arc costs.
---

The `graph/sp` package provides shortest-path computation for mixed graphs (containing
arcs, edges, or both). Costs may be negative for arcs, but the graph must not contain
negative-cost cycles.

## Algorithm

`BellmanFordMoore` is a queue-based implementation of the Bellman-Ford-Moore algorithm
(also known as SPFA — Shortest Path Faster Algorithm). It is a generalization of the
classical single-source shortest-path problem: all given sources are initialized with
distance 0 and enqueued simultaneously, which is equivalent to adding a virtual source
vertex connected to each real source with a zero-cost arc. The result is the shortest
distance from *any* source to every vertex.

**Complexity:** O(n · m) worst case, O(m) on average for typical inputs.

## Functions

```go
func BellmanFordMoore[T base.Number](
    g         *hmgraph.Graph,
    arcCosts  *hmgraph.ArcMap[T],
    edgeCosts *hmgraph.EdgeMap[T],
    sources   []*hmgraph.Vertex,
) (dist *hmgraph.VertexMap[T], predecessor *hmgraph.VertexMap[hmgraph.Link], err error)

func BellmanFordMooreSingle[T base.Number](
    g         *hmgraph.Graph,
    arcCosts  *hmgraph.ArcMap[T],
    edgeCosts *hmgraph.EdgeMap[T],
    source    *hmgraph.Vertex,
) (dist *hmgraph.VertexMap[T], predecessor *hmgraph.VertexMap[hmgraph.Link], err error)
```

`BellmanFordMooreSingle` is a convenience wrapper for the common single-source case.

### Parameters

| Parameter   | Description |
|-------------|-------------|
| `g`         | The graph to search. May be mixed (arcs and/or edges). |
| `arcCosts`  | Cost map for directed arcs. Pass `nil` if the graph has no arcs. |
| `edgeCosts` | Cost map for undirected edges. Pass `nil` if the graph has no edges. Edge costs must be non-negative. |
| `sources`   | One or more start vertices. All are assigned distance 0. |

### Return values

| Value         | Description |
|---------------|-------------|
| `dist`        | Distance from the nearest source to each vertex. Unreachable vertices hold `MaxValue[T]()`. |
| `predecessor` | The incoming `Link` (arc or directed edge) on the shortest path. `nil` for source vertices and unreachable vertices. |
| `err`         | Non-nil if a negative cycle is detected. Both `dist` and `predecessor` are disposed and `nil` in that case. |

## Usage

### Single source

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(4)
vs[0].CreateArc(vs[1])
vs[1].CreateArc(vs[2])
vs[2].CreateArc(vs[3])

costs := hmgraph.CreateArcMap(g, "cost", 0)
// ... set individual arc costs ...

dist, pred, err := sp.BellmanFordMooreSingle(g, costs, nil, vs[0])
if err != nil {
    log.Fatal(err)
}
```

### Multiple sources

```go
// Shortest distance from either vs[0] or vs[3] to every other vertex.
dist, pred, err := sp.BellmanFordMoore(g, costs, nil, []*hmgraph.Vertex{vs[0], vs[3]})
```

### Reconstructing a path

The `predecessor` map records the incoming `Link` on the shortest-path tree.
Walk it backwards from any target to recover the full path:

```go
var path []hmgraph.Link
for v := target; pred.Get(v) != nil; {
    link := pred.Get(v)
    path = append(path, link)
    v = link.Source()
}
// path is in reverse order — reverse it if needed
```

### Negative cycle detection

```go
dist, pred, err := sp.BellmanFordMoore(g, costs, nil, sources)
if err != nil {
    fmt.Println("negative cycle detected:", err)
}
```

## Constraints

- Arc costs may be negative; **edge costs must be non-negative** (an undirected negative edge would form a trivial negative cycle).
- The graph must not contain negative-cost cycles. If one exists, the function returns an error and disposes both output maps.
- Pass `nil` for `arcCosts` when the graph has no arcs, and `nil` for `edgeCosts` when it has no edges.