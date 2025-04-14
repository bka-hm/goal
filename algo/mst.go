package algo

import (
	"gitlab.lrz.de/hm/goal-core/base"
	"gitlab.lrz.de/hm/goal-core/hmgraph"
)

type MstAlgo[T base.Number] func(g *hmgraph.Graph, weight *hmgraph.EdgeMap[T]) (mst []*hmgraph.Edge)
