package jpsplus

type SimpleUnsortedPriorityQueue struct {
	m_nextFreeNode int
	m_identical    bool
	m_nodeArray    []*PathfindingNode
}

func newSimpleUnsortedPriorityQueue(arraySize int) *SimpleUnsortedPriorityQueue {
	s := new(SimpleUnsortedPriorityQueue)
	s.m_nodeArray = make([]*PathfindingNode, arraySize)
	return s
}

func (s *SimpleUnsortedPriorityQueue) Reset() {
	s.m_nextFreeNode = 0
}

func (s *SimpleUnsortedPriorityQueue) Empty() bool {
	return 0 == s.m_nextFreeNode
}
func (s *SimpleUnsortedPriorityQueue) Add(node *PathfindingNode) {
	s.m_nodeArray[s.m_nextFreeNode] = node
	s.m_nextFreeNode += 1
}

func (s *SimpleUnsortedPriorityQueue) Pop() *PathfindingNode {
	// Find the cheapest node
	cheapestNodeCost := s.m_nodeArray[0].m_finalCost
	cheapestNodeIndex := 0

	for i := 1; i < s.m_nextFreeNode; i++ {
		if s.m_nodeArray[i].m_finalCost < cheapestNodeCost {
			cheapestNodeCost = s.m_nodeArray[i].m_finalCost
			cheapestNodeIndex = i
		}
	}

	// Remember the cheapest node
	cheapestNode := s.m_nodeArray[cheapestNodeIndex]

	// Delete off list (put the last node where this one was)
	s.m_nextFreeNode -= 1
	s.m_nodeArray[cheapestNodeIndex] = s.m_nodeArray[s.m_nextFreeNode]

	return cheapestNode
}
