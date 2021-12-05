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

	agents := make([]Agent, len(n.solution))

	i := 0
	for agent := range n.solution {
		agents[i] = agent
		i++
	}
	println(8)

	var wg sync.WaitGroup

	for i := 0; i < len(agents)-1; i++ {
		for j := i + 1; j < len(agents); j++ {
			wg.Add(1)
			go findPathConflict(&wg, conflictsChan, n.solution[agents[i]], n.solution[agents[j]], i, j)
		}
	}

	wg.Wait()
	select {
	case result := <-conflictsChan:
		return result.n, &agents[result.aIndex], &agents[result.bIndex]
	default:
		return nil, nil, nil

	}
}

func findPathConflict(wg *sync.WaitGroup, conflictsChan chan *ConflictResult, a []astar.Node, b []astar.Node, aIndex int, bIndex int) {
	defer wg.Done()

	size := len(a)
	if size > len(b) {
		size = len(b)
	}

	for i := 0; i < size; i++ {
		if a[i].X == b[i].X && a[i].Y == b[i].Y {
			select {
			case conflictsChan <- &ConflictResult{&a[i], aIndex, bIndex}:
			default:
			}
			return
		}
	}
}

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
