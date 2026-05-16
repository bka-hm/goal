---
title: 'Minimum Spanning Forest'
weight: 5
summary: The `graph/msf` package provides three algorithms for computing a minimum spanning forest of an undirected graph — Prim-Jarník, Kruskal, and Borůvka.
---

The `graph/msf` package computes a **minimum spanning forest** (MSF) of an
undirected weighted graph. For a connected graph the result is a minimum
spanning tree; for a disconnected graph it is one minimum spanning tree per
connected component.

## Interface

All three implementations share the `MsfAlgo` function type:

```go
type MsfAlgo[T cmp.Ordered] func(
    g    *hmgraph.Graph,
    cost *hmgraph.EdgeMap[T],
) (msf []*hmgraph.Edge)
```

### Contract

* `g` must be non-nil and contain at least one vertex. Passing a nil or empty
  graph panics.
* `g` must be undirected — it must contain no directed arcs. Passing a graph
  with arcs panics.
* `cost` maps every edge to its weight. Costs may be negative.
* The returned slice contains exactly **n − k** edges, where *n* is the number
  of vertices and *k* the number of connected components.
* The MSF is unique when all edge costs are distinct. With ties the result is
  *a* minimum spanning forest, but which one is unspecified.

### Usage

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(4)
ab := vs[0].CreateEdge(vs[1])
bc := vs[1].CreateEdge(vs[2])
cd := vs[2].CreateEdge(vs[3])
ac := vs[0].CreateEdge(vs[2])

cost := hmgraph.CreateEdgeMap(g, "cost", 0)
cost.Set(ab, 1)
cost.Set(bc, 3)
cost.Set(cd, 2)
cost.Set(ac, 4)

tree := msf.KruskalMsf(g, cost) // or PrimMsf / BoruvkaMsf
// tree contains ab, cd, bc (total cost 6)
```

---

## Prim-Jarník

```go
func PrimMsf[T cmp.Ordered](g *hmgraph.Graph, cost *hmgraph.EdgeMap[T]) []*hmgraph.Edge
```

Grows the forest one vertex at a time. All vertices start in a priority queue
keyed with ∞. The vertex with the minimum key is extracted, its cheapest
incoming edge (if any) is added to the forest, and the queue keys of its
neighbours are updated via `DecreaseKey`.

**Complexity:** O((V + E) log V) using a leftist heap with `DecreaseKey`.

---

## Kruskal

```go
func KruskalMsf[T cmp.Ordered](g *hmgraph.Graph, cost *hmgraph.EdgeMap[T]) []*hmgraph.Edge
```

Sorts all edges by cost, then adds each edge to the forest if its endpoints
belong to different components. Components are tracked with a union-find
structure.

**Complexity:** O(E log E) — dominated by the sort. Union-find operations
contribute nearly O(E) with path compression and union by rank.
Since E ≤ V², this is also O(E log V).

---

## Borůvka

```go
func BoruvkaMsf[T cmp.Ordered](g *hmgraph.Graph, cost *hmgraph.EdgeMap[T]) []*hmgraph.Edge
```

Works in rounds. Each round scans all edges to find the cheapest edge leaving
each component, then merges components via union-find. Because each round at
least halves the number of components, there are at most O(log V) rounds.

**Complexity:** O(E · α(V) · log V), where α is the inverse Ackermann
function arising from the union-find operations (O(E · α(V)) Find calls per
round, O(log V) rounds). Since α(V) ≤ 4 for any realistic input it is often
omitted, but it is present. Note also that this implementation rescans all E
original edges every round; a contraction-based variant would eliminate
internal edges after each merge, giving a smaller constant in practice while
keeping the same asymptotic bound.