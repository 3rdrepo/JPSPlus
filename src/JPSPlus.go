package jpsplus

import (
	"fmt"
)

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

type JPSPlus struct {
	// m_width                       int
	// m_height                      int
	// m_jumpDistancesAndGoalBounds  [][]JumpDistancesAndGoalBounds
	m_simpleUnsortedPriorityQueue *SimpleUnsortedPriorityQueue
	m_fastStack                   *FastStack
	m_mapNodes                    *PathfindingNodeMap
	m_currentIteration            int
	m_goalNode                    *PathfindingNode
	m_goalRow                     int
	m_goalCol                     int
}

func newJPSPlus() *JPSPlus {
	width := DefaultJumpDistancesAndGoalBounds.width()
	height := DefaultJumpDistancesAndGoalBounds.height()

	j := new(JPSPlus)
	j.m_currentIteration = 1

	j.m_simpleUnsortedPriorityQueue = newSimpleUnsortedPriorityQueue(10000)
	j.m_fastStack = newFastStack(1000)
	j.m_mapNodes = newPathfindingNodeMap(width, height)

	return j
}

func (j *JPSPlus) GetPath(s xyLocJPS, g xyLocJPS) ([]xyLocJPS, bool) {
	fmt.Println("GetPath")
	startRow := s.y
	startCol := s.x
	j.m_goalRow = g.y
	j.m_goalCol = g.x

	// Initialize map
	j.m_goalNode = j.m_mapNodes.get(j.m_goalRow, j.m_goalCol)
	j.m_currentIteration += 1

	j.m_fastStack.Reset()
	j.m_simpleUnsortedPriorityQueue.Reset()

	// Create starting node
	startNode := j.m_mapNodes.get(startRow, startCol)
	startNode.m_parent = nil
	startNode.m_givenCost = 0
	startNode.m_finalCost = 0
	startNode.m_listStatus = PathfindingNode_OnOpen
	startNode.m_iteration = j.m_currentIteration

	// Actual search
	status := j.SearchLoop(startNode)
	// fmt.Printf("status == %#v\n", status)

	if status == PathFound {
		// fmt.Printf("jps == %#v\n", j.m_mapNodes)
		path := j.FinalizePath() //路径后续处理
		// fmt.Printf("path == %v\n", path)
		return path, true
	} else {
		// No path
		return []xyLocJPS{}, false
	}
}

func (j *JPSPlus) SearchLoop(startNode *PathfindingNode) int {
	// Special case for the starting node

	if startNode == j.m_goalNode {
		return PathFound
	}

	// jumpDistancesAndGoalBounds := &j.m_jumpDistancesAndGoalBounds[startNode.m_row][startNode.m_col]
	jump := DefaultJumpDistancesAndGoalBounds.get(startNode.m_row, startNode.m_col)

	JPSPlusExplore_AllDirections(startNode, jump, j)
	startNode.m_listStatus = PathfindingNode_OnClosed
	// fmt.Println(!j.m_simpleUnsortedPriorityQueue.Empty())
	// fmt.Println(!j.m_fastStack.Empty())

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
		jump := DefaultJumpDistancesAndGoalBounds.get(currentNode.m_row, currentNode.m_col)
		exploreDirections[(jump.blockedDirectionBitfield*8)+currentNode.m_directionFromParent](currentNode, jump, j)

		currentNode.m_listStatus = PathfindingNode_OnClosed
	}
	return NoPathExists
}

func (j *JPSPlus) FinalizePath() []xyLocJPS {
	var prevNode *PathfindingNode
	curNode := j.m_goalNode
	finalPath := make([]xyLocJPS, 0, 100)

	for nil != curNode {
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

		finalPath = append(finalPath, loc)
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

func MacroExploreDown(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jpsPlus *JPSPlus) {
	jpsPlus.SearchDown(currentNode, jump.jumpDistance.get(Down))

	// if jpsPlus.m_goalRow >= jump.bounds.get(Down, MinRow) &&
	// 	jpsPlus.m_goalRow <= jump.bounds.get(Down, MaxRow) &&
	// 	jpsPlus.m_goalCol >= jump.bounds.get(Down, MinCol) &&
	// 	jpsPlus.m_goalCol <= jump.bounds.get(Down, MaxCol) {
	// 	jpsPlus.SearchDown(currentNode, jump.jumpDistance.get(Down))
	// }
}

func MacroExploreDownRight(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jpsPlus *JPSPlus) {
	// if jpsPlus.m_goalRow >= jump.bounds.get(DownRight, MinRow) &&
	// 	jpsPlus.m_goalRow <= jump.bounds.get(DownRight, MaxRow) &&
	// 	jpsPlus.m_goalCol >= jump.bounds.get(DownRight, MinCol) &&
	// 	jpsPlus.m_goalCol <= jump.bounds.get(DownRight, MaxCol) {
	// 	jpsPlus.SearchDownRight(currentNode, jump.jumpDistance.get(DownRight))
	// }
	jpsPlus.SearchDownRight(currentNode, jump.jumpDistance.get(DownRight))

}

func MacroExploreRight(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jpsPlus *JPSPlus) {

	// if jpsPlus.m_goalRow >= jump.bounds.get(Right, MinRow) &&
	// 	jpsPlus.m_goalRow <= jump.bounds.get(Right, MaxRow) &&
	// 	jpsPlus.m_goalCol >= jump.bounds.get(Right, MinCol) &&
	// 	jpsPlus.m_goalCol <= jump.bounds.get(Right, MaxCol) {
	// 	jpsPlus.SearchRight(currentNode, jump.jumpDistance.get(Right))
	// }
	jpsPlus.SearchRight(currentNode, jump.jumpDistance.get(Right))

}

func MacroExploreUpRight(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jpsPlus *JPSPlus) {

	// if jpsPlus.m_goalRow >= jump.bounds.get(UpRight, MinRow) &&
	// 	jpsPlus.m_goalRow <= jump.bounds.get(UpRight, MaxRow) &&
	// 	jpsPlus.m_goalCol >= jump.bounds.get(UpRight, MinCol) &&
	// 	jpsPlus.m_goalCol <= jump.bounds.get(UpRight, MaxCol) {
	// 	jpsPlus.SearchUpRight(currentNode, jump.jumpDistance.get(UpRight))
	// }
	jpsPlus.SearchUpRight(currentNode, jump.jumpDistance.get(UpRight))

}

func MacroExploreUp(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jpsPlus *JPSPlus) {

	// if jpsPlus.m_goalRow >= jump.bounds.get(Up, MinRow) &&
	// 	jpsPlus.m_goalRow <= jump.bounds.get(Up, MaxRow) &&
	// 	jpsPlus.m_goalCol >= jump.bounds.get(Up, MinCol) &&
	// 	jpsPlus.m_goalCol <= jump.bounds.get(Up, MaxCol) {
	// 	jpsPlus.SearchUp(currentNode, jump.jumpDistance.get(Up))
	// }
	jpsPlus.SearchUp(currentNode, jump.jumpDistance.get(Up))

}

func MacroExploreUpLeft(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jpsPlus *JPSPlus) {

	// if jpsPlus.m_goalRow >= jump.bounds.get(UpLeft, MinRow) &&
	// 	jpsPlus.m_goalRow <= jump.bounds.get(UpLeft, MaxRow) &&
	// 	jpsPlus.m_goalCol >= jump.bounds.get(UpLeft, MinCol) &&
	// 	jpsPlus.m_goalCol <= jump.bounds.get(UpLeft, MaxCol) {
	// 	jpsPlus.SearchUpLeft(currentNode, jump.jumpDistance.get(UpLeft))
	// }
	jpsPlus.SearchUpLeft(currentNode, jump.jumpDistance.get(UpLeft))

}

func MacroExploreLeft(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jpsPlus *JPSPlus) {

	// if jpsPlus.m_goalRow >= jump.bounds.get(Left, MinRow) &&
	// 	jpsPlus.m_goalRow <= jump.bounds.get(Left, MaxRow) &&
	// 	jpsPlus.m_goalCol >= jump.bounds.get(Left, MinCol) &&
	// 	jpsPlus.m_goalCol <= jump.bounds.get(Left, MaxCol) {
	// 	jpsPlus.SearchLeft(currentNode, jump.jumpDistance.get(Left))
	// }
	jpsPlus.SearchLeft(currentNode, jump.jumpDistance.get(Left))

}

func MacroExploreDownLeft(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jpsPlus *JPSPlus) {
	// fmt.Println("MacroExploreDownLeft")

	// if jpsPlus.m_goalRow >= jump.bounds.get(DownLeft, MinRow) &&
	// 	jpsPlus.m_goalRow <= jump.bounds.get(DownLeft, MaxRow) &&
	// 	jpsPlus.m_goalCol >= jump.bounds.get(DownLeft, MinCol) &&
	// 	jpsPlus.m_goalCol <= jump.bounds.get(DownLeft, MaxCol) {
	// 	jpsPlus.SearchDownLeft(currentNode, jump.jumpDistance.get(DownLeft))
	// }
	jpsPlus.SearchDownLeft(currentNode, jump.jumpDistance.get(DownLeft))

}

func JPSPlusExplore_Null(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	// Purposely does nothing
}

func JPSPlusExplore_D(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_DR(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
}

func JPSPlusExplore_R(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_UR(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_UL(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_DL(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
}

// Adjacent Doubles

func JPSPlusExplore_D_DR(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
}

func JPSPlusExplore_DR_R(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_UR(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_UR_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_U_UL(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_UL_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_DL(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
}

func JPSPlusExplore_DL_D(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

// Non-Adjacent Cardinal Doubles

func JPSPlusExplore_D_R(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_U_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_D(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_D_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_R_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

// Adjacent Triples

func JPSPlusExplore_D_DR_R(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_DR_R_UR(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_UR_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_UR_U_UL(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UL_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_UL_L_DL(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_DL_D(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_DL_D_DR(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
}

// Non-Adjacent Cardinal Triples

func JPSPlusExplore_D_R_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_R_U_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_L_D(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_D_R(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

// Quads

func JPSPlusExplore_R_DR_D_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_R_D_DL_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UR_R_D(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_U_R_DR_D(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_UL_U_R(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_L_U_UR_R(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_D_DL_L_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_D_L_UL_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

// Quints

func JPSPlusExplore_R_DR_D_DL_L(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UR_R_DR_D(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_UL_U_UR_R(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_D_DL_L_UL_U(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_AllDirections(currentNode *PathfindingNode, jump *JumpDistancesAndGoalBounds, jps *JPSPlus) {
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
		newSuccessor := j.m_mapNodes.get(newRow, col)
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
			newSuccessor := j.m_mapNodes.get(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, DownRight, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row + jumpDistance
		newCol := col + jumpDistance
		givenCost := currentNode.m_givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := j.m_mapNodes.get(newRow, newCol)
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
		newSuccessor := j.m_mapNodes.get(row, newCol)
		j.PushNewNode(newSuccessor, currentNode, Right, givenCost)
	}
}

func (j *JPSPlus) SearchUpRight(currentNode *PathfindingNode, jumpDistance int) {
	row := currentNode.m_row
	col := currentNode.m_col

	// Check for goal in general direction (straight line to Goal or Target Jump Point)
	if row > j.m_goalRow && col < j.m_goalCol {
		absJumpDistance := jumpDistance
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
			givenCost := currentNode.m_givenCost + int64(SQRT_2*smallerDiff)
			newSuccessor := j.m_mapNodes.get(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, UpRight, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row - jumpDistance
		newCol := col + jumpDistance
		givenCost := currentNode.m_givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := j.m_mapNodes.get(newRow, newCol)
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
		newSuccessor := j.m_mapNodes.get(newRow, col)
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
			newSuccessor := j.m_mapNodes.get(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, UpLeft, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row - jumpDistance
		newCol := col - jumpDistance
		givenCost := currentNode.m_givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := j.m_mapNodes.get(newRow, newCol)
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
			givenCost := currentNode.m_givenCost + fixed_point_shift(diff)
			newSuccessor := j.m_goalNode
			j.PushNewNode(newSuccessor, currentNode, Left, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newCol := col - jumpDistance
		givenCost := currentNode.m_givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := j.m_mapNodes.get(row, newCol)
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
			newSuccessor := j.m_mapNodes.get(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, DownLeft, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		// Directly jump
		newRow := row + jumpDistance
		newCol := col - jumpDistance
		givenCost := currentNode.m_givenCost + int64(SQRT_2*jumpDistance)
		newSuccessor := j.m_mapNodes.get(newRow, newCol)
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
			j.m_fastStack.Push(newSuccessor)
		} else {
			j.m_simpleUnsortedPriorityQueue.Add(newSuccessor)
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

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
