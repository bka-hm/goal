---
title: 'Maximum Flow'
weight: 6
summary: The `graph/flow` package computes maximum flows and minimum cuts with Edmonds-Karp, with an optional capacity-scaling speedup, plus a multi-terminal (multi-source/multi-sink) wrapper.
---

The `graph/flow` package computes a **maximum flow** between two vertices of a directed
graph, together with a matching **minimum cut**. All functions require a graph from the
`hmgraph` package that contains only arcs (no undirected edges).

---

## Edmonds-Karp

```go
func EdmondsKarp[T constraints.Integer](
    g       *hmgraph.Graph,
    kapa    *hmgraph.ArcMap[T],
    source  *hmgraph.Vertex,
    target  *hmgraph.Vertex,
    scaling bool,
) (flow *hmgraph.ArcMap[T], cut []*hmgraph.Vertex)
```

Computes a maximum flow from `source` to `target` by repeatedly searching for an
augmenting path in the residual graph using breadth-first search — the Edmonds-Karp
variant of Ford-Fulkerson — and pushing as much flow as the path's bottleneck capacity
allows.

**Capacity scaling:** if `scaling` is `true`, the search is restricted in phases. Each
phase only considers residual arcs whose remaining capacity is at least a threshold `C`.
`C` starts at the largest arc capacity in the graph and is halved whenever no augmenting
path remains at the current threshold, until it drops below `1`. This bounds the number
of phases to O(log U), where U is the maximum capacity, giving an overall complexity of
**O(E² log U)** instead of plain Edmonds-Karp's **O(V·E²)**. If `scaling` is `false`, `C`
stays fixed at `1`, which is exactly plain Edmonds-Karp (every positive-residual arc is
eligible immediately).

**Minimum cut:** besides the flow, `EdmondsKarp` returns the source side of a minimum
s-t cut — the vertices still reachable from `source` in the residual graph once no
augmenting path is left. By the max-flow min-cut theorem, the arcs leaving this set form
a minimum cut whose capacity equals the value of the maximum flow.

### Parameters

| Parameter | Description |
|-----------|-------------|
| `g`       | The graph to search. Must contain only arcs; panics if it contains edges. |
| `kapa`    | Arc capacities. Must be non-negative. |
| `source`  | The flow source. |
| `target`  | The flow sink. |
| `scaling` | Whether to use capacity scaling (recommended for large capacities) or plain Edmonds-Karp. |

### Return values

| Value  | Description |
|--------|-------------|
| `flow` | Flow value on each arc, satisfying `0 ≤ flow(arc) ≤ kapa(arc)` and flow conservation at every vertex except `source`/`target`. |
| `cut`  | The vertices on the source side of a minimum cut (always includes `source`). |

### Usage

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(4)
a01 := vs[0].CreateArc(vs[1])
a02 := vs[0].CreateArc(vs[2])
a13 := vs[1].CreateArc(vs[3])
a23 := vs[2].CreateArc(vs[3])

kapa := hmgraph.CreateArcMap(g, "capacity", 0)
kapa.Set(a01, 3)
kapa.Set(a02, 2)
kapa.Set(a13, 2)
kapa.Set(a23, 3)

flow, cut := flowpkg.EdmondsKarp(g, kapa, vs[0], vs[3], true)
// value of the maximum flow:
value := flow.Get(a01) + flow.Get(a02)
```

---

## MultiTerminalMaxFlow

```go
func MultiTerminalMaxFlow[T constraints.Integer](
    g       *hmgraph.Graph,
    kapa    *hmgraph.ArcMap[T],
    support *hmgraph.VertexMap[T],
    scaling bool,
) (flow *hmgraph.ArcMap[T], cut []*hmgraph.Vertex)
```

Computes a maximum flow on a graph with several sources and sinks instead of a single
source/target pair. Each vertex `v` is annotated in `support` with a supply or demand:

* `support.Get(v) > 0` — `v` is a source supplying that many flow units.
* `support.Get(v) < 0` — `v` is a sink demanding that many flow units.
* `support.Get(v) == 0` — `v` is a plain transit vertex.

Internally, this reduces to the single-source/single-target case: a super-source and
super-target are added temporarily, connected to every source/sink by an arc whose
capacity equals its supply/demand, and `EdmondsKarp` is run between them. The
super-vertices and their arcs are removed again before returning, so `g` and `kapa` are
left as they were except for the flow computed on the original arcs.

The returned `cut` is the source side of a minimum cut among the **original** vertices —
the super-source and super-target are filtered out.

### Parameters

| Parameter | Description |
|-----------|-------------|
| `g`       | The graph to search. Must contain only arcs. |
| `kapa`    | Arc capacities. Extended internally with the arcs to/from the super-terminals. |
| `support` | Supply (positive) / demand (negative) per vertex. Zero for transit vertices. |
| `scaling` | Forwarded to `EdmondsKarp`. |

### Usage

```go
g := hmgraph.NewGraph()
vs := g.CreateVertices(2)
arc := vs[0].CreateArc(vs[1])
kapa := hmgraph.CreateArcMap(g, "capacity", 0)
kapa.Set(arc, 8)

support := hmgraph.CreateVertexMap(g, "support", 0)
support.Set(vs[0], 3)  // vs[0] supplies 3 units
support.Set(vs[1], -3) // vs[1] demands 3 units

flow, cut := flowpkg.MultiTerminalMaxFlow(g, kapa, support, true)
// flow.Get(arc) == 3
```

---

## Constraints

- `g` must contain only arcs (no undirected edges); `EdmondsKarp` panics otherwise.
- Capacities, supplies, and the resulting flow are non-negative integer types
  (`constraints.Integer`).
- `MultiTerminalMaxFlow` temporarily mutates `g` and `kapa` (adding the super-terminals
  and their arcs) but restores both before returning.