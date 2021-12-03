package main

import (
	"sync"

	"github.com/dwarkeshsp/astar"
)

type Agent struct {
	start astar.Node
	end   astar.Node
}

type CTNode struct {
	constaints map[Agent][]astar.Node
	solution   map[Agent][]astar.Node
	cost       int
}

func emptyNode(agents []Agent) CTNode {
	n := CTNode{}
	n.constaints = make(map[Agent][]astar.Node)
	for _, agent := range agents {
		n.constaints[agent] = []astar.Node{}
	}
	n.solution = make(map[Agent][]astar.Node)
	return n
}

func (n *CTNode) findSolution() {
	for agent, obstacles := range n.constaints {
		aConfig := astar.Config{
			GridWidth:    GRID_SIZE,
			GridHeight:   GRID_SIZE,
			InvalidNodes: obstacles,
		}

		algo, _ := astar.New(aConfig)
		agentPath, _ := algo.FindPath(agent.start, agent.end)
		n.solution[agent] = agentPath
	}
	n.storeCost()
}

func (n *CTNode) storeCost() {
	cost := 0
	for _, points := range n.solution {
		cost += len(points)
	}
	n.cost = cost
}

type CTree struct {
	queue []*CTNode
	mu    *sync.Mutex
}

func (t CTree) Len() int { return len(t.queue) }

func (t CTree) Less(i, j int) bool { return t.queue[i].cost < t.queue[j].cost }

func (t CTree) Swap(i, j int) {
	t.mu.Lock()
	t.queue[i], t.queue[j] = t.queue[j], t.queue[i]
	t.mu.Unlock()
}

func (t *CTree) Push(node *CTNode) {
	t.mu.Lock()
	t.queue = append(t.queue, node)
	t.mu.Unlock()
}

func (t *CTree) Pop() *CTNode {
	t.mu.Lock()
	old := t.queue
	n := len(old)
	node := old[n-1]
	old[n-1] = nil // avoid memory leak
	t.queue = old[0 : n-1]
	t.mu.Unlock()
	return node
}
