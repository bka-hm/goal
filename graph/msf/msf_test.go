package msf

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.lrz.de/hm/goal/base"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

func getAlgos[T base.Number]() []MsfAlgo[T] {
	return []MsfAlgo[T]{PrimMsf[T], KruskalMsf[T], BoruvkaMsf[T]}
}

func TestEmpty(t *testing.T) {
	g := hmgraph.NewGraph()
	cost := hmgraph.CreateEdgeMap(g, "cost", 0)
	for _, alg := range getAlgos[int]() {
		assert.Panics(t, func() {
			alg(g, cost)
		})
	}
}

func TestArcs(t *testing.T) {
	g := hmgraph.NewGraph()
	cost := hmgraph.CreateEdgeMap(g, "cost", 0)
	vs := g.CreateVertices(2)
	vs[0].CreateEdge(vs[1])
	vs[0].CreateArc(vs[1])
	for _, alg := range getAlgos[int]() {
		assert.Panics(t, func() {
			alg(g, cost)
		})
	}
}

func TestSimple(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(2)
	cost := hmgraph.CreateEdgeMap(g, "cost", 0)
	e := vs[0].CreateEdge(vs[1])
	cost.Set(e, 3)
	expected := []*hmgraph.Edge{e}
	for _, alg := range getAlgos[int]() {
		tree := alg(g, cost)
		assertEquals(t, expected, tree)
	}
}

func TestSingleVertex(t *testing.T) {
	g := hmgraph.NewGraph()
	g.CreateVertices(1)
	cost := hmgraph.CreateEdgeMap(g, "cost", 0)
	expected := []*hmgraph.Edge{}
	for _, alg := range getAlgos[int]() {
		tree := alg(g, cost)
		assertEquals(t, expected, tree)
	}
}

func TestNoEdges(t *testing.T) {
	g := hmgraph.NewGraph()
	g.CreateVertices(6)
	cost := hmgraph.CreateEdgeMap(g, "cost", 0)
	expected := []*hmgraph.Edge{}
	for _, alg := range getAlgos[int]() {
		tree := alg(g, cost)
		assertEquals(t, expected, tree)
	}
}

func TestSimpleTree1(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(3)
	a, b, c := vs[0], vs[1], vs[2]
	aEdges := a.CreateEdges([]*hmgraph.Vertex{b, c})
	ab, ac := aEdges[0], aEdges[1]
	bc := b.CreateEdge(c)
	edgeMap := hmgraph.CreateEdgeMap(g, "cost", math.MaxInt)
	edgeMap.Set(ab, 4)
	edgeMap.Set(ac, 10)
	edgeMap.Set(bc, 14)

	expected := []*hmgraph.Edge{ab, ac}
	for _, alg := range getAlgos[int]() {
		tree := alg(g, edgeMap)
		assertEquals(t, expected, tree)
	}
}

func TestSimpleTree2(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(7)
	a, b, c, d, e, f, h := vs[0], vs[1], vs[2], vs[3], vs[4], vs[5], vs[6]

	aEdges := a.CreateEdges([]*hmgraph.Vertex{b, d})
	ab, ad := aEdges[0], aEdges[1]
	bEdges := b.CreateEdges([]*hmgraph.Vertex{c, d, e})
	bc, bd, be := bEdges[0], bEdges[1], bEdges[2]
	cf := c.CreateEdge(f)
	de := d.CreateEdge(e)
	eEdges := e.CreateEdges([]*hmgraph.Vertex{f, h})
	ef, eh := eEdges[0], eEdges[1]

	edgeMap := hmgraph.CreateEdgeMap(g, "cost", math.MaxInt)
	edgeMap.Set(ab, 7)
	edgeMap.Set(ad, 4)
	edgeMap.Set(bc, 9)
	edgeMap.Set(bd, 6)
	edgeMap.Set(be, 8)
	edgeMap.Set(cf, 5)
	edgeMap.Set(de, 20)
	edgeMap.Set(ef, 10)
	edgeMap.Set(eh, 11)
	expected := []*hmgraph.Edge{ad, cf, be, eh, bc, bd}

	for _, alg := range getAlgos[int]() {
		tree := alg(g, edgeMap)
		assertEquals(t, expected, tree)
	}
}

func TestUnconnected(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(5)
	cost := hmgraph.CreateEdgeMap(g, "cost", 0)

	v0Edges := vs[0].CreateEdges([]*hmgraph.Vertex{vs[1], vs[4]})
	e1, e3 := v0Edges[0], v0Edges[1]
	e2 := vs[1].CreateEdge(vs[4])
	e4 := vs[3].CreateEdge(vs[2])
	cost.Set(e1, 3)
	cost.Set(e2, 2)
	cost.Set(e3, 1)
	cost.Set(e4, 4)

	expected := []*hmgraph.Edge{e2, e3, e4}
	for _, alg := range getAlgos[int]() {
		tree := alg(g, cost)
		assertEquals(t, expected, tree)
	}
}

func TestSpanningForest(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(8)
	a, b, c, d, e, f, h, i := vs[0], vs[1], vs[2], vs[3], vs[4], vs[5], vs[6], vs[7]

	aEdges := a.CreateEdges([]*hmgraph.Vertex{b, c, d})
	ab, ac, ad := aEdges[0], aEdges[1], aEdges[2]
	bd := b.CreateEdge(d)
	cd := c.CreateEdge(d)
	eEdges := e.CreateEdges([]*hmgraph.Vertex{f, h, i})
	ef, eh, ei := eEdges[0], eEdges[1], eEdges[2]
	fEdges := f.CreateEdges([]*hmgraph.Vertex{h, i})
	fh, fi := fEdges[0], fEdges[1]
	hi := h.CreateEdge(i)

	cost := hmgraph.CreateEdgeMap(g, "cost", math.MaxInt)
	cost.Set(ab, 3)
	cost.Set(ac, 2)
	cost.Set(ad, 1)
	cost.Set(bd, 4)
	cost.Set(cd, 5)
	cost.Set(ef, 3)
	cost.Set(eh, 2)
	cost.Set(ei, 7)
	cost.Set(fh, 5)
	cost.Set(hi, 8)
	cost.Set(fi, 9)

	expected := []*hmgraph.Edge{ad, ac, ab, eh, ef, ei}

	for _, alg := range getAlgos[int]() {
		tree := alg(g, cost)
		assertEquals(t, expected, tree)
	}
}

func assertEquals(t *testing.T, edges []*hmgraph.Edge, tree []*hmgraph.Edge) {
	assert.Equal(t, len(edges), len(tree), "wrong number of edges")
	for _, e := range edges {
		assert.Contains(t, tree, e, "missing edge")
	}
}
