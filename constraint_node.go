package main

import (
	"time"

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

func createRootNode(agents []Agent, obstacles []astar.Node) *CTNode {
	n := &CTNode{}
	n.constaints = make(map[Agent][]astar.Node)
	for _, agent := range agents {
		n.constaints[agent] = obstacles
	}
	n.solution = make(map[Agent][]astar.Node)
	if !n.findSolution() {
		return nil
	}
	return n
}

func (n *CTNode) findSolution() bool {
	for agent, obstacles := range n.constaints {
		aConfig := astar.Config{
			GridWidth:    GRID_SIZE,
			GridHeight:   GRID_SIZE,
			InvalidNodes: obstacles,
		}

		algo, _ := astar.New(aConfig)
		// TODO: deal with no solution found
		agentPath, _ := algo.FindPath(agent.start, agent.end)
		if agentPath == nil {
			return false
		}
		n.solution[agent] = agentPath
	}
	n.storeCost()
	return true
}

func (n *CTNode) storeCost() {
	cost := 0
	for _, points := range n.solution {
		cost += len(points)
	}
	n.cost = cost
}

type ConflictResult struct {
	n      *astar.Node
	aIndex int
	bIndex int
}

func (n *CTNode) findFirstConflict() (*astar.Node, *Agent, *Agent) {
	conflictsChan := make(chan *ConflictResult)
	finishedChan := make(chan bool)

	agents := make([]Agent, len(n.solution))

	i := 0
	for agent := range n.solution {
		agents[i] = agent
		i++
	}

	workers := 0
	for i := 0; i < len(agents)-1; i++ {
		for j := i + 1; j < len(agents); j++ {
			workers++
			go findPathConflict(&conflictsChan, &finishedChan, n.solution[agents[i]], n.solution[agents[j]], i, j)
		}
	}

	select {
	case result := <-conflictsChan:
		println("hererrre")
		return result.n, &agents[result.aIndex], &agents[result.bIndex]
	case <-time.After(3 * time.Second):
		return nil, nil, nil
	}
}

func findPathConflict(conflictsChan *chan *ConflictResult, finishedChan *chan bool, a []astar.Node, b []astar.Node, aIndex int, bIndex int) {

	size := len(a)
	if size > len(b) {
		size = len(b)
	}

	for i := 0; i < size; i++ {
		// select {
		// case <-*finishedChan:
		// 	return
		// default:
		// }
		if a[i].X == b[i].X && a[i].Y == b[i].Y {
			println("CONFLICT FOUND")
			println(a[i].X, a[i].Y, b[i].X, b[i].Y)
			select {
			case *conflictsChan <- &ConflictResult{&a[i], aIndex, bIndex}:
			default:
			}
			break
		}
	}

	// <-*finishedChan
}

// func sendFinished(workers int, finishedChan *chan bool) {
// 	for i := 0; i < workers; i++ {
// 		*finishedChan <- true
// 	}
// }

func (n *CTNode) fork(conflictNode *astar.Node, restrictedAgent *Agent) *CTNode {
	newNode := &CTNode{}
	newNode.constaints = make(map[Agent][]astar.Node)
	for agent, restrictions := range n.constaints {
		newRestrictions := make([]astar.Node, len(restrictions))
		copy(newRestrictions, restrictions)
		newNode.constaints[agent] = newRestrictions
	}
	newNode.constaints[*restrictedAgent] = append(newNode.constaints[*restrictedAgent], *conflictNode)

	newNode.solution = make(map[Agent][]astar.Node)
	if !newNode.findSolution() {
		return nil
	}
	return newNode
}
