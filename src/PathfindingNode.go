package jpsplus

const (
	PathfindingNode_OnNone = iota
	PathfindingNode_OnOpen
	PathfindingNode_OnClosed
)

type PathfindingNode struct {
	parent              *PathfindingNode
	row                 int
	col                 int
	givenCost           int64
	finalCost           int64
	iteration           int
	directionFromParent int
	listStatus          int
}

func newPathfindingNode(r int, c int) *PathfindingNode {
	node := new(PathfindingNode)
	node.row = r
	node.col = c
	node.listStatus = PathfindingNode_OnNone
	node.iteration = 0
	return node
}
