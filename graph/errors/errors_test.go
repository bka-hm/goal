package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	graph_errors "gitlab.lrz.de/hm/goal/graph/errors"
	"gitlab.lrz.de/hm/goal/graph/hmgraph"
)

func TestContainsEdgeErrorEdge(t *testing.T) {
	g := hmgraph.NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	e := u.CreateEdge(v)
	err := graph_errors.NewContainsEdgeError(e)
	assert.Equal(t, e, err.Edge())
}

func TestContainsEdgeErrorMessage(t *testing.T) {
	g := hmgraph.NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	e := u.CreateEdge(v)
	err := graph_errors.NewContainsEdgeError(e)
	assert.Equal(t, e.String(), err.Error())
}

func TestContainsEdgeErrorImplementsError(t *testing.T) {
	g := hmgraph.NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	e := u.CreateEdge(v)
	var err error = graph_errors.NewContainsEdgeError(e)
	assert.NotNil(t, err)
}

func TestContainsCycleErrorCycle(t *testing.T) {
	g := hmgraph.NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	a1 := u.CreateArc(v)
	a2 := v.CreateArc(w)
	a3 := w.CreateArc(u)
	cycle := []*hmgraph.Arc{a1, a2, a3}
	err := graph_errors.NewContainsCycleError(cycle)
	assert.Equal(t, cycle, err.Cycle())
}

func TestContainsCycleErrorMessageSingleArc(t *testing.T) {
	g := hmgraph.NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	a := u.CreateArc(v)
	err := graph_errors.NewContainsCycleError([]*hmgraph.Arc{a})
	assert.Equal(t, fmt.Sprintf("Cycle: [%s]", a.String()), err.Error())
}

func TestContainsCycleErrorMessageMultipleArcs(t *testing.T) {
	g := hmgraph.NewGraph()
	u := g.CreateVertex()
	v := g.CreateVertex()
	w := g.CreateVertex()
	a1 := u.CreateArc(v)
	a2 := v.CreateArc(w)
	a3 := w.CreateArc(u)
	err := graph_errors.NewContainsCycleError([]*hmgraph.Arc{a1, a2, a3})
	expected := fmt.Sprintf("Cycle: [%s, %s, %s]", a1.String(), a2.String(), a3.String())
	assert.Equal(t, expected, err.Error())
}

func TestContainsCycleErrorMessageEmptyCycle(t *testing.T) {
	err := graph_errors.NewContainsCycleError([]*hmgraph.Arc{})
	assert.Equal(t, "Cycle: []", err.Error())
}

func TestContainsCycleErrorImplementsError(t *testing.T) {
	var err error = graph_errors.NewContainsCycleError([]*hmgraph.Arc{})
	assert.NotNil(t, err)
}
