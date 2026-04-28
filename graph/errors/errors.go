package errors

import "gitlab.lrz.de/hm/goal/graph/hmgraph"

type ContainsEdgeError struct {
	containedEdge *hmgraph.Edge
}

func NewContainsEdgeError(edge *hmgraph.Edge) *ContainsEdgeError {
	return &ContainsEdgeError{containedEdge: edge}
}

func (ee ContainsEdgeError) Error() string {
	return ee.containedEdge.String()
}

func (ee ContainsEdgeError) Edge() *hmgraph.Edge {
	return ee.containedEdge
}

type ContainsCycleError struct {
	containedCycle []*hmgraph.Arc
}

func NewContainsCycleError(cycle []*hmgraph.Arc) *ContainsCycleError {
	return &ContainsCycleError{containedCycle: cycle}
}

func (ee ContainsCycleError) Error() string {
	toReturn := "Cycle: ["
	for i, a := range ee.containedCycle {
		if i > 0 {
			toReturn += ", "
		}
		toReturn += a.String()
	}
	return toReturn + "]"
}

func (ee ContainsCycleError) Cycle() []*hmgraph.Arc {
	return ee.containedCycle
}
