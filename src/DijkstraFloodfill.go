package jpsplus

const (
	FIXED_POINT_MULTIPLIER = 100000
	FIXED_POINT_ONE
	FIXED_POINT_SQRT_2 = 141421

	numberOfBuckets   = 300000
	nodesInEachBucket = 1000
	division          = 10000
)

type DijkstraFloodfill struct {
	m_width            int
	m_height           int
	m_map              []bool
	m_currentIteration int
	m_fastOpenList     *BucketPriorityQueue
	m_mapNodes         [][]DijkstraPathfindingNode
	// m_openList         PQueue
}

var Dijkstra *DijkstraFloodfill

func fixed_point_shift(x int) int {
	return x * FIXED_POINT_MULTIPLIER
}

func newDijkstraFloodfill(width int, height int, mAp []bool, distantJumpPointMap [][]DistantJumpPoints) *DijkstraFloodfill {
	d := new(DijkstraFloodfill)
	d.m_width = width
	d.m_height = height
	d.m_map = mAp
	d.m_fastOpenList = newBucketPriorityQueue(numberOfBuckets, nodesInEachBucket, division)
	d.m_mapNodes = make([]([]DijkstraPathfindingNode), height)
	// Initialize nodes
	for pos := 0; pos < height; pos++ {
		d.m_mapNodes[pos] = make([]DijkstraPathfindingNode, width)
	}
	for r := 0; r < d.m_height; r++ {
		for c := 0; c < m_width; c++ {
			m_mapNodes[r][c] = DijkstraPathfindingNode{
				m_row:                      r,
				m_col:                      c,
				m_iteration:                0,
				m_listStatus:               PathfindingNode_OnNone,
				m_blockedDirectionBitfield: 0,
			}

			for i := 0; i < 8; i++ {
				// Detect invalid movement from jump distances
				// (jump distance of zero is invalid movement)
				if distantJumpPointMap[r][c].jumpDistance[i] == 0 {
					node.m_blockedDirectionBitfield |= (1 << i)
				}
			}
		}
	}
	Dijkstra = d
	return d
}

func (d *DijkstraFloodfill) Flood(r int, c int) {
	d.m_currentIteration += 1
	d.m_fastOpenList.Reset()

	if d.IsEmpty(r, c) {
		// Begin with starting node
		node := u.m_mapNodes[r][c]
		node.m_parent = NULL
		node.m_directionFromStart = 255
		node.m_directionFromParent = 255
		node.m_givenCost = 0
		node.m_iteration = m_currentIteration
		node.m_listStatus = PathfindingNode_OnOpen
		// Explore outward in all directions on the starting node
		Explore_AllDirectionsWithChecks(node)

		node.m_listStatus = PathfindingNode_OnClosed
		u.m_mapNodes[r][c] = node
	}

	for !u.m_fastOpenList.Empty() {
		currentNode := u.m_fastOpenList.Pop()
		exploreDirectionsDijkstraFlood[(currentNode.m_blockedDirectionBitfield*8)+currentNode.m_directionFromParent](currentNode)

		currentNode.m_listStatus = PathfindingNode_OnClosed
	}
}

func (d DijkstraFloodfill) GetCurrentInteration() int {
	return d.m_currentIteration
}
func (d *DijkstraFloodfill) Explore_AllDirectionsWithChecks(DijkstraPathfindingNode *currentNode) {
	//DOWN, DOWNRIGHT, RIGHT, UPRIGHT, UP, UPLEFT, LEFT, DOWNLEFT
	offsetRow := []int{1, 1, 0, -1, -1, -1, 0, 1}
	offsetCol := []int{0, 1, 1, 1, 0, -1, -1, -1}

	for i := 0; i < 8; i++ {
		neighborRow := currentNode.m_row + offsetRow[i]

		// Out of grid bounds?
		if neighborRow >= d.m_height {
			continue
		}

		neighborCol := currentNode.m_col + offsetCol[i]

		// Out of grid bounds?
		if neighborCol >= d.m_width {
			continue
		}

		// Valid tile - get the node
		newSuccessor := &(m_mapNodes[neighborRow][neighborCol])

		// Blocked?
		if !d.m_map[neighborCol+(neighborRow*d.m_width)] {
			continue
		}

		// Diagonal blocked?
		isDiagonal := (i & 0x1) == 1
		if isDiagonal && (!d.m_map[currentNode.m_col+((currentNode.m_row+offsetRow[i])*d.m_width)] ||
			!d.m_map[currentNode.m_col+offsetCol[i]+(currentNode.m_row*d.m_width)]) {
			continue
		}
		var costToNextNode uint
		if isDiagonal {
			costToNextNode = FIXED_POINT_SQRT_2
		} else {
			costToNextNode = FIXED_POINT_ONE
		}

		PushNewNode(newSuccessor, currentNode, i, i, costToNextNode)
	}
}

func (d *DijkstraFloodfill) PushNewNode(newSuccessor *DijkstraPathfindingNode, currentNode *DijkstraPathfindingNode, startDirection int, parentDirection int, givenCost uint) {
	if newSuccessor.m_iteration != d.m_currentIteration {
		// Place node on the Open list (we've never seen it before)
		newSuccessor.m_parent = currentNode
		newSuccessor.m_directionFromStart = startDirection
		newSuccessor.m_directionFromParent = parentDirection
		newSuccessor.m_givenCost = givenCost
		newSuccessor.m_listStatus = PathfindingNode_OnOpen
		newSuccessor.m_iteration = d.m_currentIteration

		d.m_fastOpenList.Push(newSuccessor)

	} else if (givenCost < newSuccessor.m_givenCost) && (PathfindingNode_OnOpen == newSuccessor.m_listStatus) {
		// We found a cheaper way to this node - update it
		lastCost := newSuccessor.m_givenCost
		newSuccessor.m_parent = currentNode
		newSuccessor.m_directionFromStart = startDirection
		newSuccessor.m_directionFromParent = parentDirection
		newSuccessor.m_givenCost = givenCost

		d.m_fastOpenList.DecreaseKey(newSuccessor, lastCost)
	}
}

func (d DijkstraFloodfill) IsEmpty(r int, c int) bool {
	if r < 0 {
		return false
	} else {
		if c < 0 {
			return false
		} else {
			if c < d.m_width && r < d.m_height {
				return d.m_map[c+(r*d.m_width)]
			} else {
				return false
			}
		}
	}
}

func Explore_Null(currentNode *DijkstraPathfindingNode) {
	// Purposely does nothing
}

func Explore_D(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
}

func Explore_DR(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDownRight(currentNode)
}

func Explore_R(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
}

func Explore_UR(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUpRight(currentNode)
}

func Explore_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUp(currentNode)
}

func Explore_UL(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUpLeft(currentNode)
}

func Explore_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchLeft(currentNode)
}

func Explore_DL(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDownLeft(currentNode)
}

// Adjacent Doubles

func Explore_D_DR(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchDownRight(currentNode)
}

func Explore_DR_R(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDownRight(currentNode)
	Dijkstra.SearchRight(currentNode)
}

func Explore_R_UR(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchUpRight(currentNode)
}

func Explore_UR_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUpRight(currentNode)
	Dijkstra.SearchUp(currentNode)
}

func Explore_U_UL(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchUpLeft(currentNode)
}

func Explore_UL_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUpLeft(currentNode)
	Dijkstra.SearchLeft(currentNode)
}

func Explore_L_DL(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchDownLeft(currentNode)
}

func Explore_DL_D(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDownLeft(currentNode)
	Dijkstra.SearchDown(currentNode)
}

// Non-Adjacent Cardinal Doubles

func Explore_D_R(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchRight(currentNode)
}

func Explore_R_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchUp(currentNode)
}

func Explore_U_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchLeft(currentNode)
}

func Explore_L_D(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchDown(currentNode)
}

func Explore_D_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchUp(currentNode)
}

func Explore_R_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchLeft(currentNode)
}

// Adjacent Triples

func Explore_D_DR_R(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchDownRight(currentNode)
	Dijkstra.SearchRight(currentNode)
}

func Explore_DR_R_UR(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDownRight(currentNode)
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchUpRight(currentNode)
}

func Explore_R_UR_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchUpRight(currentNode)
	Dijkstra.SearchUp(currentNode)
}

func Explore_UR_U_UL(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUpRight(currentNode)
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchUpLeft(currentNode)
}

func Explore_U_UL_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchUpLeft(currentNode)
	Dijkstra.SearchLeft(currentNode)
}

func Explore_UL_L_DL(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUpLeft(currentNode)
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchDownLeft(currentNode)
}

func Explore_L_DL_D(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchDownLeft(currentNode)
	Dijkstra.SearchDown(currentNode)
}

func Explore_DL_D_DR(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDownLeft(currentNode)
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchDownRight(currentNode)
}

// Non-Adjacent Cardinal Triples

func Explore_D_R_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchUp(currentNode)
}

func Explore_R_U_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchLeft(currentNode)
}

func Explore_U_L_D(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchDown(currentNode)
}

func Explore_L_D_R(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchRight(currentNode)
}

// Quads

func Explore_R_DR_D_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchDownRight(currentNode)
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchLeft(currentNode)
}

func Explore_R_D_DL_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchDownLeft(currentNode)
	Dijkstra.SearchLeft(currentNode)
}

func Explore_U_UR_R_D(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchUpRight(currentNode)
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchDown(currentNode)
}

func Explore_U_R_DR_D(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchDownRight(currentNode)
	Dijkstra.SearchDown(currentNode)
}

func Explore_L_UL_U_R(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchUpLeft(currentNode)
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchRight(currentNode)
}

func Explore_L_U_UR_R(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchUpRight(currentNode)
	Dijkstra.SearchRight(currentNode)
}

func Explore_D_DL_L_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchDownLeft(currentNode)
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchUp(currentNode)
}

func Explore_D_L_UL_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchUpLeft(currentNode)
	Dijkstra.SearchUp(currentNode)
}

// Quints

func Explore_R_DR_D_DL_L(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchDownRight(currentNode)
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchDownLeft(currentNode)
	Dijkstra.SearchLeft(currentNode)
}

func Explore_U_UR_R_DR_D(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchUpRight(currentNode)
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchDownRight(currentNode)
	Dijkstra.SearchDown(currentNode)
}

func Explore_L_UL_U_UR_R(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchUpLeft(currentNode)
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchUpRight(currentNode)
	Dijkstra.SearchRight(currentNode)
}

func Explore_D_DL_L_UL_U(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchDownLeft(currentNode)
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchUpLeft(currentNode)
	Dijkstra.SearchUp(currentNode)
}

func Explore_AllDirections(currentNode *DijkstraPathfindingNode) {
	Dijkstra.SearchDown(currentNode)
	Dijkstra.SearchDownLeft(currentNode)
	Dijkstra.SearchLeft(currentNode)
	Dijkstra.SearchUpLeft(currentNode)
	Dijkstra.SearchUp(currentNode)
	Dijkstra.SearchUpRight(currentNode)
	Dijkstra.SearchRight(currentNode)
	Dijkstra.SearchDownRight(currentNode)
}

func (d *DijkstraFloodfill) SearchDown(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row + 1
	newCol := currentNode.m_col
	givenCost := currentNode.m_givenCost + FIXED_POINT_ONE
	newSuccessor = &m_mapNodes[newRow][newCol]
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, Down, givenCost)
}

func (d *DijkstraFloodfill) SearchDownRight(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row + 1
	newCol := currentNode.m_col + 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_SQRT_2
	newSuccessor := &m_mapNodes[newRow][newCol]
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, DownRight, givenCost)
}

func (d *DijkstraFloodfill) SearchRight(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row
	newCol := currentNode.m_col + 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_ONE
	newSuccessor = &m_mapNodes[newRow][newCol]
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, Right, givenCost)
}

func (d *DijkstraFloodfill) SearchUpRight(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row - 1
	newCol := currentNode.m_col + 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_SQRT_2
	newSuccessor = &m_mapNodes[newRow][newCol]
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, UpRight, givenCost)
}

func (d *DijkstraFloodfill) SearchUp(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row - 1
	newCol := currentNode.m_col
	givenCost := currentNode.m_givenCost + FIXED_POINT_ONE
	newSuccessor := &m_mapNodes[newRow][newCol]
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, Up, givenCost)
}

func (d *DijkstraFloodfill) SearchUpLeft(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row - 1
	newCol := currentNode.m_col - 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_SQRT_2
	newSuccessor := &m_mapNodes[newRow][newCol]
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, UpLeft, givenCost)
}

func (d *DijkstraFloodfill) SearchLeft(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row
	newCol := currentNode.m_col - 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_ONE
	newSuccessor := &m_mapNodes[newRow][newCol]
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, Left, givenCost)
}

func (d *DijkstraFloodfill) SearchDownLeft(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row + 1
	newCol := currentNode.m_col - 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_SQRT_2
	DijkstraPathfindingNode * newSuccessor = &m_mapNodes[newRow][newCol]
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, DownLeft, givenCost)
}
