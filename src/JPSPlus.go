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

type xyLocJPS struct {
	x int
	y int
}

var DefaultJumpDistancesAndGoalBounds [][]JumpDistancesAndGoalBounds

type JPSPlus struct {
	m_width                       int
	m_height                      int
	m_fastStack                   *FastStack
	m_simpleUnsortedPriorityQueue *SimpleUnsortedPriorityQueue
	// m_jumpDistancesAndGoalBounds  [][]JumpDistancesAndGoalBounds
	m_mapNodes         [][]PathfindingNode
	m_currentIteration int
	m_goalNode         *PathfindingNode
	m_goalRow          int
	m_goalCol          int
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
	jpsPlus = j
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
		path := j.FinalizePath() //路径后续处理
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

		// jumpDistancesAndGoalBounds := &j.m_jumpDistancesAndGoalBounds[startNode.m_row][startNode.m_col]
		j.Explore_AllDirections(startNode, j)
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
		// jumpDistancesAndGoalBounds := &j.m_jumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]

		exploreDirections[(jumpDistancesAndGoalBounds.blockedDirectionBitfield*8)+
			currentNode.m_directionFromParent](currentNode, j)

		currentNode.m_listStatus = PathfindingNode_OnClosed
	}
	return NoPathExists
}

func (j *JPSPlus) FinalizePath() []xyLocJPS {
	var prevNode *PathfindingNode
	curNode := j.m_goalNode
	finalPath := make([]xyLocJPS, 0, 100)

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
				finalPath = append(finalPath, locNew)
			}
		}

		finalPath = append(finalPath, locNew)
		prevNode = curNode
		curNode = curNode.m_parent
	}

	return reverse(finalPath)
}

func reverse(path []xyLocJPS) []xyLocJPS {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func MacroExploreDown(currentNode *PathfindingNode, jpsPlus *JPSPlus) {
	jum := DefaultJumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]
	if jpsPlus.m_goalRow >= jmp.bounds[Down][MinRow] &&
		jpsPlus.m_goalRow <= jmp.bounds[Down][MaxRow] &&
		jpsPlus.m_goalCol >= jmp.bounds[Down][MinCol] &&
		jpsPlus.m_goalCol <= jmp.bounds[Down][MaxCol] {
		jpsPlus.SearchDown(currentNode, jum.jumpDistance[Down])
	}
}

func MacroExploreDownRight(currentNode *PathfindingNode, jpsPlus *JPSPlus) {
	jum := DefaultJumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]
	if m_goalRow >= jmp.bounds[DownRight][MinRow] &&
		m_goalRow <= jmp.bounds[DownRight][MaxRow] &&
		m_goalCol >= jmp.bounds[DownRight][MinCol] &&
		m_goalCol <= jmp.bounds[DownRight][MaxCol] {
		jpsPlus.SearchDownRight(currentNode, jmp.jumpDistance[DownRight])
	}

}

func MacroExploreRight(currentNode *PathfindingNode, jpsPlus *JPSPlus) {
	jum := DefaultJumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]
	if m_goalRow >= jmp.bounds[Right][MinRow] &&
		m_goalRow <= jmp.bounds[Right][MaxRow] &&
		m_goalCol >= jmp.bounds[Right][MinCol] &&
		m_goalCol <= jmp.bounds[Right][MaxCol] {
		jpsPlus.SearchRight(currentNode, jmp.jumpDistance[Right])
	}
}

func MacroExploreUpRight(currentNode *PathfindingNode, jpsPlus *JPSPlus) {
	jum := DefaultJumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]
	if m_goalRow >= jmp.bounds[UpRight][MinRow] &&
		m_goalRow <= jmp.bounds[UpRight][MaxRow] &&
		m_goalCol >= jmp.bounds[UpRight][MinCol] &&
		m_goalCol <= jmp.bounds[UpRight][MaxCol] {
		jpsPlus.SearchUpRight(currentNode, jmp.jumpDistance[UpRight])
	}
}

func MacroExploreUp(currentNode *PathfindingNode, jpsPlus *JPSPlus) {
	jum := DefaultJumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]
	if m_goalRow >= jmp.bounds[Up][MinRow] &&
		m_goalRow <= jmp.bounds[Up][MaxRow] &&
		m_goalCol >= jmp.bounds[Up][MinCol] &&
		m_goalCol <= jmp.bounds[Up][MaxCol] {
		jpsPlus.SearchUp(currentNode, jmp.jumpDistance[Up])
	}
}

func MacroExploreUpLeft(currentNode *PathfindingNode, jpsPlus *JPSPlus) {
	jum := DefaultJumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]
	if m_goalRow >= jmp.bounds[UpLeft][MinRow] &&
		m_goalRow <= jmp.bounds[UpLeft][MaxRow] &&
		m_goalCol >= jmp.bounds[UpLeft][MinCol] &&
		m_goalCol <= jmp.bounds[UpLeft][MaxCol] {
		jpsPlus.SearchUpLeft(currentNode, jmp.jumpDistance[UpLeft])
	}
}

func MacroExploreLeft(currentNode *PathfindingNode, jpsPlus *JPSPlus) {
	jum := DefaultJumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]
	if m_goalRow >= jmp.bounds[Left][MinRow] &&
		m_goalRow <= jmp.bounds[Left][MaxRow] &&
		m_goalCol >= jmp.bounds[Left][MinCol] &&
		m_goalCol <= jmp.bounds[Left][MaxCol] {
		jpsPlus.SearchLeft(currentNode, jmp.jumpDistance[Left])
	}
}

func MacroExploreDownLeft(currentNode *PathfindingNode, jpsPlus *JPSPlus) {
	jum := DefaultJumpDistancesAndGoalBounds[currentNode.m_row][currentNode.m_col]
	if m_goalRow >= jmp.bounds[DownLeft][MinRow] &&
		m_goalRow <= jmp.bounds[DownLeft][MaxRow] &&
		m_goalCol >= jmp.bounds[DownLeft][MinCol] &&
		m_goalCol <= jmp.bounds[DownLeft][MaxCol] {
		jpsPlus.SearchDownLeft(currentNode, jmp.jumpDistance[DownLeft])
	}
}

func JPSPlusExplore_Null(currentNode *PathfindingNode, jps *JPSPlus) {
	// Purposely does nothing
}

func JPSPlusExplore_D(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
}

func JPSPlusExplore_DR(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jps)
}

func JPSPlusExplore_R(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
}

func JPSPlusExplore_UR(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jps)
}

func JPSPlusExplore_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUp(currentNode, jps)
}

func JPSPlusExplore_UL(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jps)
}

func JPSPlusExplore_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jps)
}

func JPSPlusExplore_DL(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jps)
}

// Adjacent Doubles

func JPSPlusExplore_D_DR(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreDownRight(currentNode, jps)
}

func JPSPlusExplore_DR_R(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jps)
	MacroExploreRight(currentNode, jps)
}

func JPSPlusExplore_R_UR(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
	MacroExploreUpRight(currentNode, jps)
}

func JPSPlusExplore_UR_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jps)
	MacroExploreUp(currentNode, jps)
}

func JPSPlusExplore_U_UL(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUp(currentNode, jps)
	MacroExploreUpLeft(currentNode, jps)
}

func JPSPlusExplore_UL_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
}

func JPSPlusExplore_L_DL(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jps)
	MacroExploreDownLeft(currentNode, jps)
}

func JPSPlusExplore_DL_D(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jps)
	MacroExploreDown(currentNode, jps)
}

// Non-Adjacent Cardinal Doubles

func JPSPlusExplore_D_R(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreRight(currentNode, jps)
}

func JPSPlusExplore_R_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
	MacroExploreUp(currentNode, jps)
}

func JPSPlusExplore_U_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUp(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
}

func JPSPlusExplore_L_D(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jps)
	MacroExploreDown(currentNode, jps)
}

func JPSPlusExplore_D_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreUp(currentNode, jps)
}

func JPSPlusExplore_R_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
}

// Adjacent Triples

func JPSPlusExplore_D_DR_R(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreDownRight(currentNode, jps)
	MacroExploreRight(currentNode, jps)
}

func JPSPlusExplore_DR_R_UR(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jps)
	MacroExploreRight(currentNode, jps)
	MacroExploreUpRight(currentNode, jps)
}

func JPSPlusExplore_R_UR_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
	MacroExploreUpRight(currentNode, jps)
	MacroExploreUp(currentNode, jps)
}

func JPSPlusExplore_UR_U_UL(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jps)
	MacroExploreUp(currentNode, jps)
	MacroExploreUpLeft(currentNode, jps)
}

func JPSPlusExplore_U_UL_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUp(currentNode, jps)
	MacroExploreUpLeft(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
}

func JPSPlusExplore_UL_L_DL(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
	MacroExploreDownLeft(currentNode, jps)
}

func JPSPlusExplore_L_DL_D(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jps)
	MacroExploreDownLeft(currentNode, jps)
	MacroExploreDown(currentNode, jps)
}

func JPSPlusExplore_DL_D_DR(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jps)
	MacroExploreDown(currentNode, jps)
	MacroExploreDownRight(currentNode, jps)
}

// Non-Adjacent Cardinal Triples

func JPSPlusExplore_D_R_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreRight(currentNode, jps)
	MacroExploreUp(currentNode, jps)
}

func JPSPlusExplore_R_U_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
	MacroExploreUp(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
}

func JPSPlusExplore_U_L_D(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUp(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
	MacroExploreDown(currentNode, jps)
}

func JPSPlusExplore_L_D_R(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jps)
	MacroExploreDown(currentNode, jps)
	MacroExploreRight(currentNode, jps)
}

// Quads

func JPSPlusExplore_R_DR_D_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
	MacroExploreDownRight(currentNode, jps)
	MacroExploreDown(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
}

func JPSPlusExplore_R_D_DL_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
	MacroExploreDown(currentNode, jps)
	MacroExploreDownLeft(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
}

func JPSPlusExplore_U_UR_R_D(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUp(currentNode, jps)
	MacroExploreUpRight(currentNode, jps)
	MacroExploreRight(currentNode, jps)
	MacroExploreDown(currentNode, jps)
}

func JPSPlusExplore_U_R_DR_D(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUp(currentNode, jps)
	MacroExploreRight(currentNode, jps)
	MacroExploreDownRight(currentNode, jps)
	MacroExploreDown(currentNode, jps)
}

func JPSPlusExplore_L_UL_U_R(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jps)
	MacroExploreUpLeft(currentNode, jps)
	MacroExploreUp(currentNode, jps)
	MacroExploreRight(currentNode, jps)
}

func JPSPlusExplore_L_U_UR_R(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jps)
	MacroExploreUp(currentNode, jps)
	MacroExploreUpRight(currentNode, jps)
	MacroExploreRight(currentNode, jps)
}

func JPSPlusExplore_D_DL_L_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreDownLeft(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
	MacroExploreUp(currentNode, jps)
}

func JPSPlusExplore_D_L_UL_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
	MacroExploreUpLeft(currentNode, jps)
	MacroExploreUp(currentNode, jps)
}

// Quints

func JPSPlusExplore_R_DR_D_DL_L(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreRight(currentNode, jps)
	MacroExploreDownRight(currentNode, jps)
	MacroExploreDown(currentNode, jps)
	MacroExploreDownLeft(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
}

func JPSPlusExplore_U_UR_R_DR_D(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreUp(currentNode, jps)
	MacroExploreUpRight(currentNode, jps)
	MacroExploreRight(currentNode, jps)
	MacroExploreDownRight(currentNode, jps)
	MacroExploreDown(currentNode, jps)
}

func JPSPlusExplore_L_UL_U_UR_R(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jps)
	MacroExploreUpLeft(currentNode, jps)
	MacroExploreUp(currentNode, jps)
	MacroExploreUpRight(currentNode, jps)
	MacroExploreRight(currentNode, jps)
}

func JPSPlusExplore_D_DL_L_UL_U(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreDownLeft(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
	MacroExploreUpLeft(currentNode, jps)
	MacroExploreUp(currentNode, jps)
}

func JPSPlusExplore_AllDirections(currentNode *PathfindingNode, jps *JPSPlus) {
	MacroExploreDown(currentNode, jps)
	MacroExploreDownLeft(currentNode, jps)
	MacroExploreLeft(currentNode, jps)
	MacroExploreUpLeft(currentNode, jps)
	MacroExploreUp(currentNode, jps)
	MacroExploreUpRight(currentNode, jps)
	MacroExploreRight(currentNode, jps)
	MacroExploreDownRight(currentNode, jps)
}

func (j *JPSPlus) SearchDown(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Consider straight line to Goal
	if col == j.m_goalCol && row < j.m_goalRow {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (row + absJumpDistance) >= j.m_goalRow {
			diff := j.m_goalRow - row
			givenCost := currentNode.m_givenCost + fixed_point_shift(diff)
			newSuccessor := j.m_goalNode
			j.PushNewNode(newSuccessor, currentNode, Down, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row + jumpDistance
		givenCost := currentNode.m_givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := &j.m_mapNodes[newRow][col]
		j.PushNewNode(newSuccessor, currentNode, Down, givenCost)
	}
}

func (j *JPSPlus) SearchDownRight(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Check for goal in general direction (straight line to Goal or Target Jump Point)
	if row < j.m_goalRow && col < j.m_goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := j.m_goalRow - row
		diffCol := j.m_goalCol - col
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row + smallerDiff
			newCol := col + smallerDiff
			givenCost := currentNode.m_givenCost + int64(SQRT_2*smallerDiff)
			newSuccessor := &j.m_mapNodes[newRow][newCol]
			j.PushNewNode(newSuccessor, currentNode, DownRight, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row + jumpDistance
		newCol := col + jumpDistance
		givenCost := currentNode.m_givenCost + (SQRT_2 * jumpDistance)
		newSuccessor := &j.m_mapNodes[newRow][newCol]
		j.PushNewNode(newSuccessor, currentNode, DownRight, givenCost)
	}
}

func (j *JPSPlus) SearchRight(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Consider straight line to Goal
	if row == j.m_goalRow && col < j.m_goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (col + absJumpDistance) >= j.m_goalCol {
			diff := j.m_goalCol - col
			givenCost := currentNode.m_givenCost + fixed_point_shift(diff)
			newSuccessor := j.m_goalNode
			j.PushNewNode(newSuccessor, currentNode, Right, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newCol := col + jumpDistance
		givenCost := currentNode.m_givenCost + fixed_point_shift(jumpDistance)
		PathfindingNode * newSuccessor = &j.m_mapNodes[row][newCol]
		j.PushNewNode(newSuccessor, currentNode, Right, givenCost)
	}
}

func (j *JPSPlus) SearchUpRight(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Check for goal in general direction (straight line to Goal or Target Jump Point)
	if row > j.m_goalRow && col < j.m_goalCol {
		absJumpDistance = jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := row - j.m_goalRow
		diffCol := j.m_goalCol - col
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row - smallerDiff
			newCol := col + smallerDiff
			givenCost := currentNode.m_givenCost + (SQRT_2 * smallerDiff)
			newSuccessor := &j.m_mapNodes[newRow][newCol]
			j.PushNewNode(newSuccessor, currentNode, UpRight, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row - jumpDistance
		newCol := col + jumpDistance
		givenCost := currentNode.m_givenCost + (SQRT_2 * jumpDistance)
		newSuccessor := &j.m_mapNodes[newRow][newCol]
		j.PushNewNode(newSuccessor, currentNode, UpRight, givenCost)
	}
}

func (j *JPSPlus) SearchUp(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Consider straight line to Goal
	if col == j.m_goalCol && row > j.m_goalRow {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (row - absJumpDistance) <= j.m_goalRow {
			diff := row - j.m_goalRow
			givenCost := currentNode.m_givenCost + fixed_point_shift(diff)
			newSuccessor := j.m_goalNode
			j.PushNewNode(newSuccessor, currentNode, Up, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row - jumpDistance
		givenCost := currentNode.m_givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := &j.m_mapNodes[newRow][col]
		j.PushNewNode(newSuccessor, currentNode, Up, givenCost)
	}
}

func (j *JPSPlus) SearchUpLeft(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Check for goal in general direction (straight line to Goal or Target Jump Point)
	if row > j.m_goalRow && col > j.m_goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := row - j.m_goalRow
		diffCol := col - j.m_goalCol
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row - smallerDiff
			newCol := col - smallerDiff
			givenCost := currentNode.m_givenCost + int64(SQRT_2*smallerDiff)
			newSuccessor := &j.m_mapNodes[newRow][newCol]
			j.PushNewNode(newSuccessor, currentNode, UpLeft, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row - jumpDistance
		newCol := col - jumpDistance
		givenCost := currentNode.m_givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := &j.m_mapNodes[newRow][newCol]
		j.PushNewNode(newSuccessor, currentNode, UpLeft, givenCost)
	}
}

func (j *JPSPlus) SearchLeft(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Consider straight line to Goal
	if row == j.m_goalRow && col > j.m_goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (col - absJumpDistance) <= j.m_goalCol {
			diff := col - j.m_goalCol
			givenCost := currentNode.m_givenCost + FIXED_POINT_SHIFT(diff)
			newSuccessor := j.m_goalNode
			j.PushNewNode(newSuccessor, currentNode, Left, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newCol := col - jumpDistance
		givenCost := currentNode.m_givenCost + FIXED_POINT_SHIFT(jumpDistance)
		newSuccessor := &j.m_mapNodes[row][newCol]
		j.PushNewNode(newSuccessor, currentNode, Left, givenCost)
	}
}

func (j *JPSPlus) SearchDownLeft(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Check for goal in general direction (straight line to Goal or Target Jump Point)
	if row < j.m_goalRow && col > j.m_goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := j.m_goalRow - row
		diffCol := col - j.m_goalCol
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row + smallerDiff
			newCol := col - smallerDiff
			givenCost := currentNode.m_givenCost + int64(SQRT_2*smallerDiff)
			newSuccessor := &j.m_mapNodes[newRow][newCol]
			j.PushNewNode(newSuccessor, currentNode, DownLeft, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row + jumpDistance
		newCol := col - jumpDistance
		givenCost := currentNode.m_givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := &j.m_mapNodes[newRow][newCol]
		j.PushNewNode(newSuccessor, currentNode, DownLeft, givenCost)
	}
}

func (j *JPSPlus) PushNewNode(newSuccessor *PathfindingNode, currentNode *PathfindingNode, parentDirection int, givenCost int64) {
	if newSuccessor.m_iteration != j.m_currentIteration {
		// Place node on the Open list (we've never seen it before)

		// Compute heuristic using octile calculation (optimized: minDiff * SQRT_2_MINUS_ONE + maxDiff)
		diffrow := abs(j.m_goalRow - newSuccessor.m_row)
		diffcolumn := abs(j.m_goalCol - newSuccessor.m_col)
		var heuristicCost int64
		if diffrow <= diffcolumn {
			heuristicCost = int64(diffrow*SQRT_2_MINUS_ONE) + fixed_point_shift(diffcolumn)
		} else {
			heuristicCost = int64(diffcolumn*SQRT_2_MINUS_ONE) + fixed_point_shift(diffrow)
		}

		newSuccessor.m_parent = currentNode
		newSuccessor.m_directionFromParent = parentDirection
		newSuccessor.m_givenCost = givenCost
		newSuccessor.m_finalCost = givenCost + heuristicCost
		newSuccessor.m_listStatus = PathfindingNode_OnOpen
		newSuccessor.m_iteration = j.m_currentIteration

		if newSuccessor.m_finalCost <= currentNode.m_finalCost {
			m_fastStack.Push(newSuccessor)
		} else {
			m_simpleUnsortedPriorityQueue.Add(newSuccessor)
		}
	} else if givenCost < newSuccessor.m_givenCost && newSuccessor.m_listStatus == PathfindingNode_OnOpen { // Might be valid to remove this 2nd condition for extra speed (a node on the closed list wouldn't be cheaper)
		// We found a cheaper way to this node - update node

		// Extract heuristic cost (was previously calculated)
		heuristicCost := newSuccessor.m_finalCost - newSuccessor.m_givenCost

		newSuccessor.m_parent = currentNode
		newSuccessor.m_directionFromParent = parentDirection
		newSuccessor.m_givenCost = givenCost
		newSuccessor.m_finalCost = givenCost + heuristicCost

		// No decrease key operation necessary (already in unsorted open list)
	}
}
