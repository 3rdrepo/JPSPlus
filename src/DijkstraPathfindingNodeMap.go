package jpsplus

type DijkstraPathfindingNodeMap [][]*DijkstraPathfindingNode

func newMapNode(w int, h int) *DijkstraPathfindingNodeMap {
	d := make(DijkstraPathfindingNodeMap, h)
	for pos := 0; pos < h; pos++ {
		d[pos] = make([]*DijkstraPathfindingNode, w)
	}
	return &d
}

func (d DijkstraPathfindingNodeMap) insert(r int, c int, node *DijkstraPathfindingNode) {
	if r >= 0 && r < d.height() {
		if c >= 0 && c < d.width() {
			d[r][c] = node
		}
	}
}

func (d DijkstraPathfindingNodeMap) node(r int, c int) *DijkstraPathfindingNode {
	if r < 0 || r >= d.height() {
		return nil
	} else {
		if c < 0 || c >= d.width() {
			return nil
		} else {
			return d[r][c]
		}
	}
}

func (d DijkstraPathfindingNodeMap) width() int {
	return len(d[0])
}

func (d DijkstraPathfindingNodeMap) height() int {
	return len(d)
}

func (d DijkstraPathfindingNodeMap) nodeBlocked(r int, c int) int {
	return d[r][c].m_blockedDirectionBitfield
}
