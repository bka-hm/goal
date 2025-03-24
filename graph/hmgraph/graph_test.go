package hmgraph

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphCreation(t *testing.T) {
	g := NewGraph()
	assert.Equal(t, 0, g.VertexCount(), "new graph should be empty.")
	assert.Equal(t, 0, g.ArcCount(), "new graph should be empty.")
	assert.Equal(t, 0, g.EdgeCount(), "new graph should be empty.")
}

func TestVertexCreation(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	assert.Equal(t, g, v.Graph(), "vertex should identify graph")
	assert.Equal(t, 0, v.Index(), "single vertex's index should be 0")
	assert.Equal(t, 0, v.InDegree())
	assert.Equal(t, 0, v.OutDegree())
	assert.Equal(t, 0, v.Degree())
}

func TestArcCreation(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	a := v.CreateArc(u)
	assert.Equal(t, 0, a.Index(), "single arc's index should be 0")
	assert.Equal(t, 0, v.InDegree())
	assert.Equal(t, 1, v.OutDegree())
	assert.Equal(t, 1, u.InDegree())
	assert.Equal(t, 0, u.OutDegree())
	assert.Equal(t, v, a.Source())
	assert.Equal(t, u, a.Target())
}

func TestEdgeCreation(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	e := v.CreateEdge(u)
	assert.Equal(t, 0, e.Index(), "first edge should have index 0")
	assert.Equal(t, 1, v.Degree())
	assert.Equal(t, 1, u.Degree())
	assert.True(t, e.IsIncident(v))
	assert.True(t, e.IsIncident(u))
}

func TestEdgeCreationWithNil(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	assert.Panics(t, func() { v.CreateEdge(nil) })
}

func TestLinkCreation(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	e := v.CreateEdge(u)
	l := e.AsLink(u)
	assert.Equal(t, u, l.Source())
	assert.Equal(t, v, l.Target())
}

func TestLoopEdgeCreation(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	assert.Panics(t, func() { v.CreateEdge(v) }, "Loop edges should not be created.")
	assert.Equal(t, 0, v.Degree(), "Loop edges should not be created.")
}

func TestEdgeGetOpposite(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	e := v.CreateEdge(u)
	assert.Equal(t, v, e.Opposite(u))
	assert.Equal(t, u, e.Opposite(v))
}

func TestGraphFormattingSimple(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	expected := fmt.Sprintf("Graph\n{\n  Vertex #%d\n}", v.index)
	assert.Equal(t, expected, g.String())
}

func TestEdgeFormattingWithoutMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	e := v.CreateEdge(u)
	assert.Equal(t, fmt.Sprintf("Edge #%d---#%d", v.index, u.index), e.String())
}

func TestArcFormattingWithoutMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	a := v.CreateArc(u)
	assert.Equal(t, fmt.Sprintf("Arc #%d-->#%d", v.index, u.index), a.String())
}

func TestEdgeFormattingWithMaps(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	e := v.CreateEdge(u)
	CreateEdgeMap(g, "weight", 1.3)
	CreateEdgeMap(g, "color", "blue")
	assert.Equal(t, fmt.Sprintf("Edge #%d---#%d (weight=1.3;color=blue;)", v.index, u.index), e.String())
}

func TestArcFormattingWithMaps(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	a := v.CreateArc(u)
	m := CreateArcMap(g, "weight", 1.3)
	CreateArcMap(g, "color", "blue")
	assert.Equal(t, fmt.Sprintf("Arc #%d-->#%d (weight=1.3;color=blue;)", v.index, u.index), a.String())
	assert.Equal(t, g, m.Graph(), "arcmap should identify its graph")
}

func TestVertexFormattingWithoutEdgesAndMaps(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	assert.Equal(t, fmt.Sprintf("Vertex #%d", v.index), v.String())
}

func TestVertexFormattingWithoutEdgesWithMaps(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	m := CreateVertexMap(g, "size", 5)
	assert.Equal(t, fmt.Sprintf("Vertex #%d (size=5;)", v.index), v.String())
	assert.Equal(t, g, m.Graph(), "vertexmap should identify its graph")
}

func TestVertexMapDispose(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	h := CreateVertexMap(g, "hokus", 2.3)
	CreateVertexMap(g, "size", 5)
	p := CreateVertexMap(g, "pokus", "la")
	p.Dispose()
	h.Dispose()
	assert.Equal(t, fmt.Sprintf("Vertex #%d (size=5;)", v.index), v.String())
}

func TestVertexFormattingWithEdgesAndArcs(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	u.CreateArc(v)
	v.CreateArc(u)
	_ = v.CreateEdge(u)
	expected := fmt.Sprintf("Vertex #%d\n  {\n    Edge #%d---#%d\n    Arc #%d-->#%d\n    Arc #%d-->#%d\n  }", v.index, v.index, u.index, v.index, u.index, u.index, v.index)
	assert.Equal(t, expected, v.String())
}

func TestVertexFormattingWithLabelledEdges(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	_ = v.CreateEdge(u)
	CreateEdgeMap(g, "height", 2.2)
	expected := fmt.Sprintf("Vertex #%d\n  {\n    Edge #%d---#%d (height=2.2;)\n  }", v.index, v.index, u.index)
	assert.Equal(t, expected, v.String())
}

func TestEdgeMapDisposal(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	_ = v.CreateEdge(u)
	em := CreateEdgeMap(g, "xxx", 2.2)
	CreateEdgeMap(g, "height", 2.2)
	em2 := CreateEdgeMap(g, "yyy", 2.2)
	em.Dispose()
	em2.Dispose()
	expected := fmt.Sprintf("Vertex #%d\n  {\n    Edge #%d---#%d (height=2.2;)\n  }", v.index, v.index, u.index)
	assert.Equal(t, expected, v.String())
}

func TestVertexFormattingWithLabelledArcs(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	v.CreateArc(u)
	CreateArcMap(g, "height", 2.2)
	expected := fmt.Sprintf("Vertex #%d\n  {\n    Arc #%d-->#%d (height=2.2;)\n  }", v.index, v.index, u.index)
	assert.Equal(t, expected, v.String())
}

func TestEdgeDisposal(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	v.CreateArc(u)
	em := CreateArcMap(g, "xxx", 2.2)
	CreateArcMap(g, "height", 2.2)
	em2 := CreateArcMap(g, "yyy", 2.2)
	em.Dispose()
	em2.Dispose()
	expected := fmt.Sprintf("Vertex #%d\n  {\n    Arc #%d-->#%d (height=2.2;)\n  }", v.index, v.index, u.index)
	assert.Equal(t, expected, v.String())
}

func TestArcMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	x := g.CreateVertex()
	em1 := CreateArcMap(g, "lbl", 17.3)
	assert.Equal(t, "lbl", em1.Label())
	a1 := v.CreateArc(x)
	em2 := CreateArcMap(g, "lbl2", 1.3)
	assert.Equal(t, "lbl2", em2.Label())
	a2 := u.CreateArc(v)
	assert.Equal(t, 17.3, em1.Get(a1))
	assert.Equal(t, 17.3, em1.Get(a2))
	assert.Equal(t, 1.3, em2.Get(a1))
	assert.Equal(t, 1.3, em2.Get(a2))
	em1.Set(a2, 88.8)
	assert.Equal(t, 88.8, em1.Get(a2))
}

func TestEdgeMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	x := g.CreateVertex()
	em1 := CreateEdgeMap(g, "lbl", 17.3)
	assert.Equal(t, "lbl", em1.Label())
	assert.Equal(t, g, em1.Graph())
	a1 := v.CreateEdge(x)
	em2 := CreateEdgeMap(g, "lbl2", 1.3)
	assert.Equal(t, "lbl2", em2.Label())
	a2 := u.CreateEdge(v)
	assert.Equal(t, 17.3, em1.Get(a1))
	assert.Equal(t, 17.3, em1.Get(a2))
	assert.Equal(t, 1.3, em2.Get(a1))
	assert.Equal(t, 1.3, em2.Get(a2))
	em1.Set(a2, 88.8)
	assert.Equal(t, 88.8, em1.Get(a2))
}

func TestVertexMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	vm1 := CreateVertexMap(g, "ll", 77)
	u := g.CreateVertex()
	vm2 := CreateVertexMap(g, "lx", 7)
	assert.Equal(t, "ll", vm1.Label())
	assert.Equal(t, 77, vm1.Get(v))
	assert.Equal(t, 77, vm1.Get(u))
	assert.Equal(t, 7, vm2.Get(v))
	assert.Equal(t, 7, vm2.Get(u))
	vm1.Set(u, 88)
	assert.Equal(t, 88, vm1.Get(u))
}

func TestVertexVisitingEmptyGraph(t *testing.T) {
	g := NewGraph()
	i := 0
	g.ForVertices(func(v *Vertex) { i++ })
	assert.Equal(t, 0, i)
}

func TestVertexVisitingFullGraph(t *testing.T) {
	g := NewGraph()
	g.CreateVertex()
	g.CreateVertex()
	g.CreateVertex()
	i, i2 := 0, 0
	g.ForVertices(func(v *Vertex) {
		i++
		i2 += v.index
	})
	assert.Equal(t, 3, i)
	assert.Equal(t, 3, i2)
}

func TestArcVisitingEmptyGraph(t *testing.T) {
	g := NewGraph()
	i := 0
	g.ForArcs(func(a *Arc) { i++ })
	assert.Equal(t, 0, i)
}

func TestArcVisitingFullGraph(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	w := g.CreateVertex()
	v.CreateArc(u)
	u.CreateArc(w)
	w.CreateArc(v)
	w.CreateArc(u)
	i, i2 := 0, 0
	g.ForArcs(func(a *Arc) {
		i++
		i2 += a.index
	})
	assert.Equal(t, 4, i)
	assert.Equal(t, 6, i2)
}

func TestEdgeVisitingEmptyGraph(t *testing.T) {
	g := NewGraph()
	i := 0
	g.ForEdges(func(e *Edge) { i++ })
	assert.Equal(t, 0, i)
}

func TestEdgeVisitingFullGraph(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	w := g.CreateVertex()
	_ = v.CreateEdge(u)
	_ = u.CreateEdge(w)
	i, i2 := 0, 0
	g.ForEdges(func(e *Edge) {
		i++
		i2 += e.index
	})
	assert.Equal(t, 2, i)
	assert.Equal(t, 1, i2)
}

func TestIncidentEdgeVisiting(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	w := g.CreateVertex()
	_ = v.CreateEdge(u)
	_ = u.CreateEdge(w)
	_ = v.CreateEdge(w)
	i, i2 := 0, 0
	v.ForEdges(func(e *Edge) {
		i++
		i2 += e.index
	})
	assert.Equal(t, 2, i)
	assert.Equal(t, 2, i2)
}

func TestIncidentArcVisiting(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	w := g.CreateVertex()
	v.CreateArc(u)
	u.CreateArc(w)
	v.CreateArc(w)
	i, i2 := 0, 0
	v.ForOutArcs(func(a *Arc) {
		i++
		i2 += a.index
	})
	assert.Equal(t, 2, i)
	assert.Equal(t, 2, i2)
	j, j2 := 0, 0
	w.ForInArcs(func(a *Arc) {
		j++
		j2 += a.index
	})
	assert.Equal(t, 2, j)
	assert.Equal(t, 3, j2)
}

func TestContainsEdgesOrArcs(t *testing.T) {
	g := NewGraph()
	assert.False(t, g.ContainsEdgesOrArcs())
	a := g.CreateVertex()
	b := g.CreateVertex()
	assert.False(t, g.ContainsEdgesOrArcs())
	a.CreateArc(b)
	assert.True(t, g.ContainsEdgesOrArcs())

	g = NewGraph()
	a = g.CreateVertex()
	b = g.CreateVertex()
	_ = a.CreateEdge(b)
	assert.True(t, g.ContainsEdgesOrArcs())
	c := g.CreateVertex()
	c.CreateArc(b)
	assert.True(t, g.ContainsEdgesOrArcs())
}

func TestPanicsOnBadCallToOpposite(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	e := u.CreateEdge(v)
	assert.Panics(t, func() { e.Opposite(w) }, "Opposite should panic on unknown vertex.")
}

func TestAnyVertex(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	assert.True(t, g.AnyVertex() == u || g.AnyVertex() == v || g.AnyVertex() == w)
}

func TestAnyVertexPanics(t *testing.T) {
	g := NewGraph()
	assert.Panics(t, func() { g.AnyVertex() })
}

func TestVertexSliceCreation(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	assert.Equal(t, 3, len(g.Vertices()))
	assert.True(t, slices.Contains(g.Vertices(), u))
	assert.True(t, slices.Contains(g.Vertices(), v))
	assert.True(t, slices.Contains(g.Vertices(), w))
}

func TestEdgeSliceCreation(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	e1 := u.CreateEdge(v)
	e2 := w.CreateEdge(v)
	assert.Equal(t, 2, len(v.Edges()))
	assert.True(t, slices.Contains(v.Edges(), e1))
	assert.True(t, slices.Contains(v.Edges(), e2))
}

func TestEdgeRemoval(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	z := g.CreateVertex()
	e1 := u.CreateEdge(v)
	e2 := w.CreateEdge(v)
	e3 := z.CreateEdge(w)
	e4 := z.CreateEdge(v)
	m := CreateEdgeMap(g, "test", 2.2)
	m.Set(e1, 1.0)
	m.Set(e2, 2.0)
	m.Set(e3, 3.0)
	m.Set(e4, 4.0)
	g.RemoveEdge(e2)
	assert.Equal(t, 3, g.EdgeCount())
	assert.Equal(t, 2, v.Degree())
	assert.Equal(t, 1, w.Degree())
	assert.Equal(t, 1.0, m.Get(e1))
	assert.Equal(t, 3.0, m.Get(e3))
	assert.Equal(t, 4.0, m.Get(e4))
	assert.Equal(t, 3, len(m.data))
	assert.True(t, slices.Contains(g.edges, e1))
	assert.False(t, slices.Contains(g.edges, e2))
	assert.True(t, slices.Contains(g.edges, e3))
	assert.True(t, slices.Contains(g.edges, e4))
}

// TestEdgeRemovalLastEdgeOfFirstVertex checks that removing an edge does not
// panic when it is the only (last) edge of vertices[0] (the a-side).
func TestEdgeRemovalLastEdgeOfFirstVertex(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	e1 := u.CreateEdge(v) // u.edges = [e1], v.edges = [e1, e2]
	v.CreateEdge(w)       // v gets a second edge so the b-side is safe
	// removing e1: a=u has only e1 (last edge → a-side would panic), b=v has two edges (safe)
	assert.NotPanics(t, func() { g.RemoveEdge(e1) })
	assert.Equal(t, 1, g.EdgeCount())
}

// TestEdgeRemovalLastEdgeOfSecondVertex checks that removing an edge does not
// panic when it is the only (last) edge of vertices[1] (the b-side).
func TestEdgeRemovalLastEdgeOfSecondVertex(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	e1 := u.CreateEdge(v) // u.edges = [e1, e2], v.edges = [e1]
	u.CreateEdge(w)       // u gets a second edge so the a-side is safe
	// removing e1: a=u has two edges (safe), b=v has only e1 (last edge → b-side would panic)
	assert.NotPanics(t, func() { g.RemoveEdge(e1) })
	assert.Equal(t, 1, g.EdgeCount())
}

func TestMultipleEdgeRemoval(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	z := g.CreateVertex()
	e1 := u.CreateEdge(v)
	e2 := w.CreateEdge(v)
	e3 := z.CreateEdge(w)
	e4 := z.CreateEdge(v)
	m := CreateEdgeMap(g, "test", 2.2)
	m.Set(e1, 1.0)
	m.Set(e2, 2.0)
	m.Set(e3, 3.0)
	m.Set(e4, 4.0)
	g.RemoveEdge(e2)
	g.RemoveEdge(e4)
	assert.Equal(t, 2, g.EdgeCount())
	assert.Equal(t, 1, v.Degree())
	assert.Equal(t, 1, w.Degree())
	assert.Equal(t, 1, z.Degree())
	assert.Equal(t, 1.0, m.Get(e1))
	assert.Equal(t, 3.0, m.Get(e3))
	assert.Equal(t, 2, len(m.data))
	assert.True(t, slices.Contains(g.edges, e1))
	assert.False(t, slices.Contains(g.edges, e2))
	assert.True(t, slices.Contains(g.edges, e3))
	assert.False(t, slices.Contains(g.edges, e4))
}

func TestArcRemoval(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	z := g.CreateVertex()
	e1 := u.CreateArc(v)
	e2 := w.CreateArc(v)
	e3 := z.CreateArc(w)
	e4 := z.CreateArc(v)
	m := CreateArcMap(g, "test", 2.2)
	m.Set(e1, 1.0)
	m.Set(e2, 2.0)
	m.Set(e3, 3.0)
	m.Set(e4, 4.0)
	g.RemoveArc(e2)
	assert.Equal(t, 3, g.ArcCount())
	assert.Equal(t, 2, v.InDegree())
	assert.Equal(t, 0, w.OutDegree())
	assert.Equal(t, 1.0, m.Get(e1))
	assert.Equal(t, 3.0, m.Get(e3))
	assert.Equal(t, 4.0, m.Get(e4))
	assert.Equal(t, 3, len(m.data))
	assert.True(t, slices.Contains(g.arcs, e1))
	assert.False(t, slices.Contains(g.arcs, e2))
	assert.True(t, slices.Contains(g.arcs, e3))
	assert.True(t, slices.Contains(g.arcs, e4))
}

func TestMultiArcRemoval(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	z := g.CreateVertex()
	e1 := u.CreateArc(v)
	e2 := w.CreateArc(v)
	e3 := z.CreateArc(w)
	e4 := z.CreateArc(v)
	m := CreateArcMap(g, "test", 2.2)
	m.Set(e1, 1.0)
	m.Set(e2, 2.0)
	m.Set(e3, 3.0)
	m.Set(e4, 4.0)
	g.RemoveArc(e2)
	g.RemoveArc(e4)
	assert.Equal(t, 2, g.ArcCount())
	assert.Equal(t, 1, v.InDegree())
	assert.Equal(t, 1, z.OutDegree())
	assert.Equal(t, 1.0, m.Get(e1))
	assert.Equal(t, 3.0, m.Get(e3))
	assert.Equal(t, 2, len(m.data))
	assert.True(t, slices.Contains(g.arcs, e1))
	assert.False(t, slices.Contains(g.arcs, e2))
	assert.True(t, slices.Contains(g.arcs, e3))
	assert.False(t, slices.Contains(g.arcs, e4))
}

func TestMultiArcRemoval2(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	z := g.CreateVertex()
	e1 := v.CreateArc(u)
	e2 := v.CreateArc(w)
	e3 := w.CreateArc(z)
	e4 := v.CreateArc(z)
	m := CreateArcMap(g, "test", 2.2)
	m.Set(e1, 1.0)
	m.Set(e2, 2.0)
	m.Set(e3, 3.0)
	m.Set(e4, 4.0)
	g.RemoveArc(e2)
	g.RemoveArc(e4)
	assert.Equal(t, 2, g.ArcCount())
	assert.Equal(t, 1, v.OutDegree())
	assert.Equal(t, 1, z.InDegree())
	assert.Equal(t, 1.0, m.Get(e1))
	assert.Equal(t, 3.0, m.Get(e3))
	assert.Equal(t, 2, len(m.data))
	assert.True(t, slices.Contains(g.arcs, e1))
	assert.False(t, slices.Contains(g.arcs, e2))
	assert.True(t, slices.Contains(g.arcs, e3))
	assert.False(t, slices.Contains(g.arcs, e4))
}

func TestVertexRemoval(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	z := g.CreateVertex()
	v.CreateArc(u)
	v.CreateArc(w)
	w.CreateEdge(z)
	v.CreateEdge(z)
	u.CreateArc(v)
	m := CreateVertexMap(g, "test", 2.2)
	m.Set(u, 1.1)
	m.Set(v, 2.1)
	m.Set(w, 3.1)
	m.Set(z, 4.1)

	g.RemoveVertex(w)
	assert.Equal(t, 3, g.VertexCount())
	assert.True(t, slices.Contains(g.Vertices(), u))
	assert.True(t, slices.Contains(g.Vertices(), v))
	assert.False(t, slices.Contains(g.Vertices(), w))
	assert.True(t, slices.Contains(g.Vertices(), z))
	assert.Equal(t, 1, u.OutDegree())
	assert.Equal(t, 1, v.InDegree())
	assert.Equal(t, 1, z.Degree())
	assert.Equal(t, 2, g.ArcCount())
	assert.Equal(t, 1, g.EdgeCount())
	assert.Equal(t, 2.1, m.Get(v))
	assert.Equal(t, 1.1, m.Get(u))
	assert.Equal(t, 4.1, m.Get(z))
	g.RemoveVertex(z)
	g.RemoveVertex(v)
}

func TestEdgeLinkConversion(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	e := u.CreateEdge(v)
	l := e.AsLink(v)
	assert.True(t, LinkIsEdge(l))
	assert.False(t, LinkIsArc(l))
	LinkGetEdge(l)
	assert.Panics(t, func() { LinkGetArc(l) })
}

func TestArcLinkConversion(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	var l Link = u.CreateArc(v)
	assert.False(t, LinkIsEdge(l))
	assert.True(t, LinkIsArc(l))
	LinkGetArc(l)
	assert.Panics(t, func() { LinkGetEdge(l) })
}

func TestReverseArcLinkConversion(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	arc := v.CreateArc(u)
	reverseLink := arc.AsReverseLink()
	assert.False(t, LinkIsEdge(reverseLink))
	assert.True(t, LinkIsReverseArc(reverseLink))
	assert.Equal(t, u, reverseLink.Source())
	assert.Equal(t, v, reverseLink.Target())
}

func TestReverseArcLinkGetArc(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	arc := v.CreateArc(u)
	reverseLink := arc.AsReverseLink()
	assert.Equal(t, arc, ReverseLinkGetArc(reverseLink))
}

func TestReverseArcLinkPanics(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	var link Link = v.CreateArc(u)
	assert.Panics(t, func() { ReverseLinkGetArc(link) })
}
