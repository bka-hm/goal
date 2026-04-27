package connectivity

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

func vertexIndices(vs []*hmgraph.Vertex) []int {
	indices := make([]int, len(vs))
	for i, v := range vs {
		indices[i] = v.Index()
	}
	return indices
}

func TestStronglyConnectedEmptyGraph(t *testing.T) {
	g := hmgraph.NewGraph()

	components := ConnectedComponents(g, StronglyConnected)

	assert.Equal(t, 0, len(components))
}

func TestSimpleArcGraph(t *testing.T) {
	// arrange
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(3)
	vs[1].CreateArc(vs[2])
	vs[2].CreateArc(vs[1])

	//act
	components := ConnectedComponents(g, StronglyConnected)

	// assert
	assert.Equal(t, 2, len(components))

	idx0 := slices.IndexFunc(components, func(c []*hmgraph.Vertex) bool {
		return slices.Contains(c, vs[0])
	})
	assert.Equal(t, 1, len(components[idx0]))

	idx12 := slices.IndexFunc(components, func(c []*hmgraph.Vertex) bool {
		return slices.Contains(c, vs[1])
	})
	assert.NotEqual(t, -1, idx12)
	assert.ElementsMatch(t, vertexIndices(components[idx12]), []int{1, 2})
}

func TestSimpleArcGraphWithTopology(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(3)
	vs[1].CreateArc(vs[2])
	vs[2].CreateArc(vs[1])
	vs[1].CreateArc(vs[0])

	components := ConnectedComponents(g, StronglyConnected)

	assert.Equal(t, 2, len(components))
	assert.ElementsMatch(t, vertexIndices(components[1]), []int{0})
	assert.ElementsMatch(t, vertexIndices(components[0]), []int{1, 2})
}

func TestGraphWithTwoComponents(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(4)

	vs[1].CreateEdge(vs[2])
	vs[2].CreateArc(vs[3])
	vs[0].CreateEdge(vs[3])

	components := ConnectedComponents(g, StronglyConnected)
	assert.Equal(t, 2, len(components))
	assert.ElementsMatch(t, vertexIndices(components[0]), []int{1, 2})
	assert.ElementsMatch(t, vertexIndices(components[1]), []int{3, 0})
}

func TestGraphWeaklyConnected(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(4)

	vs[1].CreateEdge(vs[2])
	vs[2].CreateArc(vs[3])
	vs[0].CreateEdge(vs[3])

	components := ConnectedComponents(g, WeaklyConnected)
	assert.Equal(t, 1, len(components))
	assert.ElementsMatch(t, vertexIndices(components[0]), []int{0, 1, 2, 3})
}

func TestIsStronglyConnected_SingleVertex(t *testing.T) {
	g := hmgraph.NewGraph()
	g.CreateVertices(1)

	assert.True(t, IsStronglyConnected(g))
}

func TestIsStronglyConnected_Cycle(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(3)
	vs[0].CreateArc(vs[1])
	vs[1].CreateArc(vs[2])
	vs[2].CreateArc(vs[0])

	assert.True(t, IsStronglyConnected(g))
}

func TestIsStronglyConnected_DAG(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(3)
	vs[0].CreateArc(vs[1])
	vs[1].CreateArc(vs[2])

	assert.False(t, IsStronglyConnected(g))
}

func TestIsStronglyConnected_DisconnectedVertices(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(4)
	vs[0].CreateArc(vs[1])
	vs[1].CreateArc(vs[0])
	// vs[2] and vs[3] are isolated

	assert.False(t, IsStronglyConnected(g))
}

func TestIsWeaklyConnected_ArcChain(t *testing.T) {
	// All arcs point one way — not strongly, but weakly connected
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(4)
	vs[0].CreateArc(vs[1])
	vs[1].CreateArc(vs[2])
	vs[2].CreateArc(vs[3])

	assert.True(t, IsWeaklyConnected(g))
}

func TestIsWeaklyConnected_MixedEdgesAndArcs(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(4)
	vs[0].CreateEdge(vs[1])
	vs[1].CreateArc(vs[2])
	vs[2].CreateEdge(vs[3])

	assert.True(t, IsWeaklyConnected(g))
}

func TestIsWeaklyConnected_Disconnected(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(4)
	vs[0].CreateArc(vs[1])
	// vs[2] and vs[3] are a separate pair with no connection to vs[0]/vs[1]
	vs[2].CreateArc(vs[3])

	assert.False(t, IsWeaklyConnected(g))
}

func TestIsWeaklyConnected_IsolatedVertex(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(3)
	vs[0].CreateArc(vs[1])
	// vs[2] is isolated

	assert.False(t, IsWeaklyConnected(g))
}

func TestArcGraph(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(9)
	vs[0].CreateArcs([]*hmgraph.Vertex{vs[1], vs[2]})
	vs[1].CreateArc(vs[8])
	vs[2].CreateArc(vs[3])
	vs[3].CreateArc(vs[2])
	vs[4].CreateArcs([]*hmgraph.Vertex{vs[3], vs[7]})
	vs[5].CreateArc(vs[4])
	vs[6].CreateArcs([]*hmgraph.Vertex{vs[5], vs[7]})
	vs[7].CreateArcs([]*hmgraph.Vertex{vs[2], vs[6]})
	vs[8].CreateArcs([]*hmgraph.Vertex{vs[0], vs[7]})

	components := ConnectedComponents(g, StronglyConnected)

	assert.ElementsMatch(t, vertexIndices(components[0]), []int{0, 1, 8})
	assert.ElementsMatch(t, vertexIndices(components[2]), []int{2, 3})
	assert.ElementsMatch(t, vertexIndices(components[1]), []int{4, 5, 6, 7})
}
