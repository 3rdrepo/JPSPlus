package jpsplus

const (
	PathfindingNode_OnNone = iota
	PathfindingNode_OnOpen
	PathfindingNode_OnClosed
)

type Node struct {
	parent              *Node
	row                 int
	col                 int
	pos                 int
	givenCost           int64
	finalCost           int64
	iteration           int
	directionFromParent int
	listStatus          int
	heap_index          int
}

func newNode(r int, c int) *Node {
	node := new(Node)
	node.row = r
	node.col = c
	node.pos = r*MapWidth + c
	node.listStatus = PathfindingNode_OnNone
	return node
}

func (this Node) equal(n *Node) bool {
	return this.row == n.row && this.col == n.col
}
