package main

import (
	"sync"
)

type CTree []*CTNode

var mu sync.Mutex

func (t CTree) Len() int { return len(t) }

func (t CTree) Less(i, j int) bool { return t[i].cost < t[j].cost }

func (t CTree) Swap(i, j int) {
	mu.Lock()
	t[i], t[j] = t[j], t[i]
	mu.Unlock()
}

func (t *CTree) Push(x interface{}) {
	mu.Lock()
	node := x.(*CTNode)
	*t = append(*t, node)
	mu.Unlock()
}

func (t *CTree) Pop() interface{} {
	mu.Lock()
	old := *t
	n := len(old)
	node := old[n-1]
	old[n-1] = nil // avoid memory leak
	*t = old[0 : n-1]
	mu.Unlock()
	return node
}
