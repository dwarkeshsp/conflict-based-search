package main

import (
	"container/heap"
	"time"

	"github.com/dwarkeshsp/astar"
)

func plan(agents []Agent, obstacles []astar.Node) (map[Agent][]astar.Node, int) {
	open := make(CTree, 1)
	root := createRootNode(agents, obstacles)
	if root == nil {
		return nil, 0
	}
	open[0] = root
	heap.Init(&open)

	lastPushed := time.Now()
	resultChan := make(chan *CTNode)

	for open.Len() > 0 || time.Since(lastPushed) < 2*time.Second {
		go considerBestNode(&open, resultChan, &lastPushed)
		select {
		case result := <-resultChan:
			return result.solution, result.cost
		default:
		}
	}

	select {
	case result := <-resultChan:
		return result.solution, result.cost
	case <-time.After(20 * time.Second):
		return nil, 0
	}

}

func considerBestNode(open *CTree, solutionChan chan *CTNode, lastPushed *time.Time) {

	var node *CTNode
	mu.Lock()
	if len(*open) > 0 {
		node = heap.Pop(open).(*CTNode)
	} else {
		return
	}
	mu.Unlock()

	conflictNode, agentA, agentB := node.findFirstConflict()
	if conflictNode == nil {
		solutionChan <- node
		return
	}

	newNodeA := node.fork(conflictNode, agentA)
	newNodeB := node.fork(conflictNode, agentB)

	if newNodeA != nil {
		mu.Lock()
		heap.Push(open, newNodeA)
		mu.Unlock()
	}
	if newNodeB != nil {
		mu.Lock()
		heap.Push(open, newNodeB)
		mu.Unlock()
	}
	*lastPushed = time.Now()
}
