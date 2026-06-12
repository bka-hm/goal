package sp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

// TestJohnson verifies the full Johnson workflow on a graph with a negative arc.
//
// Graph:
//
//	0 -1→ 1
//	1 -(-2)→ 2
//	0 -5→ 2
//
// Expected heights (multi-source BFM, all sources at 0): h = [0, 0, -2]
// Expected transformed costs: w'(0→1)=1, w'(1→2)=0, w'(0→2)=7 (all non-negative)
// Expected original distances from 0 via Dijkstra + JohnsonInverse: d = [0, 1, -1]
func TestJohnsonWithNegativeArc(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	vs := g.CreateVertices(3)
	a01 := vs[0].CreateArc(vs[1])
	aCost.Set(a01, 1.0)
	a12 := vs[1].CreateArc(vs[2])
	aCost.Set(a12, -2.0)
	a02 := vs[0].CreateArc(vs[2])
	aCost.Set(a02, 5.0)

	heights, transformed, err := Johnson(g, aCost)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, heights.Get(vs[0]))
	assert.Equal(t, 0.0, heights.Get(vs[1]))
	assert.Equal(t, -2.0, heights.Get(vs[2]))
	assert.Equal(t, 1.0, transformed.Get(a01))
	assert.Equal(t, 0.0, transformed.Get(a12))
	assert.Equal(t, 7.0, transformed.Get(a02))

	// all transformed costs must be non-negative
	g.ForArcs(func(a *hmgraph.Arc) {
		assert.GreaterOrEqual(t, transformed.Get(a), 0.0)
	})

	// run Dijkstra with transformed costs, then invert
	eCost := hmgraph.CreateEdgeMap(g, "cost", 0.0)
	transformedDist, pred, dijkstraErr := Dijkstra(g, transformed, eCost, vs[0])
	assert.NoError(t, dijkstraErr)
	pred.Dispose()

	dist := JohnsonInverse(g, vs[0], transformedDist, heights)
	transformedDist.Dispose()
	assert.Equal(t, 0.0, dist.Get(vs[0]))
	assert.Equal(t, 1.0, dist.Get(vs[1]))
	assert.Equal(t, -1.0, dist.Get(vs[2]))

	dist.Dispose()
	heights.Dispose()
	transformed.Dispose()
	vm, am, em := g.MapCount()
	assert.Equal(t, 0, vm)
	assert.Equal(t, 1, am) // aCost
	assert.Equal(t, 1, em) // eCost
}

// TestJohnsonPanicsWithEdges: undirected edges are not supported and cause a panic.
func TestJohnsonPanicsWithEdges(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	vs := g.CreateVertices(2)
	vs[0].CreateEdge(vs[1])
	assert.Panics(t, func() {
		_, _, err := Johnson(g, aCost)
		if err != nil {
			return
		}
	})
}

// TestJohnsonNegativeCycleReturnsError: a negative arc cycle is reported as an error;
// both output maps are nil.
func TestJohnsonNegativeCycleReturnsError(t *testing.T) {
	g := hmgraph.NewGraph()
	aCost := hmgraph.CreateArcMap(g, "cost", 0.0)
	vs := g.CreateVertices(2)
	a01 := vs[0].CreateArc(vs[1])
	aCost.Set(a01, -1.0)
	a10 := vs[1].CreateArc(vs[0])
	aCost.Set(a10, -1.0)

	heights, transformed, err := Johnson(g, aCost)
	assert.Error(t, err)
	assert.Nil(t, heights)
	assert.Nil(t, transformed)
	vm, am, _ := g.MapCount()
	assert.Equal(t, 0, vm)
	assert.Equal(t, 1, am) // only aCost remains
}
