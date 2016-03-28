package jpsplus

import (
	"fmt"
)

// func abs(a int) int {
// 	if a < 0 {
// 		return -a
// 	}
// 	return a
// }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

const (
	LAND = 1 << iota
	WALL
)

const (
	COST_STRAIGHT = 1000
	COST_DIAGONAL = 1414
)

type MapData [][]int

func NewMapData(rows, cols int) MapData {
	result := make([]([]int), rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]int, cols)
	}
	return result
}

func (m MapData) Clone() MapData {
	rows := len(m)
	cols := len(m[0])
	result := make([]([]int), rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]int, cols)
		copy(result[i], m[i])
	}
	return result
}

func str_map(data MapData, nodes []*Node) string {
	var result string
	for i, row := range data {
		for j, cell := range row {
			added := false
			for _, node := range nodes {
				if node.X == i && node.Y == j {
					result += "o"
					added = true
					break
				}
			}
			if !added {
				switch cell {
				case LAND:
					result += "."
				case WALL:
					result += "#"
				default:
					result += "?"
				}
			}
		}
		result += "\n"
	}
	return result
}

type Node struct {
	X, Y       int
	parent     *Node
	f, g, h    int
	heap_index int
}

func NewNode(x, y int) *Node {
	node := &Node{
		X:      x,
		Y:      y,
		parent: nil,
		f:      0,
		g:      0,
		h:      0,
	}
	return node
}

func (n *Node) String() string {
	return fmt.Sprintf("<Node x:%d y:%d addr:%d>", n.X, n.Y, &n)
}

type nodeList struct {
	nodes      map[int]*Node
	rows, cols int
}

func newNodeList(rows, cols int) *nodeList {
	return &nodeList{
		nodes: make(map[int]*Node, rows*cols),
		rows:  rows,
		cols:  cols,
	}
}

func (n *nodeList) addNode(node *Node) {
	n.nodes[node.X+node.Y*n.rows] = node
}

func (n *nodeList) getNode(x, y int) *Node {
	return n.nodes[x+y*n.rows]
}

func (n *nodeList) removeNode(node *Node) {
	delete(n.nodes, node.X+node.Y*n.rows)
}

func (n *nodeList) hasNode(node *Node) bool {
	if n.getNode(node.X, node.Y) != nil {
		return true
	}
	return false
}

type Graph struct {
	nodes *nodeList
	data  MapData
}

func NewGraph(map_data MapData) *Graph {
	return &Graph{
		nodes: newNodeList(len(map_data), len(map_data[0])),
		data:  map_data,
	}
}

func (g *Graph) Node(x, y int) *Node {
	var node *Node
	node = g.nodes.getNode(x, y)

	if node == nil && (g.data[x][y] != WALL) {
		node = NewNode(x, y)
		g.nodes.addNode(node)
	}
	return node
}

func retracePath(current_node *Node) []*Node {
	var path []*Node
	path = append(path, current_node)
	for current_node.parent != nil {
		path = append(path, current_node.parent)
		current_node = current_node.parent
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func Heuristic(tile, stop *Node) (h int) {
	h_diag := min(abs(tile.X-stop.X), abs(tile.Y-stop.Y))
	h_stra := abs(tile.X-stop.X) + abs(tile.Y-stop.Y)
	h = COST_DIAGONAL*h_diag + COST_STRAIGHT*(h_stra-2*h_diag)
	return
}

var adjecentDirs8 = [][3]int{
	{-1, -1, COST_DIAGONAL}, {-1, 0, COST_STRAIGHT}, {-1, 1, COST_DIAGONAL},
	{0, -1, COST_STRAIGHT}, {0, 1, COST_STRAIGHT},
	{1, -1, COST_DIAGONAL}, {1, 0, COST_STRAIGHT}, {1, 1, COST_DIAGONAL},
}

func Astar(map_data MapData, startx, starty, stopx, stopy int, dir8 bool) []*Node {
	graph := NewGraph(map_data)
	rows, cols := len(graph.data), len(graph.data[0])

	closedSet := newNodeList(rows, cols)
	openSet := newNodeList(rows, cols)
	pq := make(PriorityQueue, 0, rows*cols)
	start := NewNode(startx, starty)
	stop := NewNode(stopx, stopy)
	openSet.addNode(start)
	pq.PushNode(start)

	for len(openSet.nodes) != 0 {
		current := pq.PopNode()
		openSet.removeNode(current)
		closedSet.addNode(current)

		if current.X == stop.X && current.Y == stop.Y {
			return retracePath(current)
		}

		for _, adir := range adjecentDirs8 {
			x, y := (current.X + adir[0]), (current.Y + adir[1])

			if (x < 0) || (x >= rows) || (y < 0) || (y >= cols) {
				continue
			}

			neighbor := graph.Node(x, y)
			if neighbor == nil || closedSet.hasNode(neighbor) {
				continue
			}

			g_score := current.g + adir[2]

			if !openSet.hasNode(neighbor) {
				neighbor.parent = current
				neighbor.g = g_score
				neighbor.f = neighbor.g + Heuristic(neighbor, stop)
				openSet.addNode(neighbor)
				pq.PushNode(neighbor)
			} else if g_score < neighbor.g {
				pq.RemoveNode(neighbor)
				neighbor.parent = current
				neighbor.g = g_score
				neighbor.f = neighbor.g + Heuristic(neighbor, stop)
				pq.PushNode(neighbor)
			}

		}
	}

	return nil
}
