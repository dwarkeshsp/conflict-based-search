package main

import (
	"container/heap"

	"github.com/dwarkeshsp/astar"
)

func plan(agents []Agent, obstacles []astar.Node) map[Agent][]astar.Node {
	open := &CTree{queue: make([]*CTNode, 100)}
	heap.Init(open)

	root := createRootNode(agents, obstacles)
	heap.Push(open, root)

	for open.Len() > 0 {
		// TODO: make concurrent
		node := heap.Pop(open).(*CTNode)
		conflictNode, agentA, agentB := node.findFirstConflict()
		if conflictNode == nil {
			return node.solution
		}
		newNodeA := node.fork(conflictNode, agentA)
		newNodeB := node.fork(conflictNode, agentB)
		heap.Push(open, newNodeA)
		heap.Push(open, newNodeB)
	}

	return nil
}
