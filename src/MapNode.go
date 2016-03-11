package jpsplus

type MapNode [][]*DijkstraPathfindingNode

func newMapNode(w int, h int) *MapNode {
	m := make(MapNode, h)
	for pos := 0; pos < h; pos++ {
		m[pos] = make([]*DijkstraPathfindingNode, w)
	}
	return &m
}

func (m MapNode) insert(r int, c int, node *DijkstraPathfindingNode) {
	if r >= 0 && r < m.height() {
		if c >= 0 && c < m.width() {
			m[r][c] = node
		}
	}
}

func (m MapNode) node(r int, c int) *DijkstraPathfindingNode {
	if r < 0 || r >= m.height() {
		return nil
	} else {
		if c < 0 || c >= m.width() {
			return nil
		} else {
			return m[r][c]
		}
	}
}

func (m MapNode) width() int {
	return len(m[0])
}

func (m MapNode) height() int {
	return len(m)
}
