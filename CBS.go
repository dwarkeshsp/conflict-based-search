package main

import (
	"container/heap"

	"github.com/dwarkeshsp/astar"
)

func plan(agents []Agent, obstacles []astar.Node) map[Agent][]astar.Node {
	open := make(CTree, 1)
	root := createRootNode(agents, obstacles)
	if root == nil {
		return nil
	}
	open[0] = root

	heap.Init(&open)

	// heap.Push(&open, root)

	for open.Len() > 0 {
		// TODO: make concurrent
		println("5")
		println("len", open.Len())

		node := heap.Pop(&open).(*CTNode)
		println(6)

		conflictNode, agentA, agentB := node.findFirstConflict()
		println(7)
		if conflictNode == nil {
			println("null conflictnode")
			return node.solution
		}
		newNodeA := node.fork(conflictNode, agentA)
		newNodeB := node.fork(conflictNode, agentB)
		if newNodeA != nil {
			heap.Push(&open, newNodeA)
		}
		if newNodeB != nil {
			heap.Push(&open, newNodeB)
		}
	}

	return nil
}
