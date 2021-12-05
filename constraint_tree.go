package main

import (
	"sync"
)

type CTree []*CTNode

var mu sync.Mutex

func (t CTree) Len() int { return len(t) }

func (t CTree) Less(i, j int) bool { return t[i].cost < t[j].cost }

func (t CTree) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t *CTree) Push(x interface{}) {
	node := x.(*CTNode)
	*t = append(*t, node)
}

func (t *CTree) Pop() interface{} {
	old := *t
	n := len(old)
	node := old[n-1]
	old[n-1] = nil // avoid memory leak
	*t = old[0 : n-1]
	return node
}
