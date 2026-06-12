---
title: 'Shortest Paths'
weight: 4
summary: The `graph/sp` package provides single- and multi-source shortest-path algorithms for mixed graphs, including Bellman-Ford-Moore (negative arcs, multi-source) and an SPFA-style Dijkstra (non-negative weights with optional target early-stop).
---

The `graph/sp` package computes shortest paths in mixed graphs (containing arcs, edges,
or both). Arc costs may be negative as long as the graph contains no negative-weight
cycles; edge costs must always be non-negative (a negative undirected edge would form a
trivial negative-weight cycle).

Both algorithms return a distance map and a predecessor map that together encode the
shortest-path tree. Unreachable vertices hold `MaxValue[T]()` in `dist` and `nil` in
`predecessor`.

### Reconstructing a path

```go
var path []hmgraph.Link
for v := target; pred.Get(v) != nil; {
    link := pred.Get(v)
    path = append(path, link)
    v = link.Source()
}
// path is in reverse order — reverse it if needed
```

---

## Bellman-Ford-Moore

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

A queue-based implementation of the Bellman-Ford-Moore algorithm. All given sources are
initialized with distance 0 and
enqueued simultaneously, equivalent to a virtual source connected to each real source by
a zero-cost arc. The result is the shortest distance from *any* source to every vertex.

`BellmanFordMooreSingle` is a convenience wrapper for the common single-source case.

**Complexity:** O(V · E) worst case, O(E) on average for typical inputs.

### Parameters

| Parameter   | Description |
|-------------|-------------|
| `g`         | The graph to search. May be mixed (arcs and/or edges). |
| `arcCosts`  | Cost map for directed arcs. Arc costs may be negative. |
| `edgeCosts` | Cost map for undirected edges. Edge costs must be non-negative. |
| `sources`   | One or more start vertices. All are assigned distance 0. |

### Return values

| Value         | Description |
|---------------|-------------|
| `dist`        | Distance from the nearest source to each vertex. |
| `predecessor` | The incoming `Link` on the shortest path. `nil` for source vertices and unreachable vertices. |
| `err`         | Non-nil if a negative cycle is detected. Both `dist` and `predecessor` are disposed and `nil` in that case. |

### Usage

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(4)
vs[0].CreateArc(vs[1])
vs[1].CreateArc(vs[2])
vs[2].CreateArc(vs[3])

costs := hmgraph.CreateArcMap(g, "cost", 0)
// ... set individual arc costs ...

// Single source
dist, pred, err := sp.BellmanFordMooreSingle(g, costs, nil, vs[0])
if err != nil {
    log.Fatal(err) // negative cycle detected
}

// Multiple sources — shortest distance from either vs[0] or vs[3]
dist, pred, err = sp.BellmanFordMoore(g, costs, nil, []*hmgraph.Vertex{vs[0], vs[3]})
```

---

## Dijkstra

```go
func Dijkstra[N base.Number](
    g         *hmgraph.Graph,
    arcCosts  *hmgraph.ArcMap[N],
    edgeCosts *hmgraph.EdgeMap[N],
    source    *hmgraph.Vertex,
    targets   ...*hmgraph.Vertex,
) (dist *hmgraph.VertexMap[N], predecessor *hmgraph.VertexMap[hmgraph.Link], err error)
```

An SPFA-style variant of Dijkstra's algorithm. Vertices are inserted into the priority
queue **lazily** — only when first reached — and may be re-inserted when a shorter path
is found later. This makes the algorithm correct for graphs with negative-weight arcs, at
the cost of losing the classical Dijkstra guarantee that each vertex is settled exactly
once. Negative-weight arcs are detected at runtime via per-vertex relaxation counts; a
vertex relaxed more than |V| times implies a negative cycle.

**Early stopping with targets:** if one or more target vertices are given *and* all arc
costs are non-negative, the algorithm terminates as soon as any target is extracted from
the priority queue — at that point the standard Dijkstra argument guarantees its distance
is optimal. With negative arc costs early stopping is suppressed, since a later negative
arc could still yield a shorter path to the target.

**Complexity:**
- O((V + E) log V) when all arc costs are non-negative — classic Dijkstra behavior.
- O(V · E) worst case with negative arc costs — SPFA worst case.

### Parameters

| Parameter   | Description |
|-------------|-------------|
| `g`         | The graph to search. May be mixed (arcs and/or edges). |
| `arcCosts`  | Cost map for directed arcs. Arc costs may be negative. |
| `edgeCosts` | Cost map for undirected edges. Edge costs must be non-negative. |
| `source`    | The single start vertex, assigned distance 0. |
| `targets`   | Optional. When provided and no arc costs are negative, the algorithm stops as soon as any target is settled. |

### Return values

| Value         | Description |
|---------------|-------------|
| `dist`        | Shortest distance from `source` to each vertex. |
| `predecessor` | The incoming `Link` on the shortest-path tree. `nil` for the source vertex and unreachable vertices. |
| `err`         | Non-nil if a negative cycle is detected. Both `dist` and `predecessor` are disposed and `nil` in that case. |

### Usage

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(4)
a01 := vs[0].CreateArc(vs[1])
a12 := vs[1].CreateArc(vs[2])
a23 := vs[2].CreateArc(vs[3])

costs := hmgraph.CreateArcMap(g, "cost", 0.0)
costs.Set(a01, 2.0)
costs.Set(a12, 1.0)
costs.Set(a23, 3.0)
edgeCosts := hmgraph.CreateEdgeMap(g, "cost", 0.0)

// All distances from vs[0]
dist, pred, err := sp.Dijkstra(g, costs, edgeCosts, vs[0])
// dist.Get(vs[3]) == 6.0

// Stop as soon as vs[2] is settled (requires non-negative arc costs)
dist, pred, err = sp.Dijkstra(g, costs, edgeCosts, vs[0], vs[2])

// Negative arcs are allowed; a negative cycle returns an error
costs.Set(a12, -5.0)
dist, pred, err = sp.Dijkstra(g, costs, edgeCosts, vs[0])
if err != nil {
    log.Fatal(err) // negative cycle detected
}
```

---

## DijkstraSinglePairShortestPath

```go
func DijkstraSinglePairShortestPath[N base.Number](
    g         *hmgraph.Graph,
    arcCosts  *hmgraph.ArcMap[N],
    edgeCosts *hmgraph.EdgeMap[N],
    source    *hmgraph.Vertex,
    target    *hmgraph.Vertex,
) (path []hmgraph.Link, err error)
```

Convenience wrapper around `Dijkstra` for the single-pair case. Runs Dijkstra with
`target` as the early-stop vertex, extracts the path from the predecessor map, then
disposes both internal maps. The caller receives only the path slice.

The path is returned in **forward order** (source → target). Returns `nil` if `target`
is unreachable or `source == target`. Returns an error if a negative cycle is detected.

### Usage

```go
path, err := sp.DijkstraSinglePairShortestPath(g, costs, edgeCosts, vs[0], vs[3])
if err != nil {
    log.Fatal(err)
}
// path[0] is the first link (out of source), path[len-1] is the last link (into target)
for _, link := range path {
    fmt.Println(link.Source(), "→", link.Target())
}
```

---

## Johnson Transformation

```go
func Johnson[N base.Number](
    g        *hmgraph.Graph,
    arcCosts *hmgraph.ArcMap[N],
) (heights *hmgraph.VertexMap[N], transformedArcCosts *hmgraph.ArcMap[N], err error)

func JohnsonInverse[N base.Number](
    g              *hmgraph.Graph,
    source         *hmgraph.Vertex,
    transformedDist *hmgraph.VertexMap[N],
    heights        *hmgraph.VertexMap[N],
) (dist *hmgraph.VertexMap[N])
```

`Johnson` reweights arc costs so that all transformed costs are non-negative, enabling
Dijkstra to be used on graphs with negative arc costs. It runs multi-source
Bellman-Ford-Moore with all vertices as sources (equivalent to adding a virtual
zero-cost source), yielding height values h(v). The transformed cost of each arc is:

> w'(u, v) = w(u, v) + h(u) − h(v) ≥ 0

`JohnsonInverse` recovers original distances from distances computed on the transformed
graph, using the formula:

> d(source, v) = d'(source, v) − h(source) + h(v)

Only directed graphs (no undirected edges) are supported by `Johnson`. Returns an error
if a negative cycle is detected.

### Usage

```go
// 1. Compute heights and transformed costs
heights, transformed, err := sp.Johnson(g, arcCosts)
if err != nil {
    log.Fatal(err) // negative cycle
}

edgeCosts := hmgraph.CreateEdgeMap(g, "cost", 0.0) // no edges in this graph

// 2. Run Dijkstra from each source using transformed costs
transformedDist, pred, err := sp.Dijkstra(g, transformed, edgeCosts, source)
pred.Dispose()

// 3. Recover original distances
dist := sp.JohnsonInverse(g, source, transformedDist, heights)
transformedDist.Dispose()
```

---

## Constraints

- **Edge costs must be non-negative** for both algorithms — a negative undirected edge would form a trivial negative cycle and causes a panic.
- **Arc costs may be negative**, but the graph must not contain negative-weight cycles. Both algorithms return an error on detection and dispose their output maps.
- `g` and all source/target vertices must be non-nil.