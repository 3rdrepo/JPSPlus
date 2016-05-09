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
	MapWidth  = 100
	MapHeight = 100
)

const (
	COST_STRAIGHT = 1000
	COST_DIAGONAL = 1414
)

type LocJPS struct {
	Row int
	Col int
}

type Node struct {
	row       int
	col       int
	givenCost int
	finalCost int
	parent    *Node
	direction int
	heapIndex int
}

func newNode(r int, c int) *Node {
	node := new(Node)
	node.row = r
	node.col = c
	return node
}

func (this Node) equal(n *Node) bool {
	return this.row == n.row && this.col == n.col
}

type nodeList map[int]*Node

func newNodeList() nodeList {
	return make(nodeList)
}

func (n nodeList) insert(node *Node) {
	n[node.row*MapWidth+node.col] = node
}

func (n nodeList) lookup(row int, col int) (node *Node, ok bool) {
	node, ok = n[row*MapWidth+col]
	return
}

func (n nodeList) remove(row int, col int) {
	delete(n, row*MapWidth+col)
}

func (n nodeList) len() int {
	return len(n)
}

type JPSPlus struct {
	startNode *Node
	goalNode  *Node
	fastStack *PriorityQueue
	openSet   nodeList
	closeSet  nodeList
}

func NewJPSPlus(sRow, sCol, gRow, gCol int) *JPSPlus {
	j := new(JPSPlus)
	j.startNode = newNode(sRow, sCol)
	j.goalNode = newNode(gRow, gCol)
	j.fastStack = newPriorityQueue()
	j.closeSet = newNodeList()
	j.openSet = newNodeList()
	j.closeSet.insert(j.startNode)

	return j
}

func (this *JumpMap) GetPath(sRow, sCol, gRow, gCol int) (path []LocJPS, isFind bool) {
	if sRow == gRow && sCol == gCol {
		isFind = true
	} else {
		jps := NewJPSPlus(sRow, sCol, gRow, gCol)
		status := this.SearchLoop(jps)
		if status == PathFound {
			path = jps.FinalizePath()
			isFind = true
		}
	}
	return
}

func (this *JumpMap) SearchLoop(jps *JPSPlus) int {
	jumpNode := this.Jump(jps.startNode.row, jps.startNode.col)
	JPSPlusExplore_AllDirections(jps.startNode, jumpNode, jps)

	for 0 != jps.fastStack.Len() {
		cur := jps.fastStack.PopNode()
		jps.openSet.remove(cur.row, cur.col)
		if cur.equal(jps.goalNode) {
			return PathFound
		} else {
			jps.closeSet.insert(cur)
			jump := this.Jump(cur.row, cur.col)
			exploreDirections[(jump.blocked*8)+cur.direction](cur, jump, jps)
		}
	}
	return NoPathExists
}

func (j *JPSPlus) FinalizePath() (path []LocJPS) {
	// var prevNode *Node
	// curNode := j.goalNode

	// for nil != curNode {
	// 	loc := LocJPS{curNode.row, curNode.col}

	// 	if prevNode != nil {
	// 		xDiff := curNode.col - prevNode.col
	// 		yDiff := curNode.row - prevNode.row

	// 		xInc := 0
	// 		yInc := 0

	// 		if xDiff > 0 {
	// 			xInc = 1
	// 		} else if xDiff < 0 {
	// 			xInc = -1
	// 			xDiff = -xDiff
	// 		}

	// 		if yDiff > 0 {
	// 			yInc = 1
	// 		} else if yDiff < 0 {
	// 			yInc = -1
	// 			yDiff = -yDiff
	// 		}

	// 		x := prevNode.col
	// 		y := prevNode.row
	// 		steps := xDiff - 1
	// 		if yDiff > xDiff {
	// 			steps = yDiff - 1
	// 		}

	// 		for i := 0; i < steps; i++ {
	// 			x += xInc
	// 			y += yInc

	// 			locNew := LocJPS{y, x}
	// 			path = append(path, locNew)
	// 		}
	// 	}

	// 	path = append(path, loc)
	// 	prevNode = curNode
	// 	curNode = curNode.parent
	// }
	fmt.Println(j.openSet.len(), j.closeSet.len(), j.fastStack.Len())
	for cur := j.goalNode; nil != cur; cur = cur.parent {
		path = append(path, LocJPS{cur.row, cur.col})
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return
}

func MacroExploreDown(currentNode *Node, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchDown(currentNode, jump.distant[Down])
}

func MacroExploreDownRight(currentNode *Node, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchDownRight(currentNode, jump.distant[DownRight])

}

func MacroExploreRight(currentNode *Node, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchRight(currentNode, jump.distant[Right])

}

func MacroExploreUpRight(currentNode *Node, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchUpRight(currentNode, jump.distant[UpRight])

}

func MacroExploreUp(currentNode *Node, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchUp(currentNode, jump.distant[Up])

}

func MacroExploreUpLeft(currentNode *Node, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchUpLeft(currentNode, jump.distant[UpLeft])

}

func MacroExploreLeft(currentNode *Node, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchLeft(currentNode, jump.distant[Left])

}

func MacroExploreDownLeft(currentNode *Node, jump *Jump, jpsPlus *JPSPlus) {
	jpsPlus.SearchDownLeft(currentNode, jump.distant[DownLeft])

}

func JPSPlusExplore_Null(currentNode *Node, jump *Jump, jps *JPSPlus) {
}

func JPSPlusExplore_D(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_DR(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
}

func JPSPlusExplore_R(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_UR(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_UL(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_DL(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
}

func JPSPlusExplore_D_DR(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
}

func JPSPlusExplore_DR_R(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_UR(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_UR_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_U_UL(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_UL_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_DL(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
}

func JPSPlusExplore_DL_D(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_D_R(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_U_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_D(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_D_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_R_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_D_DR_R(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_DR_R_UR(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_UR_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_UR_U_UL(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UL_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_UL_L_DL(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
}

func JPSPlusExplore_L_DL_D(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_DL_D_DR(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
}

func JPSPlusExplore_D_R_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_R_U_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_L_D(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_D_R(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_R_DR_D_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_R_D_DL_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UR_R_D(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_U_R_DR_D(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_UL_U_R(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_L_U_UR_R(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_D_DL_L_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_D_L_UL_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_R_DR_D_DL_L(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
}

func JPSPlusExplore_U_UR_R_DR_D(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
	MacroExploreDown(currentNode, jump, jps)
}

func JPSPlusExplore_L_UL_U_UR_R(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
}

func JPSPlusExplore_D_DL_L_UL_U(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
}

func JPSPlusExplore_AllDirections(currentNode *Node, jump *Jump, jps *JPSPlus) {
	MacroExploreDown(currentNode, jump, jps)
	MacroExploreDownLeft(currentNode, jump, jps)
	MacroExploreLeft(currentNode, jump, jps)
	MacroExploreUpLeft(currentNode, jump, jps)
	MacroExploreUp(currentNode, jump, jps)
	MacroExploreUpRight(currentNode, jump, jps)
	MacroExploreRight(currentNode, jump, jps)
	MacroExploreDownRight(currentNode, jump, jps)
}

func (j *JPSPlus) SearchDown(currentNode *Node, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if col == j.goalNode.col && row < j.goalNode.row {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (row + absJumpDistance) >= j.goalNode.row {
			diff := j.goalNode.row - row
			givenCost := currentNode.givenCost + fixed_point_shift(diff)
			newSuccessor := j.goalNode
			j.PushNewNode(newSuccessor, currentNode, Down, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row + jumpDistance
		givenCost := currentNode.givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := newNode(newRow, col)
		j.PushNewNode(newSuccessor, currentNode, Down, givenCost)
	}
}

func (j *JPSPlus) SearchDownRight(currentNode *Node, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row < j.goalNode.row && col < j.goalNode.col {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := j.goalNode.row - row
		diffCol := j.goalNode.col - col
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row + smallerDiff
			newCol := col + smallerDiff
			givenCost := currentNode.givenCost + COST_DIAGONAL*smallerDiff
			newSuccessor := newNode(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, DownRight, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row + jumpDistance
		newCol := col + jumpDistance
		givenCost := currentNode.givenCost + COST_DIAGONAL*jumpDistance
		newSuccessor := newNode(newRow, newCol)
		j.PushNewNode(newSuccessor, currentNode, DownRight, givenCost)
	}
}

func (j *JPSPlus) SearchRight(currentNode *Node, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row == j.goalNode.row && col < j.goalNode.col {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (col + absJumpDistance) >= j.goalNode.col {
			diff := j.goalNode.col - col
			givenCost := currentNode.givenCost + fixed_point_shift(diff)
			newSuccessor := j.goalNode
			j.PushNewNode(newSuccessor, currentNode, Right, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newCol := col + jumpDistance
		givenCost := currentNode.givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := newNode(row, newCol)
		j.PushNewNode(newSuccessor, currentNode, Right, givenCost)
	}
}

func (j *JPSPlus) SearchUpRight(currentNode *Node, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row > j.goalNode.row && col < j.goalNode.col {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := row - j.goalNode.row
		diffCol := j.goalNode.col - col
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row - smallerDiff
			newCol := col + smallerDiff
			givenCost := currentNode.givenCost + COST_DIAGONAL*smallerDiff
			newSuccessor := newNode(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, UpRight, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row - jumpDistance
		newCol := col + jumpDistance
		givenCost := currentNode.givenCost + COST_DIAGONAL*jumpDistance
		newSuccessor := newNode(newRow, newCol)
		j.PushNewNode(newSuccessor, currentNode, UpRight, givenCost)
	}
}

func (j *JPSPlus) SearchUp(currentNode *Node, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if col == j.goalNode.col && row > j.goalNode.row {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (row - absJumpDistance) <= j.goalNode.row {
			diff := row - j.goalNode.row
			givenCost := currentNode.givenCost + fixed_point_shift(diff)
			newSuccessor := j.goalNode
			j.PushNewNode(newSuccessor, currentNode, Up, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row - jumpDistance
		givenCost := currentNode.givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := newNode(newRow, col)
		j.PushNewNode(newSuccessor, currentNode, Up, givenCost)
	}
}

func (j *JPSPlus) SearchUpLeft(currentNode *Node, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row > j.goalNode.row && col > j.goalNode.col {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := row - j.goalNode.row
		diffCol := col - j.goalNode.col
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row - smallerDiff
			newCol := col - smallerDiff
			givenCost := currentNode.givenCost + COST_DIAGONAL*smallerDiff
			newSuccessor := newNode(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, UpLeft, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row - jumpDistance
		newCol := col - jumpDistance
		givenCost := currentNode.givenCost + COST_DIAGONAL*jumpDistance
		newSuccessor := newNode(newRow, newCol)
		j.PushNewNode(newSuccessor, currentNode, UpLeft, givenCost)
	}
}

func (j *JPSPlus) SearchLeft(currentNode *Node, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row == j.goalNode.row && col > j.goalNode.col {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		if (col - absJumpDistance) <= j.goalNode.col {
			diff := col - j.goalNode.col
			givenCost := currentNode.givenCost + fixed_point_shift(diff)
			newSuccessor := j.goalNode
			j.PushNewNode(newSuccessor, currentNode, Left, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newCol := col - jumpDistance
		givenCost := currentNode.givenCost + fixed_point_shift(jumpDistance)
		newSuccessor := newNode(row, newCol)
		j.PushNewNode(newSuccessor, currentNode, Left, givenCost)
	}
}

func (j *JPSPlus) SearchDownLeft(currentNode *Node, jumpDistance int) {
	row := currentNode.row
	col := currentNode.col

	if row < j.goalNode.row && col > j.goalNode.col {
		absJumpDistance := jumpDistance
		if absJumpDistance < 0 {
			absJumpDistance = -absJumpDistance
		}

		diffRow := j.goalNode.row - row
		diffCol := col - j.goalNode.col
		smallerDiff := diffRow
		if diffCol < smallerDiff {
			smallerDiff = diffCol
		}

		if smallerDiff <= absJumpDistance {
			newRow := row + smallerDiff
			newCol := col - smallerDiff
			givenCost := currentNode.givenCost + COST_DIAGONAL*smallerDiff
			newSuccessor := newNode(newRow, newCol)
			j.PushNewNode(newSuccessor, currentNode, DownLeft, givenCost)
			return
		}
	}

	if jumpDistance > 0 {
		newRow := row + jumpDistance
		newCol := col - jumpDistance
		givenCost := currentNode.givenCost + COST_DIAGONAL*jumpDistance
		newSuccessor := newNode(newRow, newCol)
		j.PushNewNode(newSuccessor, currentNode, DownLeft, givenCost)
	}
}

func (j *JPSPlus) PushNewNode(newJump *Node, cur *Node, parentDirection int, givenCost int) {
	_, ok := j.closeSet.lookup(newJump.row, newJump.col)
	if !ok {
		old, ok := j.openSet.lookup(newJump.row, newJump.col)
		if ok {
			if old.givenCost > givenCost {
				old.givenCost = givenCost
				old.parent = cur
				old.finalCost = givenCost + heuristicDistance(newJump, j.goalNode)
			}
		} else {
			newJump.parent = cur
			newJump.direction = parentDirection
			newJump.givenCost = givenCost
			newJump.finalCost = givenCost + heuristicDistance(newJump, j.goalNode)
			j.openSet.insert(newJump)
			j.fastStack.PushNode(newJump)
		}
	}
}

func heuristicDistance(cur *Node, stop *Node) int {
	row := abs(cur.row - stop.row)
	col := abs(cur.col - stop.col)
	h_dia := min(row, col)
	h_str := abs(row - col)
	return COST_DIAGONAL*h_dia + COST_STRAIGHT*h_str
}

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func fixed_point_shift(x int) int {
	return x * COST_STRAIGHT
}
