package jpsplus

type FastStack struct {
	nextFreeNode int
	nodeArray    []*PathfindingNode
}

func newFastStack(size int) *FastStack {
	f := new(FastStack)
	f.nodeArray = make([]*PathfindingNode, size)
	return f
}

func (f *FastStack) Reset() {
	f.nextFreeNode = 0
}

func (f FastStack) Empty() bool {
	return 0 == f.nextFreeNode
}

func (f *FastStack) Push(node *PathfindingNode) {
	f.nextFreeNode += 1
	f.nodeArray[f.nextFreeNode] = node
}

func (f *FastStack) Pop() *PathfindingNode {
	t := f.nextFreeNode
	f.nextFreeNode -= 1
	return f.nodeArray[t]
}
