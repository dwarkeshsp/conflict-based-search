package main

import "sync"

type CTNode struct {
	constaints map[Agent]map[int]*Point
	solution   map[Agent][]*Point
	cost       int
}

func (n *CTNode) storeCost() {
	cost := 0
	for _, points := range n.solution {
		cost += len(points)
	}
	n.cost = cost
}

type CT struct {
	queue []*CTNode
	mu    *sync.Mutex
}

func (t CT) Len() int { return len(t.queue) }

func (t CT) Less(i, j int) bool { return t.queue[i].cost < t.queue[j].cost }

func (t CT) Swap(i, j int) {
	t.mu.Lock()
	t.queue[i], t.queue[j] = t.queue[j], t.queue[i]
	t.mu.Unlock()
}

func (t *CT) Push(x interface{}) {
	t.mu.Lock()
	node := x.(*CTNode)
	t.queue = append(t.queue, node)
	t.mu.Unlock()
}

func (t *CT) Pop() interface{} {
	t.mu.Lock()
	old := t.queue
	n := len(old)
	node := old[n-1]
	old[n-1] = nil // avoid memory leak
	t.queue = old[0 : n-1]
	t.mu.Unlock()
	return node
}
