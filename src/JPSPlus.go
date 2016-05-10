package jpsplus

const (
	MapWidth = 100
	MapHeight
	TileWidth = 100
	TileHeight
)

const (
	COST_STRAIGHT = 1000
	COST_DIAGONAL = 1414
	CloseSetMax   = 200
)

type LocJPS struct {
	row int
	col int
}

func newLocJPS(row int, col int) *LocJPS {
	l := new(LocJPS)
	l.row = row
	l.col = col
	return l
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

func (this *Node) equal(n *Node) bool {
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

func NewJPSPlus(sRow int, sCol int, gRow int, gCol int) *JPSPlus {
	j := new(JPSPlus)
	j.startNode = newNode(sRow, sCol)
	j.goalNode = newNode(gRow, gCol)
	j.fastStack = newPriorityQueue()
	j.closeSet = newNodeList()
	j.openSet = newNodeList()
	j.closeSet.insert(j.startNode)

	return j
}

func (this *JumpMap) GetPath(sRow int, sCol int, gRow int, gCol int) (path map[int]*LocJPS, isFind bool) {
	jps := NewJPSPlus(sRow, sCol, gRow, gCol)
	if this.SearchLoop(jps) {
		path = jps.FinalizePath()
		isFind = true
	}
	return
}

func logicToTile(x int, y int) (r int, c int) {
	r = y / TileHeight
	c = x / TileWidth
	return
}

func tileToLogic(r int, c int) (x int, y int) {
	y = r * TileHeight
	x = c * TileWidth
	return
}

func (this *JumpMap) SearchLoop(jps *JPSPlus) bool {
	jumpNode := this[jps.startNode.row][jps.startNode.col]
	JPSPlusExplore_AllDirections(jps.startNode, jumpNode, jps)
	for 0 != jps.fastStack.Len() && jps.closeSet.len() < CloseSetMax {
		cur := jps.fastStack.PopNode()
		jps.openSet.remove(cur.row, cur.col)
		if cur.equal(jps.goalNode) {
			return true
		} else {
			jps.closeSet.insert(cur)
			jump := this[cur.row][cur.col]
			exploreDirections[(jump.blocked*8)+cur.direction](cur, jump, jps)
		}
	}
	return false
}

func (j *JPSPlus) FinalizePath() map[int]*LocJPS {
	jump := make(map[int]*LocJPS)
	for i, cur := 0, j.goalNode; nil != cur; cur, i = cur.parent, i+1 {
		jump[i] = newLocJPS(cur.row, cur.col)
	}
	return jump
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
