package taskflow

import (
	"fmt"
	"github.com/CloudGoSight/cloudgosight_infrastructure/taskflow/utils"
	"sync"
	"sync/atomic"
)

type eGraph struct { // execution graph
	name          string
	nodes         []*innerNode
	joinCounter   utils.RC
	entries       []*innerNode
	scheCond      *sync.Cond
	instancelized bool
	canceled      atomic.Bool // only changes when task in graph panic
}

func newGraph(name string) *eGraph {
	return &eGraph{
		name:     name,
		nodes:    make([]*innerNode, 0),
		scheCond: sync.NewCond(&sync.Mutex{}),
	}
}

func (g *eGraph) JoinCounter() int {
	return g.joinCounter.Value()
}

func (g *eGraph) reset() {
	g.joinCounter.Set(0)
	g.entries = g.entries[:0]
	for _, n := range g.nodes {
		n.joinCounter.Set(0)
	}
}

func (g *eGraph) push(n ...*innerNode) {
	g.nodes = append(g.nodes, n...)
	for _, node := range n {
		node.g = g
	}
}

func (g *eGraph) setup() {
	g.reset()

	for _, node := range g.nodes {
		node.joinCounter.Set(len(node.dependents))

		if len(node.dependents) == 0 {
			g.entries = append(g.entries, node)
		}
	}
}

// only for visualizer
func (g *eGraph) topologicalSort() (sorted []*innerNode, err error) {
	indegree := map[*innerNode]int{} // Node -> indegree
	zeros := make([]*innerNode, 0)   // zero deps
	sorted = make([]*innerNode, 0, len(g.nodes))

	for _, node := range g.nodes {
		set := map[*innerNode]struct{}{}
		for _, dep := range node.dependents {
			set[dep] = struct{}{}
		}
		indegree[node] = len(set)
		if len(set) == 0 {
			zeros = append(zeros, node)
		}
	}

	for len(zeros) > 0 {
		node := zeros[0]
		zeros = zeros[1:]
		sorted = append(sorted, node)

		for _, succeesor := range node.successors {
			in := indegree[succeesor]
			in = in - 1
			if in <= 0 { // successor has no deps, put into zeros list
				zeros = append(zeros, succeesor)
			}
			indegree[succeesor] = in
		}
	}

	for _, node := range g.nodes {
		if indegree[node] > 0 {
			return nil, fmt.Errorf("graph has cycles")
		}
	}

	return
}
