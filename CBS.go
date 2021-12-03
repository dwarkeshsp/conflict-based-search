package main

import "github.com/dwarkeshsp/astar"

func plan(agents []Agent) map[Agent][]astar.Node {
	open := &CTree{queue: make([]*CTNode, 100)}
	root := emptyNode(agents)
	root.findSolution()

	open.Push(&root)

	for open.Len() > 0 {

	}

	return nil

}
