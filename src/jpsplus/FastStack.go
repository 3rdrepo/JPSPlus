package jpsplus

type FastStack struct {
	m_nextFreeNode int
	m_nodeArray    []*PathfindingNode
}

func newFastStack(size int) *FastStack {
	f := new(FastStack)
	f.m_nodeArray = make([]*PathfindingNode, size)
	return f
}

func (f *FastStack) Reset() {
	f.m_nextFreeNode = 0
}

func (f FastStack) Empty() bool {
	return 0 == f.m_nextFreeNode
}

func (f *FastStack) Push(node *PathfindingNode) {
	f.m_nextFreeNode += 1
	f.m_nodeArray[f.m_nextFreeNode] = node
}

func (f *FastStack) Pop() *PathfindingNode {
	t := f.m_nextFreeNode
	f.m_nextFreeNode -= 1
	return f.m_nodeArray[t]
}
