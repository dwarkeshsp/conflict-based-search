package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dwarkeshsp/astar"
)

var GRID_SIZE = 100

func main() {
	var gridSize = flag.Int("s", 100, "Grid size")
	var agentsFile = flag.String("a", "", "Name of agents file")
	var obstaclesFile = flag.String("o", "", "Name of obstacles file")
	flag.Parse()

	GRID_SIZE = *gridSize
	agents := parseAgents(*agentsFile)
	obstacles := parseObstacles(*obstaclesFile)

	solution := plan(agents, obstacles)

	printAgentNodeMap(&solution)
}

func printAgentNodeMap(nodeMap *map[Agent][]astar.Node) {
	if nodeMap != nil {
		// i := 0
		// for _, path := range *nodeMap {
		// 	fmt.Printf("Agent %d path\n", i)
		// 	for _, node := range path {
		// 		fmt.Println(node)
		// 	}
		// 	i++
		// }
	} else {
		fmt.Println("No solution found")
	}
}

func parseAgents(agentsFile string) []Agent {
	agents := []Agent{}

	scanner := createScanner(agentsFile)
	for scanner.Scan() {
		coordinates := parseLine(scanner.Text())
		start := astar.Node{X: coordinates[0], Y: coordinates[1]}
		end := astar.Node{X: coordinates[1], Y: coordinates[2]}
		agent := Agent{start, end}
		agents = append(agents, agent)
	}

	return agents
}

func parseObstacles(obstaclesFile string) []astar.Node {
	nodes := []astar.Node{}

	scanner := createScanner(obstaclesFile)
	for scanner.Scan() {
		coordinates := parseLine(scanner.Text())
		node := astar.Node{X: coordinates[0], Y: coordinates[1]}
		nodes = append(nodes, node)
	}

	return nodes
}

func createScanner(filename string) *bufio.Scanner {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	return bufio.NewScanner(file)
}

func parseLine(line string) []int {
	stringList := strings.Split(line, " ")
	coordinates := []int{}
	for _, stringNum := range stringList {
		num, _ := strconv.Atoi(stringNum)
		coordinates = append(coordinates, num)
	}
	return coordinates
}
