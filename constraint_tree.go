package main

import (
	"sync"
)

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

func (t *CTree) Push(x interface{}) {
	t.mu.Lock()
	node := x.(*CTNode)
	t.queue = append(t.queue, node)
	t.mu.Unlock()
}

func (t *CTree) Pop() interface{} {
	t.mu.Lock()
	old := t.queue
	n := len(old)
	node := old[n-1]
	old[n-1] = nil // avoid memory leak
	t.queue = old[0 : n-1]
	t.mu.Unlock()
	return node
}
