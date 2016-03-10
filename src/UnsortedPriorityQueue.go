package jpsplus

type UnsortedPriorityQueue struct {
	m_nextFreeNode int
	m_iteration    int
	m_identical    bool
	m_nodeArray    []*DijkstraPathfindingNode
}

func newUnsortedPriorityQueue(arraySize int) *UnsortedPriorityQueue {
	u := new(UnsortedPriorityQueue)
	u.m_nodeArray = make([]*DijkstraPathfindingNode, arraySize)
}

func (u UnsortedPriorityQueue) Empty(int iteration) bool {
	if u.m_iteration == u.iteration {
		return 0 == u.m_nextFreeNode
	} else {
		return true
	}
}

func (u UnsortedPriorityQueue) GetIteration() int {
	return u.m_iteration
}

func (u *UnsortedPriorityQueue) Push(node *DijkstraPathfindingNode) {
	if u.m_iteration != node.m_iteration {
		u.m_nextFreeNode = 0
		u.m_identical = true
		u.m_iteration = node.m_iteration
	}
	u.m_nextFreeNode += 1
	u.m_nodeArray[u.m_nextFreeNode] = node
}

func (u *UnsortedPriorityQueue) Remove(node *DijkstraPathfindingNode) {
	for i := 0; i < u.m_nextFreeNode; i++ {
		if u.m_nodeArray[i] == node {
			// Delete off Open list (put last node where this one was)
			u.m_nextFreeNode -= 1
			u.m_nodeArray[i] = u.m_nodeArray[u.m_nextFreeNode]
			return
		}
	}
}

func (u *UnsortedPriorityQueue) Pop() *DijkstraPathfindingNode {
	// Find cheapest node
	cheapestNodeCostFinal := u.m_nodeArray[0].m_givenCost
	cheapestNodeIndex := 0

	for i := 1; i < u.m_nextFreeNode; i++ {
		if u.m_nodeArray[i].m_givenCost < cheapestNodeCostFinal {
			cheapestNodeCostFinal = u.m_nodeArray[i].m_givenCost
			cheapestNodeIndex = i
		}
	}

	// Remember cheapest node
	cheapestNode := u.m_nodeArray[cheapestNodeIndex]

	// Delete off Open list (put last node where this one was)
	u.m_nextFreeNode -= 1
	u.m_nodeArray[cheapestNodeIndex] = u.m_nodeArray[u.m_nextFreeNode]

	return cheapestNode
}
