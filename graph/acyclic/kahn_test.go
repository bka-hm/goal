package acyclic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	graph_errors "gitlab.lrz.de/hm/goal/graph/errors"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

func TestSmallStandardCase(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[1].CreateArc(vs[4])
	vs[0].CreateArc(vs[5])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	sorting, err := TopologicalOrder(g)
	assert.NoError(t, err)
	assertValidOrdering(g, sorting, t)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed.")
}

func TestErrorOnEdge(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[1].CreateArc(vs[4])
	vs[0].CreateArc(vs[5])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	vs[5].CreateEdge(vs[4])
	_, err := TopologicalOrder(g)
	assert.Error(t, err)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed on error.")
}

func TestErrorOnCycle(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[1].CreateArc(vs[4])
	vs[0].CreateArc(vs[5])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	vs[4].CreateArc(vs[3])
	_, err := TopologicalOrder(g)
	var cycleErr *graph_errors.ContainsCycleError
	assert.ErrorAs(t, err, &cycleErr)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed on cycle.")
}

func TestIsDagTrue(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[0].CreateArc(vs[5])
	vs[5].CreateArc(vs[4])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	vs[4].CreateArc(vs[3])
	isDagResult := IsDag(g)
	assert.True(t, isDagResult)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed on cycle.")
}

func TestIsDagFalseHasCycle(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[0].CreateArc(vs[5])
	vs[5].CreateArc(vs[4])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	vs[1].CreateArcs([]*hmgraph.Vertex{vs[4], vs[2]})
	vs[4].CreateArc(vs[3])
	isDagResult := IsDag(g)
	assert.False(t, isDagResult)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed on cycle.")
}

func TestIsDagFalseHasEdge(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[3].CreateArc(vs[1])
	vs[0].CreateArc(vs[5])
	vs[5].CreateArc(vs[4])
	vs[2].CreateArcs([]*hmgraph.Vertex{vs[4], vs[1]})
	vs[4].CreateEdge(vs[3])
	vs[4].CreateArc(vs[3])
	isDagResult := IsDag(g)
	assert.False(t, isDagResult)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed on cycle.")
}

func TestSelfLoop(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(3)
	vs[0].CreateArc(vs[1])
	vs[1].CreateArc(vs[1])
	vs[2].CreateArc(vs[0])

	_, err := TopologicalOrder(g)
	assert.Error(t, err)

	_, ok := err.(*graph_errors.ContainsCycleError)
	assert.True(t, ok, "Error should be of type CycleError")

	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed on cycle.")
}

func TestIsDagEmpty(t *testing.T) {
	g := hmgraph.NewGraph()
	g.CreateVertices(0)

	assert.True(t, IsDag(g), "Empty graph should be a valid DAG")

	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed.")
}

func TestComplexCycleWithBranching(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(7)
	vs[0].CreateArc(vs[1])
	vs[0].CreateArc(vs[2])
	vs[1].CreateArc(vs[3])
	vs[2].CreateArc(vs[3])
	vs[3].CreateArc(vs[4])
	vs[4].CreateArc(vs[5])
	vs[5].CreateArc(vs[6])
	vs[6].CreateArc(vs[1])

	_, err := TopologicalOrder(g)
	assert.Error(t, err)

	cycleErr, ok := err.(*graph_errors.ContainsCycleError)
	assert.True(t, ok, "Error should be of type CycleError")
	assert.Greater(t, len(cycleErr.Cycle()), 0)

	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed.")
}

func assertValidOrdering(g *hmgraph.Graph, sorting []*hmgraph.Vertex, t *testing.T) {
	numbering := hmgraph.CreateVertexMap(g, "testNumber", -1)
	for i, vertex := range sorting {
		numbering.Set(vertex, i)
		vertex.ForOutArcs(func(arc *hmgraph.Arc) {
			assert.True(t, numbering.Get(arc.Target()) == -1)
		})
	}
	numbering.Dispose()
	assert.Equal(t, g.VertexCount(), len(sorting))
}
