package sp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// TestDijkstraSimple: directed graph, shortest paths via arc relaxation.
//
//	0 -2→ 1 -5→ 3
//	0 -6→ 2
//	1 -1→ 2 -1→ 3
//
// Expected: dist = [0, 2, 3, 4], path to 3 via 0→1→2→3.
func TestDijkstraSimple(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(4)
	a01 := vs[0].CreateArc(vs[1])
	aCost.Set(a01, 2.0)
	a02 := vs[0].CreateArc(vs[2])
	aCost.Set(a02, 6.0)
	a12 := vs[1].CreateArc(vs[2])
	aCost.Set(a12, 1.0)
	a13 := vs[1].CreateArc(vs[3])
	aCost.Set(a13, 5.0)
	a23 := vs[2].CreateArc(vs[3])
	aCost.Set(a23, 1.0)
	dist, pred, err := Dijkstra(g, aCost, eCost, vs[0])
	assert.NoError(t, err)
	assert.Equal(t, 0.0, dist.Get(vs[0]))
	assert.Equal(t, 2.0, dist.Get(vs[1]))
	assert.Equal(t, 3.0, dist.Get(vs[2]))
	assert.Equal(t, 4.0, dist.Get(vs[3]))
	assert.Nil(t, pred.Get(vs[0]))
	assert.Equal(t, vs[0], pred.Get(vs[1]).Source())
	assert.Equal(t, vs[1], pred.Get(vs[2]).Source())
	assert.Equal(t, vs[2], pred.Get(vs[3]).Source())
	vm, am, em := g.MapCount()
	assert.Equal(t, 2, vm) // dist + predecessor; handles is disposed
	assert.Equal(t, 1, am)
	assert.Equal(t, 1, em)
}

// TestDijkstraUndirectedEdge: undirected edges are traversable in both directions.
func TestDijkstraUndirectedEdge(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(3)
	e01 := vs[0].CreateEdge(vs[1])
	eCost.Set(e01, 3.0)
	e12 := vs[1].CreateEdge(vs[2])
	eCost.Set(e12, 2.0)
	dist, pred, err := Dijkstra(g, aCost, eCost, vs[0])
	assert.NoError(t, err)
	assert.Equal(t, 0.0, dist.Get(vs[0]))
	assert.Equal(t, 3.0, dist.Get(vs[1]))
	assert.Equal(t, 5.0, dist.Get(vs[2]))
	assert.Nil(t, pred.Get(vs[0]))
	assert.Equal(t, vs[0], pred.Get(vs[1]).Source())
	assert.Equal(t, vs[1], pred.Get(vs[2]).Source())
	vm, am, em := g.MapCount()
	assert.Equal(t, 2, vm)
	assert.Equal(t, 1, am)
	assert.Equal(t, 1, em)
}

// TestDijkstraUnreachable: a vertex with no path from the source keeps MaxValue distance
// and a nil predecessor.
func TestDijkstraUnreachable(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(3)
	a01 := vs[0].CreateArc(vs[1])
	aCost.Set(a01, 1.0)
	// vs[2] has no incoming arcs or edges
	dist, pred, err := Dijkstra(g, aCost, eCost, vs[0])
	assert.NoError(t, err)
	assert.Equal(t, 0.0, dist.Get(vs[0]))
	assert.Equal(t, 1.0, dist.Get(vs[1]))
	assert.Equal(t, base.MaxValue[float64](), dist.Get(vs[2]))
	assert.Nil(t, pred.Get(vs[0]))
	assert.NotNil(t, pred.Get(vs[1]))
	assert.Nil(t, pred.Get(vs[2]))
}

// TestDijkstraSingleVertex: source with no neighbours has distance 0 and no predecessor.
func TestDijkstraSingleVertex(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(1)
	dist, pred, err := Dijkstra(g, aCost, eCost, vs[0])
	assert.NoError(t, err)
	assert.Equal(t, 0.0, dist.Get(vs[0]))
	assert.Nil(t, pred.Get(vs[0]))
}

// TestDijkstraNilGraphPanics: passing a nil graph panics.
func TestDijkstraNilGraphPanics(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(1)
	assert.Panics(t, func() {
		_, _, err := Dijkstra(nil, aCost, eCost, vs[0])
		if err != nil {
			return
		}
	})
}

// TestDijkstraNilSourcePanics: passing a nil source panics.
func TestDijkstraNilSourcePanics(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	assert.Panics(t, func() {
		_, _, err := Dijkstra(g, aCost, eCost, nil)
		if err != nil {
			return
		}
	})
}

// TestDijkstraNegativeEdgePanics: a negative edge cost violates the algorithm's precondition.
func TestDijkstraNegativeEdgePanics(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(2)
	edge := vs[0].CreateEdge(vs[1])
	eCost.Set(edge, -1.0)
	assert.Panics(t, func() {
		_, _, err := Dijkstra(g, aCost, eCost, vs[0])
		if err != nil {
			return
		}
	})
}

// TestDijkstraSinglePair: path is returned in forward order (source → target).
//
//	0 -1→ 1 -2→ 2 -3→ 3
func TestDijkstraSinglePair(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(4)
	a01 := vs[0].CreateArc(vs[1])
	aCost.Set(a01, 1.0)
	a12 := vs[1].CreateArc(vs[2])
	aCost.Set(a12, 2.0)
	a23 := vs[2].CreateArc(vs[3])
	aCost.Set(a23, 3.0)
	path, err := DijkstraSinglePairShortestPath(g, aCost, eCost, vs[0], vs[3])
	assert.NoError(t, err)
	assert.Equal(t, 3, len(path))
	assert.Equal(t, vs[0], path[0].Source())
	assert.Equal(t, vs[1], path[1].Source())
	assert.Equal(t, vs[2], path[2].Source())
	vm, am, em := g.MapCount()
	assert.Equal(t, 0, vm) // dist and predecessor both disposed
	assert.Equal(t, 1, am)
	assert.Equal(t, 1, em)
}

// TestDijkstraSinglePairUnreachable: unreachable target returns nil path without error.
func TestDijkstraSinglePairUnreachable(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(2)
	// no arc from vs[0] to vs[1]
	path, err := DijkstraSinglePairShortestPath(g, aCost, eCost, vs[0], vs[1])
	assert.NoError(t, err)
	assert.Nil(t, path)
	vm, _, _ := g.MapCount()
	assert.Equal(t, 0, vm)
}

// TestDijkstraSinglePairSameVertex: source == target returns nil path without error.
func TestDijkstraSinglePairSameVertex(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(2)
	vs[0].CreateArc(vs[1])
	path, err := DijkstraSinglePairShortestPath(g, aCost, eCost, vs[0], vs[0])
	assert.NoError(t, err)
	assert.Nil(t, path)
	vm, _, _ := g.MapCount()
	assert.Equal(t, 0, vm)
}

// TestDijkstraSinglePairNegativeCycle: negative cycle returns error and no path.
func TestDijkstraSinglePairNegativeCycle(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(3)
	a01 := vs[0].CreateArc(vs[1])
	aCost.Set(a01, 1.0)
	a12 := vs[1].CreateArc(vs[2])
	aCost.Set(a12, -1.0)
	a21 := vs[2].CreateArc(vs[1])
	aCost.Set(a21, -1.0)
	path, err := DijkstraSinglePairShortestPath(g, aCost, eCost, vs[0], vs[2])
	assert.Error(t, err)
	assert.Nil(t, path)
	vm, _, _ := g.MapCount()
	assert.Equal(t, 0, vm)
}

// TestDijkstraNegativeCycleReturnsError: a negative arc cycle is detected and reported as an error;
// both output maps are disposed and nil.
func TestDijkstraNegativeCycleReturnsError(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	vs := g.CreateVertices(3)
	a01 := vs[0].CreateArc(vs[1])
	aCost.Set(a01, 1.0)
	a12 := vs[1].CreateArc(vs[2])
	aCost.Set(a12, -1.0)
	a21 := vs[2].CreateArc(vs[1])
	aCost.Set(a21, -1.0)
	dist, pred, err := Dijkstra(g, aCost, eCost, vs[0])
	assert.Error(t, err)
	assert.Nil(t, dist)
	assert.Nil(t, pred)
	vm, am, em := g.MapCount()
	assert.Equal(t, 0, vm)
	assert.Equal(t, 1, am)
	assert.Equal(t, 1, em)
}
