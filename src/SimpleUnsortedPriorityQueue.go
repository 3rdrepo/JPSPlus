package jpsplus

type SimpleUnsortedPriorityQueue struct {
	nextFreeNode int
	nodeArray    []*PathfindingNode
}

func newSimpleUnsortedPriorityQueue(arraySize int) *SimpleUnsortedPriorityQueue {
	s := new(SimpleUnsortedPriorityQueue)
	s.nodeArray = make([]*PathfindingNode, arraySize)
	return s
}

func (s *SimpleUnsortedPriorityQueue) Reset() {
	s.nextFreeNode = 0
}

func (s *SimpleUnsortedPriorityQueue) Empty() bool {
	return 0 == s.nextFreeNode
}
func (s *SimpleUnsortedPriorityQueue) Add(node *PathfindingNode) {
	s.nodeArray[s.nextFreeNode] = node
	s.nextFreeNode += 1
}

func (s *SimpleUnsortedPriorityQueue) Pop() *PathfindingNode {
	cheapestNodeCost := s.nodeArray[0].finalCost
	cheapestNodeIndex := 0

	for i := 1; i < s.nextFreeNode; i++ {
		if s.nodeArray[i].finalCost < cheapestNodeCost {
			cheapestNodeCost = s.nodeArray[i].finalCost
			cheapestNodeIndex = i
		}
	}

	cheapestNode := s.nodeArray[cheapestNodeIndex]

	s.nextFreeNode -= 1
	s.nodeArray[cheapestNodeIndex] = s.nodeArray[s.nextFreeNode]

	return cheapestNode
}
