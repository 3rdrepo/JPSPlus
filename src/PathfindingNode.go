package jpsplus

const (
	PathfindingNodeStatus_OnNone = iota
	PathfindingNodeStatus_OnOpen
	PathfindingNodeStatus_OnClosed
)

type PathfindingNode struct {
	m_parent              *PathfindingNode
	m_row                 int
	m_col                 int
	m_givenCost           uint
	m_finalCost           uint
	m_iteration           uint
	m_directionFromParent uint8
	m_listStatus          uint8
}

type DijkstraPathfindingNode struct {
	m_parent                   *DijkstraPathfindingNode
	m_row                      int
	m_col                      int
	m_givenCost                uint
	m_iteration                uint
	m_directionFromStart       uint8
	m_directionFromParent      uint8
	m_blockedDirectionBitfield uint8 // highest bit [DownLeft, Left, UpLeft, Up, UpRight, Right, DownRight, Down] lowest bit
	m_listStatus               uint8
}
