package astar

import (
	"fmt"
	"math"
	"slices"
)

func Test1() {

	A := Node{Parent: nil, Position: Position{Col: 0, Row: 0}}
	B := Node{Parent: nil, Position: Position{Col: 0, Row: 14}}
	C := Node{Parent: nil, Position: Position{Col: 14, Row: 0}}
	D := Node{Parent: nil, Position: Position{Col: 14, Row: 14}}

	fmt.Println("node A:", A)
	fmt.Println("node B:", B)
	fmt.Println("node C:", C)
	fmt.Println("node D:", D)

	fmt.Println("A == C", A.eq(C))

	mazeB := CreateMaze()
	mazeC := CreateMaze()
	mazeD := CreateMaze()

	pathB := AStar(mazeB, A, B)
	fmt.Println("\n\nA* Path to B: ", pathB)
	PrintPathOnMaze(&mazeB, A, pathB)

	pathC := AStar(mazeC, A, C)
	fmt.Println("\n\nA* Path to C: ", pathC)
	PrintPathOnMaze(&mazeC, A, pathC)

	pathD := AStar(mazeD, A, D)
	fmt.Println("\n\nA* Path to D: ", pathD)
	PrintPathOnMaze(&mazeD, A, pathD)

}

func CreateMaze() Maze {
	return Maze{
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func PrintPathOnMaze(maze *Maze, start Node, path []Position) {

	newMaze := *maze
	newMaze[start.Position.Row][start.Position.Col] = 4
	fmt.Println("\nA* Path:")
	for _, p := range path {
		newMaze[p.Row][p.Col] = 4
	}
	for _, row := range newMaze {
		fmt.Println(row)
	}
}

type Position struct {
	Col int32
	Row int32
}

func (p Position) String() string {
	return fmt.Sprintf("(col=%v, row=%v)", p.Col, p.Row)
}

type Node struct {
	Parent   *Node
	Position Position
	f, g, h  float64
}

func (n Node) String() string {
	// return fmt.Sprintf(`Node{(%v, %v), f: %v, g: %v, h: %v} `,
	// 	n.Position.col, n.Position.row,
	// 	n.f, n.g, n.h,
	// )
	return fmt.Sprintf(`Node(col=%v, row=%v) `, n.Position.Col, n.Position.Row)
}

func (n *Node) eq(otherNode Node) bool {
	return n.Position.Row == otherNode.Position.Row &&
		n.Position.Col == otherNode.Position.Col
}

type Maze = [][]int

func AStar(maze Maze, start_node Node, end_node Node) []Position {

	start_node.f = 0
	start_node.g = 0
	start_node.h = 0
	start_node.Parent = nil

	end_node.f = 0
	end_node.g = 0
	end_node.h = 0

	open_list := []Node{}
	closed_list := []Node{}
	path := []Position{}

	open_list = append(open_list, start_node)

	adjacentMoves := []Position{
		{0, 1},   // north
		{1, 1},   // north-east
		{1, 0},   // east
		{1, -1},  // south-east
		{0, -1},  // south
		{-1, -1}, // south-west
		{-1, 0},  // west
		{-1, 1},  // north-west
	}
	// adjacentMoves := []Position{
	// 	{0, 1},   // north
	// 	{1, 0},   // east
	// 	{0, -1},  // south
	// 	{-1, 0},  // west
	// }

	// while loop
	isFinished := func() bool { return len(open_list) > 0 }

	for ok := isFinished(); ok; ok = isFinished() {

		current_node := open_list[0]
		current_index := 0

		// Get the current node
		for index, item := range open_list {
			if item.f < current_node.f {
				current_node = item
				current_index = index
			}
		}

		// Pop current off open list, add to closed list
		open_list, _ = SlicePop(open_list, current_index)
		closed_list = append(closed_list, current_node)

		// Generate Children
		children := []Node{}

		for _, move := range adjacentMoves {
			// Get node position
			node_position := Position{
				Row: current_node.Position.Row + move.Row,
				Col: current_node.Position.Col + move.Col,
			}
			// Make sure within range
			if !isInMap(maze, node_position) {
				continue
			}
			// Make sure walkable terrain
			if maze[node_position.Row][node_position.Col] != 0 {
				continue
			}

			// create new node
			new_node := Node{
				Parent:   &current_node,
				Position: node_position,
			}
			// append
			children = append(children, new_node)
		}

		for _, child := range children {
			// child is on the closed list
			if isNodeOnClosedList(closed_list, child) {
				continue
			}
			// Create the f, g, and h values
			child.g = current_node.g + 1
			child.h = nodeDistance(child, end_node)
			child.f = child.g + child.h
			// Child is already in the open list
			if isNodeOnOpenList(closed_list, child) {
				continue
			}
			// Add the child to the open list
			open_list = append(open_list, child)

		}

		// Found the goal
		if current_node.eq(end_node) {
			for current := current_node; !current.eq(start_node); current = *current.Parent {
				path = append(path, current.Position)
				if current.Parent == nil {
					break
				}
			}
			slices.Reverse(path)
			break
		}

	}

	return path
}

func isNodeOnClosedList(closed_list []Node, current_node Node) bool {
	for _, closed_child := range closed_list {
		if closed_child.eq(current_node) {
			return true
		}
	}
	return false
}
func isNodeOnOpenList(open_list []Node, child_node Node) bool {
	for _, open_node := range open_list {
		if child_node == open_node && child_node.g > open_node.g {
			return true
		}
	}
	return false
}

func nodeDistance(s Node, e Node) float64 {
	x1 := float64(s.Position.Row)
	x2 := float64(e.Position.Row)
	y1 := float64(s.Position.Col)
	y2 := float64(e.Position.Col)
	return math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2)
}

func isInMap(maze Maze, position Position) bool {

	row := position.Row
	col := position.Col

	maze_max_row := int32(len(maze) - 1)
	maze_max_col := int32(len(maze[len(maze)-1]) - 1)

	if row < 0 || row > maze_max_row {
		return false
	}
	if col < 0 || col > maze_max_col {
		return false
	}
	return true
}

func SlicePop[T any](s []T, i int) ([]T, T) {
	elem := s[i]
	s = append(s[:i], s[i+1:]...)
	return s, elem
}
