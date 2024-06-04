package astar

import (
	"fmt"
	"math"
	"slices"
)

func Test1() {
	fmt.Println("astar!!!!!")

	A := Node{Parent: nil, Position: Position{col: 0, row: 0}}
	B := Node{Parent: nil, Position: Position{col: 0, row: 14}}
	C := Node{Parent: nil, Position: Position{col: 14, row: 0}}
	D := Node{Parent: nil, Position: Position{col: 14, row: 14}}

	fmt.Println("node A:", A)
	fmt.Println("node B:", B)
	fmt.Println("node C:", C)
	fmt.Println("node D:", D)

	fmt.Println("A == C", A.eq(C))

	mazeB := createMaze()
	mazeC := createMaze()
	mazeD := createMaze()

	pathB := AStar(mazeB, A, B)
	fmt.Println("\n\nA* Path to B: ", pathB)
	printPathOnMaze(&mazeB, A, pathB)

	pathC := AStar(mazeC, A, C)
	fmt.Println("\n\nA* Path to C: ", pathC)
	printPathOnMaze(&mazeC, A, pathC)

	pathD := AStar(mazeD, A, D)
	fmt.Println("\n\nA* Path to D: ", pathD)
	printPathOnMaze(&mazeD, A, pathD)

}

func createMaze() Maze {
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

func printPathOnMaze(maze *Maze, start Node, path []Position) {

	newMaze := *maze
	newMaze[start.Position.row][start.Position.col] = 4
	fmt.Println("\nA* Path:")
	for _, p := range path {
		newMaze[p.row][p.col] = 4
	}
	for _, row := range newMaze {
		fmt.Println(row)
	}
}

type Position struct {
	col int32
	row int32
}

func (p Position) String() string {
	return fmt.Sprintf("(col=%v, row=%v)", p.col, p.row)
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
	return fmt.Sprintf(`Node(col=%v, row=%v) `, n.Position.col, n.Position.row)
}

func (n *Node) eq(otherNode Node) bool {
	return n.Position.row == otherNode.Position.row &&
		n.Position.col == otherNode.Position.col
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

	var adjacentMoves []Position
	adjacentMoves = []Position{
		Position{0, 1},   // north
		Position{1, 1},   // north-east
		Position{1, 0},   // east
		Position{1, -1},  // south-east
		Position{0, -1},  // south
		Position{-1, -1}, // south-west
		Position{-1, 0},  // west
		Position{-1, 1},  // north-west
	}
	// adjacentMoves = []Position{
	// 	Position{0, 1},   // north
	// 	Position{1, 0},   // east
	// 	Position{0, -1},  // south
	// 	Position{-1, 0},  // west
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
				row: current_node.Position.row + move.row,
				col: current_node.Position.col + move.col,
			}
			// Make sure within range
			if !isInMap(maze, node_position) {
				continue
			}
			// Make sure walkable terrain
			if maze[node_position.row][node_position.col] != 0 {
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
	x1 := float64(s.Position.row)
	x2 := float64(e.Position.row)
	y1 := float64(s.Position.col)
	y2 := float64(e.Position.col)
	return math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2)
}

func isInMap(maze Maze, position Position) bool {

	row := position.row
	col := position.col

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
