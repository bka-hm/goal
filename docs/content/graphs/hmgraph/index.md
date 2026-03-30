---
title: 'HMGraph'
weight: 1
summary: The `graphs/hmgraph` package contains a fast and versatile data structure for mixed directed/undirected graphs.
---

The `graphs/hmgraph` package contains a data structure for
hybrid/mixed graphs (possibly containing simultaneously both arcs and edges)
with arbitrary many generic O(1)-accessible maps for vertices, edges and arcs.

A *Graph* (`hmgraph.Graph`) is a graph potentially containing both arcs and
edges, i.e. a tuple G=(V, A, E) of a set of vertices V,
a set of (directed) arcs A and a set of (undirected) edges E.

Algorithms may or may not generalize to graphs containing both connector types.

## Creating graphs

Graphs, vertices, arcs and edges can be created via functions  `NewGraph`, `CreateVertex`, `CreateArc` and `CreateEdge`,
respectively:

```go
g := hmgraph.NewGraph()
u := g.CreateVertex()
v := g.CreateVertex()
w := g.CreateVertex()
u.CreateArc(v)
w.CreateEdge(v) // same as v.CreateEdge(w)
```

All operations run in (amortized) constant time (O(1)). Creation of loop edges will result in an
error; loop arcs are allowed as well as parallel edges/arcs (multi-edges/arcs). A given
number of vertices can be created at once using the function ```CreateVertices(n int)```. Creating multiple edges or
arcs is supported via a ```CreateArcs``` and ```CreateEdges```, respectively.

The difference between ```Arc```s and ```Edge```s is that an ```Arc``` connect two Vertices directionally, i.e. they
have a *source* and a *target* vertex whereas an ```Edge``` connects two unordered vertices.

## Conversion to `string`

Graphs support standard conversion to `string`:

```go
fmt.Println("%v", g)
```

will, for the above example, result in

```
Graph
{
  Vertex #0 
  {
    Arc #0-->#1
  }
  Vertex #1 
  {
    Edge #2---#1
    Arc #0-->#1
  }
  Vertex #2 
  {
    Edge #2---#1
  }
}
```

Vertices are denoted with their *iteration index*.
This index may change when vertices are deleted and should not be depended on!

## Deletion of edges, arcs, and vertices

Edges and arcs can be deleted using the `RemoveEdge` and `RemoveArc` in O(1), respectively.
Vertices can be removed with all their incident arcs and edges using the
function `RemoveVertex` in O(1+inDeg+outDeg+Deg).

```go
g := hmgraph.NewGraph()
...
g.RemoveEdge(e)
g.RemoveArc(a)
g.RemoveVertex(v) // removes all incident arcs and edges, too
```

## Maps

Storing additional data for vertices and arcs/edges is possible using generic maps with O(1) access, improving over standard maps. 
They can be created at any time after creation of the graph. 
That is, maps can be created before all
vertices/arcs/edges are created and will work even when vertices/arcs/edges are added or removed.
Each map increases memory consumption by O(|V|), O(|E|) od O(|A|).

Maps are created with functions `hmgraph.CreateVertexMap`, `hmgraph.CreateArcMap` and `hmgraph.CreateEdgeMap`.
In addition to the graph itself, each map stores a label and a default value, which also defines the map's value type.

```go
g := hmgraph.NewGraph()
u := g.CreateVertex()
v := g.CreateVertex()
vertexMap := hmgraph.CreateVertexMap(g, "weight", 1.0)
vertexMap.Set(u, 3.2)
vertexMap.Set(v, 6.2)
w := g.CreateVertex()
a := u.CreateArc(v)
e, _ := w.CreateEdge(v)
arcMap := hmgraph.CreateArcMap(g, "intCost", 0)
arcMap.Set(a, 2)
edgeMap := hmgraph.CreateEdgeMap(g, "color", "blue")
edgeMap.Set(e, "teal")
```

When printing graphs, all values stored in maps will be contained in the output:

```
Graph
{
  Vertex #0 (weight=3.2;)
  {
    Arc #0-->#1 (intCost=2;)
  }
  Vertex #1 (weight=6.2;)
  {
    Edge #2---#1 (color=teal;)
    Arc #0-->#1 (intCost=2;)
  }
  Vertex #2 (weight=1;)
  {
    Edge #2---#1 (color=teal;)
  }
}
```

### Disposing Maps

Maps are registered with a graph to adapt to changes (deletion/addition of elements). Therefore, they continue to exist
until they are explicitly disposed.

If a map is no longer needed, it should be de-registered by calling the respective ```Dispose()``` method. 
After disposing, a map should not be accessed.

```go
g := hmgraph.NewGraph()
u := g.CreateVertex()
vertexMap := hmgraph.CreateVertexMap(g, "weight", 1.0)
vertexMap.Set(u, 3.2)

// do stuff here

vertexMap.Dispose()

// map is no longer kept in-sync with graph changes
```

In functions using a map to store helper values, it is best practice to add a ```defer```-block right after creating the map:

```go
func foo(g *hmgraph.Graph) {
	visited := hmgraph.CreateVertexMap(g, "visited", false)
	defer visited.Dispose()
```

## Query methods

Graphs, vertices, arcs and edges provide standard operations. To protect the integrity of the internal structure, access
follows the visitor pattern:

### Graph queries

* `NodeCount() int`, `ArcCount() int` and `EdgeCount() int` provide vertex, arc and edge counts, respectively.
* `ForVertices(visitor func(*Vertex))`, `ForArcs(visitor func(*Arc))`, `ForEdges(visitor func(*Edge))` call a given
  function for every vertex, arc or edge

### Vertex queries

* `InDegree() int`, `OutDegree() int` und  `Degree() int` provide the number of a single node's inbound/outbound arcs or
  the number of incident edges.* `ForOutArcs(visitor func(*Arc))`, `ForInArcs(visitor func(*Arc))` und
  `ForEdges(visitor func(*Edge))` call a given function for every inbound arc, outbound arc or incident edge. 

