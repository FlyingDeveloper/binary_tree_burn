package main

import (
	"container/list"
	"fmt"
	"log"
)

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// buildGraph builds an undirected graph based on a root node of a binary tree
// by creating a map with all of the tree's nodes as keys and a slice of each
// node's neighbors as values
func buildGraph(root *Node, graph map[*Node][]*Node) {
	if root == nil {
		return
	}

	if root.Left != nil {
		graph[root.Left] = append(graph[root.Left], root)
		graph[root] = append(graph[root], root.Left)
		buildGraph(root.Left, graph)
	}
	if root.Right != nil {
		graph[root.Right] = append(graph[root.Right], root)
		graph[root] = append(graph[root], root.Right)
		buildGraph(root.Right, graph)
	}
}

// burnGraph "burns" through an undirected graph, with each burn stage containing
// the neighboring nodes of the nodes burning in the previous stage. The return
// value is a slice of stages of the burn, with each stage containing a slice of
// Nodes burning in that stage
func burnGraph(startNode *Node, graph map[*Node][]*Node) [][]*Node {
	if graph[startNode] == nil {
		log.Fatal("Unable to find start node")
	}

	visited := map[*Node]interface{}{}
	aflame := list.New()
	aflame.PushBack(startNode)

	burnStages := [][]*Node{}
	for aflame.Len() > 0 {
		c := aflame.Front()
		burnStage := []*Node{}
		for c != nil {
			burnStage = append(burnStage, c.Value.(*Node))
			c = c.Next()
		}
		burnStages = append(burnStages, burnStage)

		numAflame := aflame.Len()
		for i := 0; i < numAflame; i++ {
			current := aflame.Front().Value.(*Node)
			aflame.Remove(aflame.Front())
			visited[current] = struct{}{}
			currentNeighbors := graph[current]
			for _, n := range currentNeighbors {
				if _, alreadyVisited := visited[n]; alreadyVisited {
					continue
				}
				aflame.PushBack(n)
			}
		}
	}

	return burnStages
}

func example1() [][]*Node {
	n21 := Node{Value: 21}
	n24 := Node{Value: 24}
	n22 := Node{Value: 22}
	n23 := Node{Value: 23}
	n14 := Node{Value: 14, Left: &n21, Right: &n24}
	n15 := Node{Value: 15, Left: &n22, Right: &n23}
	n10 := Node{Value: 10, Left: &n14, Right: &n15}
	n13 := Node{Value: 13}
	n12 := Node{Value: 12, Left: &n13, Right: &n10}

	graph := map[*Node][]*Node{}
	buildGraph(&n12, graph)

	return burnGraph(&n14, graph)
}

func example2() [][]*Node {
	n2 := Node{Value: 2}
	n21 := Node{Value: 21}
	n7 := Node{Value: 7}
	n16 := Node{Value: 16}
	n41 := Node{Value: 41, Right: &n2}
	n15 := Node{Value: 15, Left: &n21}
	n95 := Node{Value: 95, Left: &n7, Right: &n16}
	n19 := Node{Value: 19, Left: &n41}
	n82 := Node{Value: 82, Left: &n15, Right: &n95}
	n12 := Node{Value: 12, Left: &n19, Right: &n82}

	graph := map[*Node][]*Node{}
	buildGraph(&n12, graph)
	return burnGraph(&n41, graph)
}

func printBurnStages(burnStages [][]*Node) {
	for _, stage := range burnStages {
		values := mapSlice(stage, func(item *Node) int { return item.Value })
		fmt.Printf("%v\n", values[0:])
	}
}

func mapSlice[T any, U any](input []T, f func(T) U) []U {
	output := []U{}
	for _, item := range input {
		output = append(output, f(item))
	}
	return output
}

func main() {
	burnStages1 := example1()
	printBurnStages(burnStages1)

	fmt.Println("=====")

	printBurnStages(example2())
}
