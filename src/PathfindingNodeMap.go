package jpsplus

type PathfindingNodeMap [][]*PathfindingNode

func newPathfindingNodeMap(width int, height int) *PathfindingNodeMap {
	p := make(PathfindingNodeMap, height)
	for r := 0; r < height; r++ {
		p[r] = make([]*PathfindingNode, width)
		for c := 0; c < width; c++ {
			node := new(PathfindingNode)
			node.m_row = r
			node.m_col = c
			node.m_listStatus = PathfindingNode_OnNone
			node.m_iteration = 0
			p[r][c] = node
		}
	}
	return &p
}
