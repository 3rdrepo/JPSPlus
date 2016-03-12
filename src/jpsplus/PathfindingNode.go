package jpsplus

const (
	PathfindingNode_OnNone = iota
	PathfindingNode_OnOpen
	PathfindingNode_OnClosed
)

type PathfindingNode struct {
	m_parent              *PathfindingNode
	m_row                 int
	m_col                 int
	m_givenCost           int64
	m_finalCost           int64
	m_iteration           int
	m_directionFromParent int
	m_listStatus          int
}

type DijkstraPathfindingNode struct {
	m_parent                   *DijkstraPathfindingNode
	m_row                      int
	m_col                      int
	m_givenCost                int64
	m_iteration                int
	m_directionFromStart       int
	m_directionFromParent      int
	m_blockedDirectionBitfield int // highest bit [DownLeft, Left, UpLeft, Up, UpRight, Right, DownRight, Down] lowest bit
	m_listStatus               int
}
