package hmgraph

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"testing"
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
	assert.Equal(t, g, v.GetGraph(), "vertex should idetify graph")
	assert.Equal(t, 0, v.GetIndex(), "single vertex's index should be 0")
	assert.True(t, v.InDegree() == 0 && v.OutDegree() == 0 && v.Degree() == 0,
		"New graph shouldn't have any edges or arcs.")
}

func TestArcCreation(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	a := v.CreateArc(u)
	assert.Equal(t, 0, a.GetIndex(), "single node's index should be 0")

	if v.InDegree() != 0 || v.OutDegree() != 1 || u.InDegree() != 1 || u.OutDegree() != 0 {
		t.Errorf("Created arc should contribute to outdeg of source and indeg of target.")
	}
	if a.Source() != v || a.Target() != u {
		t.Errorf("Created arc doesn't identify source and target correctly.")
	}
}

func TestEdgeCreation(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	e := v.CreateEdge(u)
	assert.Equal(t, 0, e.GetIndex(), "first edge should have index 0")
	if v.Degree() != 1 || u.Degree() != 1 {
		t.Errorf("Created edge should contribute to deg of connected edges.")
	}
	if (e.Vertices()[0] != v && e.Vertices()[1] != v) || (e.Vertices()[0] != u && e.Vertices()[1] != u) {
		t.Errorf("Created edge doesn't identify vertices correctly.")
	}
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
	assert.Equal(t, l.Source(), u)
	assert.Equal(t, l.Target(), v)
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
	if e.Opposite(u) != v || e.Opposite(v) != u {
		t.Errorf("Created edge doesn't identify opposites correctly.")
	}
}

func TestGraphFormattingSimple(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	expected := fmt.Sprintf("Graph\n{\n  Vertex #%d\n}", v.index)
	if g.String() != expected {
		t.Errorf("Graph not correctly formatted: \n>>%s<<, expected:\n>>%s<<", g.String(), expected)
	}
}

func TestEdgeFormattingWithoutMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	e := v.CreateEdge(u)
	if e.String() != fmt.Sprintf("Edge #%d---#%d", v.index, u.index) {
		t.Errorf("Edge not correctly formatted: >>%s<<", e.String())
	}
}

func TestArcFormattingWithoutMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	a := v.CreateArc(u)
	assert.Equal(t, fmt.Sprintf("Arc #%d-->#%d", v.index, u.index), a.String(),
		"Edge not correctly formatted: >>%s<<", a.String())
}

func TestEdgeFormattingWithMaps(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	e := v.CreateEdge(u)
	CreateEdgeMap(g, "weight", 1.3)
	CreateEdgeMap(g, "color", "blue")
	if e.String() != fmt.Sprintf("Edge #%d---#%d (weight=1.3;color=blue;)", v.index, u.index) {
		t.Errorf("Edge not correctly formatted: >>%s<<", e.String())
	}
}
func TestArcFormattingWithMaps(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	a := v.CreateArc(u)
	m := CreateArcMap(g, "weight", 1.3)
	CreateArcMap(g, "color", "blue")
	if a.String() != fmt.Sprintf("Arc #%d-->#%d (weight=1.3;color=blue;)", v.index, u.index) {
		t.Errorf("Edge not correctly formatted: >>%s<<", a.String())
	}
	assert.Equal(t, g, m.GetGraph(), "arcmap should identify its graph")
}

func TestVertexFormattingWithoutEdgesAndMaps(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	if v.String() != fmt.Sprintf("Vertex #%d", v.index) {
		t.Errorf("Edge not correctly formatted: >>%s<<", v.String())
	}
}

func TestVertexFormattingWithoutEdgesWithMaps(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	m := CreateVertexMap(g, "size", 5)
	if v.String() != fmt.Sprintf("Vertex #%d (size=5;)", v.index) {
		t.Errorf("Edge not correctly formatted: >>%s<<", v.String())
	}
	assert.Equal(t, g, m.GetGraph(), "vertexmap should identify its graph")
}

func TestVertexMapDispose(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	h := CreateVertexMap(g, "hokus", 2.3)
	CreateVertexMap(g, "size", 5)
	p := CreateVertexMap(g, "pokus", "la")
	p.Dispose()
	h.Dispose()
	if v.String() != fmt.Sprintf("Vertex #%d (size=5;)", v.index) {
		t.Errorf("Edge not correctly formatted: >>%s<<", v.String())
	}
}

func TestVertexFormattingWithEdgesAndArcs(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	u.CreateArc(v)
	v.CreateArc(u)
	_ = v.CreateEdge(u)
	expected := fmt.Sprintf("Vertex #%d\n  {\n    Edge #%d---#%d\n    Arc #%d-->#%d\n    Arc #%d-->#%d\n  }", v.index, v.index, u.index, v.index, u.index, u.index, v.index)
	assert.Equal(t, expected, v.String(), "Edge not correctly formatted: >>%s<<\n expected: >>%s<<", v.String(), expected)
}

func TestVertexFormattingWithLabelledEdges(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	_ = v.CreateEdge(u)
	CreateEdgeMap(g, "height", 2.2)
	expected := fmt.Sprintf("Vertex #%d\n  {\n    Edge #%d---#%d (height=2.2;)\n  }", v.index, v.index, u.index)
	if v.String() != expected {
		t.Errorf("Edge not correctly formatted: >>%s<<\n expected: >>%s<<", v.String(), expected)
	}
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
	if v.String() != expected {
		t.Errorf("Edge not correctly formatted: >>%s<<\n expected: >>%s<<", v.String(), expected)
	}
}

func TestVertexFormattingWithLabelledArcs(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	v.CreateArc(u)
	CreateArcMap(g, "height", 2.2)
	expected := fmt.Sprintf("Vertex #%d\n  {\n    Arc #%d-->#%d (height=2.2;)\n  }", v.index, v.index, u.index)
	assert.Equal(t, expected, v.String(),
		"Edge not correctly formatted: >>%s<<\n expected: >>%s<<", v.String(), expected)
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
	if v.String() != expected {
		t.Errorf("Edge not correctly formatted: >>%s<<\n expected: >>%s<<", v.String(), expected)
	}
}

func TestArcMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	x := g.CreateVertex()
	em1 := CreateArcMap(g, "lbl", 17.3)
	if em1.GetLabel() != "lbl" {
		t.Errorf("Arcmap returns wrong label")
	}
	a1 := v.CreateArc(x)
	em2 := CreateArcMap(g, "lbl2", 1.3)
	if em2.GetLabel() != "lbl2" {
		t.Errorf("Arcmap returns wrong label")
	}
	a2 := u.CreateArc(v)
	if em1.Get(a1) != 17.3 || em1.Get(a2) != 17.3 {
		t.Errorf("ArcMap returns wrong default %f, %f, expected %f", em1.Get(a1), em1.Get(a2), 17.3)
	}
	if em2.Get(a2) != 1.3 || em2.Get(a2) != 1.3 {
		t.Errorf("ArcMap returns wrong default %f, %f, expected %f", em1.Get(a1), em1.Get(a2), 1.3)
	}
	em1.Set(a2, 88.8)
	if em1.Get(a2) != 88.8 {
		t.Errorf("ArcMap returns wrong value %f, expected %f", em1.Get(a2), 88.8)
	}
}

func TestEdgeMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	x := g.CreateVertex()
	em1 := CreateEdgeMap(g, "lbl", 17.3)
	assert.Equal(t, "lbl", em1.GetLabel(),
		"EdgeMap returns wrong label")
	assert.Equal(t, g, em1.GetGraph())
	a1 := v.CreateEdge(x)
	em2 := CreateEdgeMap(g, "lbl2", 1.3)
	if em2.GetLabel() != "lbl2" {
		t.Errorf("EdgeMap returns wrong label")
	}
	a2 := u.CreateEdge(v)
	if em1.Get(a1) != 17.3 || em1.Get(a2) != 17.3 {
		t.Errorf("EdgeMap returns wrong default %f, %f, expected %f", em1.Get(a1), em1.Get(a2), 17.3)
	}
	if em2.Get(a2) != 1.3 || em2.Get(a2) != 1.3 {
		t.Errorf("EdgeMap returns wrong default %f, %f, expected %f", em1.Get(a1), em1.Get(a2), 1.3)
	}
	em1.Set(a2, 88.8)
	if em1.Get(a2) != 88.8 {
		t.Errorf("EdgeMap returns wrong value %f, expected %f", em1.Get(a2), 88.8)
	}
}

func TestVertexMap(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	vm1 := CreateVertexMap(g, "ll", 77)
	u := g.CreateVertex()
	vm2 := CreateVertexMap(g, "lx", 7)
	if vm1.GetLabel() != "ll" {
		t.Errorf("VertexMap returns wrong label")
	}
	if vm1.Get(v) != 77 || vm1.Get(u) != 77 {
		t.Errorf("VertexMap returns wrong default %d, %d expected %d", vm1.Get(v), vm1.Get(u), 77)
	}
	if vm2.Get(v) != 7 || vm2.Get(u) != 7 {
		t.Errorf("VertexMap returns wrong default %d, %d, expected %d", vm1.Get(v), vm1.Get(u), 7)
	}
	vm1.Set(u, 88)
	if vm1.Get(u) != 88 {
		t.Errorf("VertexMap returns wrong value %d, expected %d", vm1.Get(u), 88)
	}
}

func TestVertexVisitingEmptyGraph(t *testing.T) {
	g := NewGraph()
	i := 0
	g.ForVertices(func(v *Vertex) { i++ })
	if i != 0 {
		t.Errorf("Shouldnt visit any vertices in empty graph.")
	}
}

func TestVertexVisitingFullGraph(t *testing.T) {
	g := NewGraph()
	g.CreateVertex()
	g.CreateVertex()
	g.CreateVertex()
	i := 0
	i2 := 0
	g.ForVertices(func(v *Vertex) {
		i++
		i2 += v.index
	})
	if i != 3 || i2 != 3 {
		t.Errorf("Didn't visit all vertices.")
	}
}

func TestArcVisitingEmptyGraph(t *testing.T) {
	g := NewGraph()
	i := 0
	g.ForArcs(func(a *Arc) { i++ })
	if i != 0 {
		t.Errorf("Shouldnt visit any arcs in empty graph.")
	}
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
	i := 0
	i2 := 0
	g.ForArcs(func(a *Arc) {
		i++
		i2 += a.index
	})
	if i != 4 || i2 != 6 {
		t.Errorf("Didn't visit all arcs.")
	}
}

func TestEdgeVisitingEmptyGraph(t *testing.T) {
	g := NewGraph()
	i := 0
	g.ForEdges(func(e *Edge) { i++ })
	if i != 0 {
		t.Errorf("Shouldnt visit any edges in empty graph.")
	}
}

func TestEdgeVisitingFullGraph(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	w := g.CreateVertex()
	_ = v.CreateEdge(u)
	_ = u.CreateEdge(w)
	i := 0
	i2 := 0
	g.ForEdges(func(e *Edge) {
		i++
		i2 += e.index
	})
	if i != 2 || i2 != 1 {
		t.Errorf("Didn't visit all edges.")
	}
}

func TestIncidentEdgeVisiting(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	w := g.CreateVertex()
	_ = v.CreateEdge(u)
	_ = u.CreateEdge(w)
	_ = v.CreateEdge(w)
	i := 0
	i2 := 0
	v.ForEdges(func(e *Edge) {
		i++
		i2 += e.index
	})
	if i != 2 || i2 != 2 {
		t.Errorf("Didn't visit all edges %d %d.", i, i2)
	}
}

func TestIncidentArcVisiting(t *testing.T) {
	g := NewGraph()
	v := g.CreateVertex()
	u := g.CreateVertex()
	w := g.CreateVertex()
	v.CreateArc(u)
	u.CreateArc(w)
	v.CreateArc(w)
	i := 0
	i2 := 0
	v.ForOutArcs(func(a *Arc) {
		i++
		i2 += a.index
	})
	if i != 2 || i2 != 2 {
		t.Errorf("Didn't visit all outArcs %d %d.", i, i2)
	}
	j := 0
	j2 := 0
	w.ForInArcs(func(a *Arc) {
		j++
		j2 += a.index
	})
	if j != 2 || j2 != 3 {
		t.Errorf("Didn't visit all inArcs %d %d.", j, j2)
	}

}

func TestContainsEdgesOrArcs(t *testing.T) {
	g := NewGraph()
	if g.ContainsEdgesOrArcs() {
		t.Errorf("Graph should not contain any edge or arcs but was: %v.", g.ContainsEdgesOrArcs())
	}

	a := g.CreateVertex()
	b := g.CreateVertex()
	if g.ContainsEdgesOrArcs() {
		t.Errorf("Graph should not contain any edge or arcs but was: %v.", g.ContainsEdgesOrArcs())
	}

	a.CreateArc(b)
	if !g.ContainsEdgesOrArcs() {
		t.Errorf("Graph should contain any edge or arcs but was: %v.", g.ContainsEdgesOrArcs())
	}

	g = NewGraph()
	a = g.CreateVertex()
	b = g.CreateVertex()
	_ = a.CreateEdge(b)
	if !g.ContainsEdgesOrArcs() {
		t.Errorf("Graph should contain any edge or arcs but was: %v.", g.ContainsEdgesOrArcs())
	}

	c := g.CreateVertex()
	c.CreateArc(b)
	if !g.ContainsEdgesOrArcs() {
		t.Errorf("Graph should contain any edge or arcs but was: %v.", g.ContainsEdgesOrArcs())
	}
}

func TestContainsLoops(t *testing.T) {
	g := NewGraph()
	assert.False(t, g.ContainsLoops())
	v1 := g.CreateVertex()
	v2 := g.CreateVertex()
	assert.False(t, g.ContainsLoops())
	_ = v1.CreateEdge(v2)
	assert.False(t, g.ContainsLoops())
	_ = v2.CreateArc(v1)
	assert.False(t, g.ContainsLoops())
	_ = v1.CreateArc(v1)
	assert.True(t, g.ContainsLoops())
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
	assert.True(t, g.GetAnyVertex() == u || g.GetAnyVertex() == v || g.GetAnyVertex() == w)
}

func TestAnyVertexPanics(t *testing.T) {
	g := NewGraph()
	assert.Panics(t, func() { g.GetAnyVertex() })
}

func TestVertexSliceCreation(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	assert.Equal(t, 3, len(g.GetVertices()))
	assert.True(t, slices.Contains(g.GetVertices(), u))
	assert.True(t, slices.Contains(g.GetVertices(), v))
	assert.True(t, slices.Contains(g.GetVertices(), w))
}

func TestEdgeSliceCreation(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	e1 := u.CreateEdge(v)
	e2 := w.CreateEdge(v)
	assert.Equal(t, 2, len(v.GetEdges()))
	assert.True(t, slices.Contains(v.GetEdges(), e1))
	assert.True(t, slices.Contains(v.GetEdges(), e2))
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
	assert.True(t, slices.Contains(g.GetVertices(), u))
	assert.True(t, slices.Contains(g.GetVertices(), v))
	assert.False(t, slices.Contains(g.GetVertices(), w))
	assert.True(t, slices.Contains(g.GetVertices(), z))
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
	assert.True(t, LinkIsEdge(&l))
	assert.False(t, LinkIsArc(&l))
	LinkGetEdge(&l)
	assert.Panics(t, func() { LinkGetArc(&l) })
}

func TestArcLinkConversion(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	var l Link = u.CreateArc(v)
	assert.False(t, LinkIsEdge(&l))
	assert.True(t, LinkIsArc(&l))
	LinkGetArc(&l)
	assert.Panics(t, func() { LinkGetEdge(&l) })
}

func TestReverseArcLinkConversion(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	arc := v.CreateArc(u)
	reverseLink := arc.AsReverseLink()
	assert.False(t, LinkIsEdge(&reverseLink))
	assert.True(t, LinkIsReverseArc(&reverseLink))
	assert.Equal(t, u, reverseLink.Source())
	assert.Equal(t, v, reverseLink.Target())
}

func TestReverseArcLinkGetArc(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	arc := v.CreateArc(u)
	reverseLink := arc.AsReverseLink()
	assert.Equal(t, arc, ReverseLinkGetArc(&reverseLink))
}

func TestReverseArcLinkPanics(t *testing.T) {
	g := NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	var link Link = v.CreateArc(u)
	assert.Panics(t, func() { ReverseLinkGetArc(&link) })
}
