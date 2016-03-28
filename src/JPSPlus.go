package jpsplus

const (
	Working = iota
	PathFound
	NoPathExists
)

const (
	SQRT_2                 = 3363
	SQRT_2_MINUS_ONE       = 985
	FIXED_POINT_MULTIPLIER = 100000
	FIXED_POINT_ONE
	FIXED_POINT_SQRT_2 = 141421
)

type LocJPS struct {
	Row int
	Col int
}

type JPSPlus struct {
	simpleUnsortedPriorityQueue *SimpleUnsortedPriorityQueue
	fastStack                   *FastStack
	currentIteration            int
	goalNode                    *PathfindingNode
	goalRow                     int
	goalCol                     int
}

func NewJPSPlus() *JPSPlus {
	j := new(JPSPlus)
	j.currentIteration = 1
	j.simpleUnsortedPriorityQueue = newSimpleUnsortedPriorityQueue(10000)
	j.fastStack = newFastStack(1000)
	return j
}

func (j *JPSPlus) GetPath(sRow, sCol, gRow, gCol int) (*[][2]int, bool) {
	// fmt.Println("GetPath")
	startRow := sRow
	startCol := sCol
	j.goalRow = gRow
	j.goalCol = gCol

	j.goalNode = newPathfindingNode(j.goalRow, j.goalCol)
	j.currentIteration += 1

	j.fastStack.Reset()
	j.simpleUnsortedPriorityQueue.Reset()

	startNode := newPathfindingNode(startRow, startCol)
	startNode.parent = nil
	startNode.givenCost = 0
	startNode.finalCost = 0
	startNode.listStatus = PathfindingNode_OnOpen
	startNode.iteration = j.currentIteration

	status := j.SearchLoop(startNode)
	// fmt.Printf("status == %#v\n", status)

	if status == PathFound {
		path := j.FinalizePath() //路径后续处理
		// fmt.Printf("path == %v\n", path)
		return path, true
	} else {
		return &([][2]int{}), false
	}
}

func (j *JPSPlus) SearchLoop(startNode *PathfindingNode) int {

	if startNode == j.goalNode {
		return PathFound
	}

	// jumpDistancesAndGoalBounds := &j.m_jumpDistancesAndGoalBounds[startNode.row][startNode.col]
	jump := DefaultJumpMap.Jump(startNode.row, startNode.col)

	JPSPlusExplore_AllDirections(startNode, jump, j)
	startNode.listStatus = PathfindingNode_OnClosed
	// fmt.Println(!j.simpleUnsortedPriorityQueue.Empty())
	// fmt.Println(!j.fastStack.Empty())

	for !j.simpleUnsortedPriorityQueue.Empty() || !j.fastStack.Empty() {
		var currentNode *PathfindingNode

		if !j.fastStack.Empty() {
			currentNode = j.fastStack.Pop()
		} else {
			currentNode = j.simpleUnsortedPriorityQueue.Pop()
		}

		if currentNode == j.goalNode {
			return PathFound
		}

		// jumpDistancesAndGoalBounds := &j.m_jumpDistancesAndGoalBounds[currentNode.row][currentNode.col]
		jump := DefaultJumpMap.Jump(currentNode.row, currentNode.col)
		exploreDirections[(jump.blocked*8)+currentNode.directionFromParent](currentNode, jump, j)

		currentNode.listStatus = PathfindingNode_OnClosed
	}
	return NoPathExists
}

func (j *JPSPlus) FinalizePath() *([][2]int) {
	var prevNode *PathfindingNode
	curNode := j.goalNode
	finalPath := make([][2]int, 0, 10000)

	for nil != curNode {
		loc := [2]int{curNode.row, curNode.col}

		if prevNode != nil {
			xDiff := curNode.col - prevNode.col
			yDiff := curNode.row - prevNode.row

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

			x := prevNode.col
			y := prevNode.row
			steps := xDiff - 1
			if yDiff > xDiff {
				steps = yDiff - 1
			}

			for i := 0; i < steps; i++ {
				x += xInc
				y += yInc

				locNew := [2]int{y, x}
				finalPath = append(finalPath, locNew)
			}
		}

		finalPath = append(finalPath, loc)
		prevNode = curNode
		curNode = curNode.parent
	}

	return reverse(&finalPath)
}

func reverse(path *([][2]int)) *([][2]int) {
	for i, j := 0, len(*path)-1; i < j; i, j = i+1, j-1 {
		(*path)[i], (*path)[j] = (*path)[j], (*path)[i]
	}
	return path
}

func MacroExploreDown(currentNode *PathfindingNode, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchDown(currentNode, jump.distant[Down])

	// if jpsPlus.goalRow >= jump.bounds.get(Down, MinRow) &&
	// 	jpsPlus.goalRow <= jump.bounds.get(Down, MaxRow) &&
	// 	jpsPlus.goalCol >= jump.bounds.get(Down, MinCol) &&
	// 	jpsPlus.goalCol <= jump.bounds.get(Down, MaxCol) {
	// 	jpsPlus.SearchDown(currentNode, jump.distant[Down])
	// }
}

func MacroExploreDownRight(currentNode *PathfindingNode, jump *Jump, jpsPlus *JPSPlus) {
	// if jpsPlus.goalRow >= jump.bounds.get(DownRight, MinRow) &&
	// 	jpsPlus.goalRow <= jump.bounds.get(DownRight, MaxRow) &&
	// 	jpsPlus.goalCol >= jump.bounds.get(DownRight, MinCol) &&
	// 	jpsPlus.goalCol <= jump.bounds.get(DownRight, MaxCol) {
	// 	jpsPlus.SearchDownRight(currentNode, jump.distant[DownRight])
	// }
	jpsPlus.SearchDownRight(currentNode, jump.distant[DownRight])

}

func MacroExploreRight(currentNode *PathfindingNode, jump *Jump, jpsPlus *JPSPlus) {

	// if jpsPlus.goalRow >= jump.bounds.get(Right, MinRow) &&
	// 	jpsPlus.goalRow <= jump.bounds.get(Right, MaxRow) &&
	// 	jpsPlus.goalCol >= jump.bounds.get(Right, MinCol) &&
	// 	jpsPlus.goalCol <= jump.bounds.get(Right, MaxCol) {
	// 	jpsPlus.SearchRight(currentNode, jump.distant[Right])
	// }
	jpsPlus.SearchRight(currentNode, jump.distant[Right])

}

func MacroExploreUpRight(currentNode *PathfindingNode, jump *Jump, jpsPlus *JPSPlus) {

	// if jpsPlus.goalRow >= jump.bounds.get(UpRight, MinRow) &&
	// 	jpsPlus.goalRow <= jump.bounds.get(UpRight, MaxRow) &&
	// 	jpsPlus.goalCol >= jump.bounds.get(UpRight, MinCol) &&
	// 	jpsPlus.goalCol <= jump.bounds.get(UpRight, MaxCol) {
	// 	jpsPlus.SearchUpRight(currentNode, jump.distant[UpRight])
	// }
	jpsPlus.SearchUpRight(currentNode, jump.distant[UpRight])

}

func MacroExploreUp(currentNode *PathfindingNode, jump *Jump, jpsPlus *JPSPlus) {

	// if jpsPlus.goalRow >= jump.bounds.get(Up, MinRow) &&
	// 	jpsPlus.goalRow <= jump.bounds.get(Up, MaxRow) &&
	// 	jpsPlus.goalCol >= jump.bounds.get(Up, MinCol) &&
	// 	jpsPlus.goalCol <= jump.bounds.get(Up, MaxCol) {
	// 	jpsPlus.SearchUp(currentNode, jump.distant[Up])
	// }
	jpsPlus.SearchUp(currentNode, jump.distant[Up])

}

func MacroExploreUpLeft(currentNode *PathfindingNode, jump *Jump, jpsPlus *JPSPlus) {

	// if jpsPlus.goalRow >= jump.bounds.get(UpLeft, MinRow) &&
	// 	jpsPlus.goalRow <= jump.bounds.get(UpLeft, MaxRow) &&
	// 	jpsPlus.goalCol >= jump.bounds.get(UpLeft, MinCol) &&
	// 	jpsPlus.goalCol <= jump.bounds.get(UpLeft, MaxCol) {
	// 	jpsPlus.SearchUpLeft(currentNode, jump.distant[UpLeft])
	// }
	jpsPlus.SearchUpLeft(currentNode, jump.distant[UpLeft])

}

func MacroExploreLeft(currentNode *PathfindingNode, jump *Jump, jpsPlus *JPSPlus) {

	// if jpsPlus.goalRow >= jump.bounds.get(Left, MinRow) &&
	// 	jpsPlus.goalRow <= jump.bounds.get(Left, MaxRow) &&
	// 	jpsPlus.goalCol >= jump.bounds.get(Left, MinCol) &&
	// 	jpsPlus.goalCol <= jump.bounds.get(Left, MaxCol) {
	// 	jpsPlus.SearchLeft(currentNode, jump.distant[Left])
	// }
	jpsPlus.SearchLeft(currentNode, jump.distant[Left])

}

func MacroExploreDownLeft(currentNode *PathfindingNode, jump *Jump, jpsPlus *JPSPlus) {

	// if jpsPlus.goalRow >= jump.bounds.get(DownLeft, MinRow) &&
	// 	jpsPlus.goalRow <= jump.bounds.get(DownLeft, MaxRow) &&
	// 	jpsPlus.goalCol >= jump.bounds.get(DownLeft, MinCol) &&
	// 	jpsPlus.goalCol <= jump.bounds.get(DownLeft, MaxCol) {
	// 	jpsPlus.SearchDownLeft(currentNode, jump.distant[DownLeft])
	// }
	jpsPlus.SearchDownLeft(currentNode, jump.distant[DownLeft])

}

func JPSPlusExplore_Null(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	// Purposely does nothing
}

func JPSPlusExplore_D(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_DR(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
}

func JPSPlusExplore_R(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_UR(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_UL(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_DL(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
}

func JPSPlusExplore_D_DR(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
}

func JPSPlusExplore_DR_R(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_UR(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_UR_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_U_UL(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_UL_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_DL(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
}

func JPSPlusExplore_DL_D(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_D_R(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_U_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_D(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_D_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_R_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_D_DR_R(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_DR_R_UR(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_UR_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_UR_U_UL(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UL_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_UL_L_DL(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_DL_D(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_DL_D_DR(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
}

func JPSPlusExplore_D_R_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_R_U_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_L_D(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_D_R(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_DR_D_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_R_D_DL_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UR_R_D(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_U_R_DR_D(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_UL_U_R(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_L_U_UR_R(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_D_DL_L_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_D_L_UL_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_R_DR_D_DL_L(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UR_R_DR_D(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_UL_U_UR_R(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_D_DL_L_UL_U(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_AllDirections(currentNode *PathfindingNode, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
}

func (j *JPSPlus) SearchDown(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if col == j.goalCol && row < j.goalRow {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (row + absJumpDistance) >= j.goalRow {
			diff := j.goalRow - row
			givenCost := currentNode.givenCost + fixed_point_shift(diff)
			newSuccessor := j.goalNode
			j.PushNewNode(newSuccessor, currentNode, Down, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row + jumpDistance
		givenCost := currentNode.givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := newPathfindingNode(newRow, col)
		j.PushNewNode(newSuccessor, currentNode, Down, givenCost)
	}
}

func (j *JPSPlus) SearchDownRight(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row < j.goalRow && col < j.goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := j.goalRow - row
		diffCol := j.goalCol - col
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row + smallerDiff
			newCol := col + smallerDiff
			givenCost := currentNode.givenCost + int64(SQRT_2*smallerDiff)
			newSuccessor := newPathfindingNode(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, DownRight, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row + jumpDistance
		newCol := col + jumpDistance
		givenCost := currentNode.givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := newPathfindingNode(newRow, newCol)
		j.PushNewNode(newSuccessor, currentNode, DownRight, givenCost)
	}
}

func (j *JPSPlus) SearchRight(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row == j.goalRow && col < j.goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (col + absJumpDistance) >= j.goalCol {
			diff := j.goalCol - col
			givenCost := currentNode.givenCost + fixed_point_shift(diff)
			newSuccessor := j.goalNode
			j.PushNewNode(newSuccessor, currentNode, Right, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newCol := col + jumpDistance
		givenCost := currentNode.givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := newPathfindingNode(row, newCol)
		j.PushNewNode(newSuccessor, currentNode, Right, givenCost)
	}
}

func (j *JPSPlus) SearchUpRight(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row > j.goalRow && col < j.goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := row - j.goalRow
		diffCol := j.goalCol - col
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row - smallerDiff
			newCol := col + smallerDiff
			givenCost := currentNode.givenCost + int64(SQRT_2*smallerDiff)
			newSuccessor := newPathfindingNode(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, UpRight, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row - jumpDistance
		newCol := col + jumpDistance
		givenCost := currentNode.givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := newPathfindingNode(newRow, newCol)
		j.PushNewNode(newSuccessor, currentNode, UpRight, givenCost)
	}
}

func (j *JPSPlus) SearchUp(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if col == j.goalCol && row > j.goalRow {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (row - absJumpDistance) <= j.goalRow {
			diff := row - j.goalRow
			givenCost := currentNode.givenCost + fixed_point_shift(diff)
			newSuccessor := j.goalNode
			j.PushNewNode(newSuccessor, currentNode, Up, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row - jumpDistance
		givenCost := currentNode.givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := newPathfindingNode(newRow, col)
		j.PushNewNode(newSuccessor, currentNode, Up, givenCost)
	}
}

func (j *JPSPlus) SearchUpLeft(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row > j.goalRow && col > j.goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := row - j.goalRow
		diffCol := col - j.goalCol
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row - smallerDiff
			newCol := col - smallerDiff
			givenCost := currentNode.givenCost + int64(SQRT_2*smallerDiff)
			newSuccessor := newPathfindingNode(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, UpLeft, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row - jumpDistance
		newCol := col - jumpDistance
		givenCost := currentNode.givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := newPathfindingNode(newRow, newCol)
		j.PushNewNode(newSuccessor, currentNode, UpLeft, givenCost)
	}
}

func (j *JPSPlus) SearchLeft(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row == j.goalRow && col > j.goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (col - absJumpDistance) <= j.goalCol {
			diff := col - j.goalCol
			givenCost := currentNode.givenCost + fixed_point_shift(diff)
			newSuccessor := j.goalNode
			j.PushNewNode(newSuccessor, currentNode, Left, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newCol := col - jumpDistance
		givenCost := currentNode.givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := newPathfindingNode(row, newCol)
		j.PushNewNode(newSuccessor, currentNode, Left, givenCost)
	}
}

func (j *JPSPlus) SearchDownLeft(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row < j.goalRow && col > j.goalCol {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := j.goalRow - row
		diffCol := col - j.goalCol
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row + smallerDiff
			newCol := col - smallerDiff
			givenCost := currentNode.givenCost + int64(SQRT_2*smallerDiff)
			newSuccessor := newPathfindingNode(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, DownLeft, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row + jumpDistance
		newCol := col - jumpDistance
		givenCost := currentNode.givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := newPathfindingNode(newRow, newCol)
		j.PushNewNode(newSuccessor, currentNode, DownLeft, givenCost)
	}
}

func (j *JPSPlus) PushNewNode(newSuccessor *PathfindingNode, currentNode *PathfindingNode, parentDirection int, givenCost int64) {
	if newSuccessor.iteration != j.currentIteration {

		diffrow := abs(j.goalRow - newSuccessor.row)
		diffcolumn := abs(j.goalCol - newSuccessor.col)
		var heuristicCost int64
		if diffrow <= diffcolumn {
			heuristicCost = int64(diffrow*SQRT_2_MINUS_ONE) + fixed_point_shift(diffcolumn)
		} else {
			heuristicCost = int64(diffcolumn*SQRT_2_MINUS_ONE) + fixed_point_shift(diffrow)
		}

		newSuccessor.parent = currentNode
		newSuccessor.directionFromParent = parentDirection
		newSuccessor.givenCost = givenCost
		newSuccessor.finalCost = givenCost + heuristicCost
		newSuccessor.listStatus = PathfindingNode_OnOpen
		newSuccessor.iteration = j.currentIteration

		if newSuccessor.finalCost <= currentNode.finalCost {
			j.fastStack.Push(newSuccessor)
		} else {
			j.simpleUnsortedPriorityQueue.Add(newSuccessor)
		}
	} else if givenCost < newSuccessor.givenCost && newSuccessor.listStatus == PathfindingNode_OnOpen {
		heuristicCost := newSuccessor.finalCost - newSuccessor.givenCost

		newSuccessor.parent = currentNode
		newSuccessor.directionFromParent = parentDirection
		newSuccessor.givenCost = givenCost
		newSuccessor.finalCost = givenCost + heuristicCost

	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

func fixed_point_shift(x int) int64 {
	return int64(x) * int64(FIXED_POINT_MULTIPLIER)
}
