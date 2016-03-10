package jpsplus

const (
	Working = iota
	PathFound
	NoPathExists
)

const (
	SQRT_2           = 3363
	SQRT_2_MINUS_ONE = 985
)

type FunctionPointer func(*PathfindingNode, *JumpDistancesAndGoalBounds)

type xyLocJPS struct {
	x int
	y int
}

type JPSPlus struct {
	m_width                       int
	m_height                      int
	m_fastStack                   *FastStack
	m_simpleUnsortedPriorityQueue *SimpleUnsortedPriorityQueue
	m_jumpDistancesAndGoalBounds  []*JumpDistancesAndGoalBounds
	m_mapNodes                    [][]PathfindingNode
	m_currentIteration            int
	m_goalNode                    *PathfindingNode
	m_goalRow                     int
	m_goalCol                     int
}

func newJPSPlus(jumpDistancesAndGoalBoundsMap []*JumpDistancesAndGoalBounds, rawMap []bool, w int, h int) *JPSPlus {
	j := new(JPSPlus)
	j.m_width = w
	j.m_height = h
	j.m_simpleUnsortedPriorityQueue = newSimpleUnsortedPriorityQueue(10000)
	j.m_fastStack = newFastStack(1000)
	j.m_jumpDistancesAndGoalBounds = jumpDistancesAndGoalBoundsMap
	j.m_currentIteration = 1
	j.m_mapNodes = make([]([]PathfindingNode), h)
	for pos := 0; pos < h; pos++ {
		j.m_mapNodes[pos] = make([]PathfindingNode, w)
	}
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			node := &j.m_mapNodes[r][c]
			node.m_row = r
			node.m_col = c
			node.m_listStatus = PathfindingNode_OnNone
			node.m_iteration = 0
		}
	}
	return j
}

func (j *JPSPlus) GetPath(s xyLocJPS, g xyLocJPS) ([]xyLocJPS, bool) {

	startRow := s.y
	startCol := s.x
	j.m_goalRow = g.y
	j.m_goalCol = g.x

	{
		// Initialize map

		j.m_goalNode = &j.m_mapNodes[j.m_goalRow][j.m_goalCol]
		j.m_currentIteration += 1

		j.m_fastStack.Reset()
		j.m_simpleUnsortedPriorityQueue.Reset()
	}

	// Create starting node
	startNode := &j.m_mapNodes[startRow][startCol]
	startNode.m_parent = nil
	startNode.m_givenCost = 0
	startNode.m_finalCost = 0
	startNode.m_listStatus = PathfindingNode_OnOpen
	startNode.m_iteration = j.m_currentIteration

	// Actual search
	status := j.SearchLoop(startNode)

	if status == PathFound {
		path := j.FinalizePath([]xyLocJPS{}) //路径后续处理
		return path, true
	} else {
		// No path
		return []xyLocJPS{}, false
	}
}

func (j *JPSPlus) SearchLoop(startNode *PathfindingNode) int {
	{
		// Special case for the starting node

		if startNode == j.m_goalNode {
			return PathFound
		}

		jumpDistancesAndGoalBounds := &j.m_jumpDistancesAndGoalBounds[startNode.m_row][startNode.m_col]
		j.Explore_AllDirections(startNode, jumpDistancesAndGoalBounds)
		startNode.m_listStatus = PathfindingNode_OnClosed
	}

	for !j.m_simpleUnsortedPriorityQueue.Empty() || !j.m_fastStack.Empty() {
		var currentNode *PathfindingNode

		if !j.m_fastStack.Empty() {
			currentNode = j.m_fastStack.Pop()
		} else {
			currentNode = j.m_simpleUnsortedPriorityQueue.Pop()
		}

		if currentNode == j.m_goalNode {
			return PathFound
		}

		// Explore nodes based on parent
		jumpDistancesAndGoalBounds := &j.m_jumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]

		exploreDirections[(jumpDistancesAndGoalBounds.blockedDirectionBitfield*8)+
			currentNode.m_directionFromParent](currentNode, jumpDistancesAndGoalBounds)

		currentNode.m_listStatus = PathfindingNode_OnClosed
	}
	return NoPathExists
}

func (j *JPSPlus) FinalizePath(finalPath []xyLocJPS) {
	var prevNode *PathfindingNode
	curNode := j.m_goalNode

	while(curNode != nil)
	{
		loc := xyLocJPS{curNode.m_col, curNode.m_row}

		if prevNode != nil {
			// Insert extra nodes if needed (may not be neccessary depending on final path use)
			xDiff := curNode.m_col - prevNode.m_col
			yDiff := curNode.m_row - prevNode.m_row

			xInc := 0
			yInc := 0

			if xDiff > 0 {
				xInc = 1
			} else if xDiff < 0 {
				xInc = -1
				xDiff = -xDiff
			}

			if yDiff > 0 {
				yInc = 1
			} else if yDiff < 0 {
				yInc = -1
				yDiff = -yDiff
			}

			x := prevNode.m_col
			y := prevNode.m_row
			steps := xDiff - 1
			if yDiff > xDiff {
				steps = yDiff - 1
			}

			for i := 0; i < steps; i++ {
				x += xInc
				y += yInc

				locNew := xyLocJPS{x, y}
				finalPath.push_back(locNew)
			}
		}

		finalPath.push_back(loc)
		prevNode = curNode
		curNode = curNode.m_parent
	}
	reverse(finalPath.begin(), finalPath.end())
}
