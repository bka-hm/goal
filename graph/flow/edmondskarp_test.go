package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

func TestPanicOnEdges(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(2)
	vs[0].CreateEdge(vs[1])
	kapa := hmgraph.CreateArcMap[int](g, "kapa", 0)
	assert.Panics(t, func() {
		EdmondsKarp(g, kapa, vs[0], vs[1], true)
	})
	x, y, _ := g.MapCount()
	assert.Equal(t, 0, x)
	assert.Equal(t, 1, y)
}

func TestSimple(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(2)
	vs[0].CreateArc(vs[1])
	kapa := hmgraph.CreateArcMap[int](g, "kapa", 2)
	flow, cut := EdmondsKarp(g, kapa, vs[0], vs[1], true)
	assert.Equal(t, 2, checkFlow(t, g, kapa, flow, vs[0], vs[1]))
	assert.ElementsMatch(t, []*hmgraph.Vertex{vs[0]}, cut)
	x, y, _ := g.MapCount()
	assert.Equal(t, 0, x)
	assert.Equal(t, 2, y)
}

func TestMinimalRelabel(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(3)
	vs[0].CreateArc(vs[1])
	vs[1].CreateArc(vs[2])
	kapa := hmgraph.CreateArcMap[int](g, "kapa", 2)
	flow, cut := EdmondsKarp(g, kapa, vs[0], vs[2], true)
	assert.Equal(t, 2, checkFlow(t, g, kapa, flow, vs[0], vs[2]))
	x, y, _ := g.MapCount()
	assert.Equal(t, 0, x)
	assert.Equal(t, 2, y)
	assert.ElementsMatch(t, []*hmgraph.Vertex{vs[0]}, cut)
}

func TestExerciseExample(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	kapa := hmgraph.CreateArcMap[int](g, "kapa", 2)
	arc01 := vs[0].CreateArc(vs[2])
	kapa.Set(arc01, 8)
	arc04 := vs[0].CreateArc(vs[4])
	kapa.Set(arc04, 4)
	arc21 := vs[2].CreateArc(vs[1])
	kapa.Set(arc21, 2)
	arc23 := vs[2].CreateArc(vs[3])
	kapa.Set(arc23, 5)
	arc34 := vs[3].CreateArc(vs[4])
	kapa.Set(arc34, 5)
	arc35 := vs[3].CreateArc(vs[5])
	kapa.Set(arc35, 3)
	arc45 := vs[4].CreateArc(vs[5])
	kapa.Set(arc45, 5)
	arc52 := vs[5].CreateArc(vs[2])
	kapa.Set(arc52, 3)
	arc51 := vs[5].CreateArc(vs[1])
	kapa.Set(arc51, 9)

	flow, cut := EdmondsKarp(g, kapa, vs[0], vs[1], true)
	assert.Equal(t, 10, checkFlow(t, g, kapa, flow, vs[0], vs[1]))
	x, y, _ := g.MapCount()
	assert.Equal(t, 0, x)
	assert.Equal(t, 2, y)
	assert.ElementsMatch(t, []*hmgraph.Vertex{vs[0], vs[2], vs[4], vs[3]}, cut)
}

func checkFlow(t *testing.T, g *hmgraph.Graph, kapa *hmgraph.ArcMap[int], flow *hmgraph.ArcMap[int],
	source *hmgraph.Vertex, target *hmgraph.Vertex) int {
	// check all kapa
	g.ForArcs(func(arc *hmgraph.Arc) {
		assert.True(t, flow.Get(arc) >= 0 && flow.Get(arc) <= kapa.Get(arc))
	})
	// check zero sum
	sflow := 0
	tflow := 0
	g.ForVertices(func(v *hmgraph.Vertex) {
		in := 0
		out := 0
		v.ForOutArcs(func(arc *hmgraph.Arc) {
			out += flow.Get(arc)
		})
		v.ForInArcs(func(arc *hmgraph.Arc) {
			in += flow.Get(arc)
		})
		if v == source {
			sflow = out - in
		} else if v == target {
			tflow = in - out
		} else {
			assert.Equal(t, 0, out-in)
		}
	})
	assert.Equal(t, 0, sflow-tflow)
	return sflow
}

func TestBackflowExample(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(13)
	kapa := hmgraph.CreateArcMap[int](g, "kapa", 0)

	arc01 := vs[0].CreateArc(vs[1])
	kapa.Set(arc01, 2)

	arc02 := vs[0].CreateArc(vs[2])
	kapa.Set(arc02, 3)

	arc03 := vs[0].CreateArc(vs[3])
	kapa.Set(arc03, 5)

	arc14 := vs[1].CreateArc(vs[4])
	kapa.Set(arc14, 2)

	arc24 := vs[2].CreateArc(vs[4])
	kapa.Set(arc24, 2)

	arc23 := vs[2].CreateArc(vs[3])
	kapa.Set(arc23, 2)

	arc46 := vs[4].CreateArc(vs[6])
	kapa.Set(arc46, 6)

	arc54 := vs[5].CreateArc(vs[4])
	kapa.Set(arc54, 2)

	arc35 := vs[3].CreateArc(vs[5])
	kapa.Set(arc35, 4)

	arc56 := vs[5].CreateArc(vs[6])
	kapa.Set(arc56, 1)

	arc57 := vs[5].CreateArc(vs[7])
	kapa.Set(arc57, 2)

	arc68 := vs[6].CreateArc(vs[8])
	kapa.Set(arc68, 5)

	arc87 := vs[8].CreateArc(vs[7])
	kapa.Set(arc87, 2)

	arc810 := vs[8].CreateArc(vs[10])
	kapa.Set(arc810, 5)

	arc79 := vs[7].CreateArc(vs[9])
	kapa.Set(arc79, 6)

	arc910 := vs[9].CreateArc(vs[10])
	kapa.Set(arc910, 3)

	arc1012 := vs[10].CreateArc(vs[12])
	kapa.Set(arc1012, 6)

	arc911 := vs[9].CreateArc(vs[11])
	kapa.Set(arc911, 5)

	arc1112 := vs[11].CreateArc(vs[12])
	kapa.Set(arc1112, 6)

	flow, cut := EdmondsKarp(g, kapa, vs[0], vs[12], true)
	assert.Equal(t, 7, checkFlow(t, g, kapa, flow, vs[0], vs[12]))
	x, y, _ := g.MapCount()
	assert.Equal(t, 0, x)
	assert.Equal(t, 2, y)
	assert.ElementsMatch(t, []*hmgraph.Vertex{vs[0], vs[1], vs[2], vs[3], vs[4], vs[5], vs[6]}, cut)
}

func TestBackflowSimple(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(4)
	kapa := hmgraph.CreateArcMap[int](g, "kapa", 0)

	arc01 := vs[0].CreateArc(vs[1])
	kapa.Set(arc01, 8)

	arc02 := vs[0].CreateArc(vs[2])
	kapa.Set(arc02, 4)

	arc12 := vs[1].CreateArc(vs[2])
	kapa.Set(arc12, 8)

	arc13 := vs[1].CreateArc(vs[3])
	kapa.Set(arc13, 4)

	arc23 := vs[2].CreateArc(vs[3])
	kapa.Set(arc23, 8)

	flow, cut := EdmondsKarp(g, kapa, vs[0], vs[3], true)
	assert.Equal(t, 12, checkFlow(t, g, kapa, flow, vs[0], vs[3]))
	x, y, _ := g.MapCount()
	assert.Equal(t, 0, x)
	assert.Equal(t, 2, y)
	assert.ElementsMatch(t, []*hmgraph.Vertex{vs[0]}, cut)
}

func TestSimpleMultiTerminalFlow(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(2)
	kapa := hmgraph.CreateArcMap(g, "kapa", 0)
	arc01 := vs[0].CreateArc(vs[1])
	kapa.Set(arc01, 8)
	support := hmgraph.CreateVertexMap(g, "support", 2)
	support.Set(vs[1], -3)
	flow, cut := MultiTerminalMaxFlow(g, kapa, support, true)
	x, y, _ := g.MapCount()
	assert.Equal(t, 2, g.VertexCount())
	assert.Equal(t, 1, x)
	assert.Equal(t, 2, y)
	assert.ElementsMatch(t, []*hmgraph.Vertex{}, cut)
	assert.Equal(t, 2, flow.Get(arc01))
}
