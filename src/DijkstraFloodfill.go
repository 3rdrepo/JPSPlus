package jpsplus

import (
// "fmt"
)

const (
	FIXED_POINT_MULTIPLIER = 100000
	FIXED_POINT_ONE
	FIXED_POINT_SQRT_2 = 141421

	numberOfBuckets   = 300000
	nodesInEachBucket = 1000
	division          = 10000
)

type DijkstraFloodfill struct {
	// m_width  int
	// m_height int
	// m_map              *BoolMap
	m_currentIteration int
	m_fastOpenList     *BucketPriorityQueue
	m_mapNodes         *DijkstraPathfindingNodeMap
	// m_openList         PQueue
}

var DefaultDijkstra = new(DijkstraFloodfill)

func fixed_point_shift(x int) int64 {
	return int64(x) * int64(FIXED_POINT_MULTIPLIER)
}

func (d *DijkstraFloodfill) init() {
	width := DefaultDistantJumpPoint.width()
	height := DefaultDistantJumpPoint.height()
	d.m_fastOpenList = newBucketPriorityQueue(numberOfBuckets, nodesInEachBucket, division)
	d.m_mapNodes = newMapNode(width, height)
	// Initialize nodes

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			node := &DijkstraPathfindingNode{
				m_row:                      r,
				m_col:                      c,
				m_iteration:                0,
				m_listStatus:               PathfindingNode_OnNone,
				m_blockedDirectionBitfield: 0,
			}
			var i uint
			for i = 0; i < 8; i++ {
				// Detect invalid movement from jump distances
				// (jump distance of zero is invalid movement)

				if DefaultDistantJumpPoint.getJumpdistance(r, c, int(i)) == 0 {
					node.m_blockedDirectionBitfield |= (1 << i)
				}
			}
			d.m_mapNodes.insert(r, c, node)
		}
	}
}

func (d *DijkstraFloodfill) Flood(r int, c int) {
	d.m_currentIteration += 1
	d.m_fastOpenList.Reset()

	if d.IsEmpty(r, c) {
		// Begin with starting node
		node := d.m_mapNodes.node(r, c)
		node.m_parent = nil
		node.m_directionFromStart = 255
		node.m_directionFromParent = 255
		node.m_givenCost = 0
		node.m_iteration = d.m_currentIteration
		node.m_listStatus = PathfindingNode_OnOpen
		// Explore outward in all directions on the starting node
		d.Explore_AllDirectionsWithChecks(node)

		node.m_listStatus = PathfindingNode_OnClosed
	}

	for !d.m_fastOpenList.Empty() {
		// fmt.Printf("m_fastOpenList %#v\n", d.m_fastOpenList.m_numNodesTracked)
		currentNode := d.m_fastOpenList.Pop()
		// fmt.Printf("fun %#v\n", exploreDirectionsDijkstraFlood[(currentNode.m_blockedDirectionBitfield*8)+currentNode.m_directionFromParent])
		exploreDirectionsDijkstraFlood[(currentNode.m_blockedDirectionBitfield*8)+currentNode.m_directionFromParent](currentNode)

		currentNode.m_listStatus = PathfindingNode_OnClosed
	}
}

func (d DijkstraFloodfill) GetCurrentInteration() int {
	return d.m_currentIteration
}
func (d *DijkstraFloodfill) Explore_AllDirectionsWithChecks(currentNode *DijkstraPathfindingNode) {
	//DOWN, DOWNRIGHT, RIGHT, UPRIGHT, UP, UPLEFT, LEFT, DOWNLEFT
	offsetRow := []int{1, 1, 0, -1, -1, -1, 0, 1}
	offsetCol := []int{0, 1, 1, 1, 0, -1, -1, -1}
	width := d.m_mapNodes.width()
	height := d.m_mapNodes.height()
	for i := 0; i < 8; i++ {
		neighborRow := currentNode.m_row + offsetRow[i]

		// Out of grid bounds?
		if (neighborRow >= 0) && (neighborRow < height) {
			neighborCol := currentNode.m_col + offsetCol[i]

			// Out of grid bounds?
			if (neighborCol >= 0) && (neighborCol < width) {
				// Valid tile - get the node
				// fmt.Printf("neighborRow = %v, neighborCol = %v\n", neighborRow, neighborCol)
				newSuccessor := d.m_mapNodes.node(neighborRow, neighborCol)
				// Blocked?
				if IsEmpty(neighborRow, neighborCol) {
					// Diagonal blocked?
					isDiagonal := (i & 0x1) == 1
					if isDiagonal && (!IsEmpty(currentNode.m_row+offsetRow[i], currentNode.m_col) ||
						!IsEmpty(currentNode.m_row, currentNode.m_col+offsetCol[i])) {
						continue
					}
					var costToNextNode int64
					if isDiagonal {
						costToNextNode = FIXED_POINT_SQRT_2
					} else {
						costToNextNode = FIXED_POINT_ONE
					}

					d.PushNewNode(newSuccessor, currentNode, i, i, costToNextNode)
				}
			}
		}

	}
}

func (d *DijkstraFloodfill) PushNewNode(newSuccessor *DijkstraPathfindingNode, currentNode *DijkstraPathfindingNode, startDirection int, parentDirection int, givenCost int64) {
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
	return IsEmpty(r, c)
}

func Explore_Null(currentNode *DijkstraPathfindingNode) {
	// Purposely does nothing
}

func Explore_D(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
}

func Explore_DR(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDownRight(currentNode)
}

func Explore_R(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
}

func Explore_UR(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUpRight(currentNode)
}

func Explore_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUp(currentNode)
}

func Explore_UL(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUpLeft(currentNode)
}

func Explore_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchLeft(currentNode)
}

func Explore_DL(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDownLeft(currentNode)
}

// Adjacent Doubles

func Explore_D_DR(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchDownRight(currentNode)
}

func Explore_DR_R(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDownRight(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
}

func Explore_R_UR(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchUpRight(currentNode)
}

func Explore_UR_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUpRight(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
}

func Explore_U_UL(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchUpLeft(currentNode)
}

func Explore_UL_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUpLeft(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
}

func Explore_L_DL(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchDownLeft(currentNode)
}

func Explore_DL_D(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDownLeft(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
}

// Non-Adjacent Cardinal Doubles

func Explore_D_R(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
}

func Explore_R_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
}

func Explore_U_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
}

func Explore_L_D(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
}

func Explore_D_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
}

func Explore_R_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
}

// Adjacent Triples

func Explore_D_DR_R(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchDownRight(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
}

func Explore_DR_R_UR(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDownRight(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchUpRight(currentNode)
}

func Explore_R_UR_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchUpRight(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
}

func Explore_UR_U_UL(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUpRight(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchUpLeft(currentNode)
}

func Explore_U_UL_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchUpLeft(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
}

func Explore_UL_L_DL(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUpLeft(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchDownLeft(currentNode)
}

func Explore_L_DL_D(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchDownLeft(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
}

func Explore_DL_D_DR(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDownLeft(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchDownRight(currentNode)
}

// Non-Adjacent Cardinal Triples

func Explore_D_R_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
}

func Explore_R_U_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
}

func Explore_U_L_D(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
}

func Explore_L_D_R(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
}

// Quads

func Explore_R_DR_D_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchDownRight(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
}

func Explore_R_D_DL_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchDownLeft(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
}

func Explore_U_UR_R_D(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchUpRight(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
}

func Explore_U_R_DR_D(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchDownRight(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
}

func Explore_L_UL_U_R(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchUpLeft(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
}

func Explore_L_U_UR_R(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchUpRight(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
}

func Explore_D_DL_L_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchDownLeft(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
}

func Explore_D_L_UL_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchUpLeft(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
}

// Quints

func Explore_R_DR_D_DL_L(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchDownRight(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchDownLeft(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
}

func Explore_U_UR_R_DR_D(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchUpRight(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchDownRight(currentNode)
	DefaultDijkstra.SearchDown(currentNode)
}

func Explore_L_UL_U_UR_R(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchUpLeft(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchUpRight(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
}

func Explore_D_DL_L_UL_U(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchDownLeft(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchUpLeft(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
}

func Explore_AllDirections(currentNode *DijkstraPathfindingNode) {
	DefaultDijkstra.SearchDown(currentNode)
	DefaultDijkstra.SearchDownLeft(currentNode)
	DefaultDijkstra.SearchLeft(currentNode)
	DefaultDijkstra.SearchUpLeft(currentNode)
	DefaultDijkstra.SearchUp(currentNode)
	DefaultDijkstra.SearchUpRight(currentNode)
	DefaultDijkstra.SearchRight(currentNode)
	DefaultDijkstra.SearchDownRight(currentNode)
}

func (d *DijkstraFloodfill) SearchDown(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row + 1
	newCol := currentNode.m_col
	givenCost := currentNode.m_givenCost + FIXED_POINT_ONE
	newSuccessor := d.m_mapNodes.node(newRow, newCol)
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, Down, givenCost)
}

func (d *DijkstraFloodfill) SearchDownRight(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row + 1
	newCol := currentNode.m_col + 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_SQRT_2
	newSuccessor := d.m_mapNodes.node(newRow, newCol)
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, DownRight, givenCost)
}

func (d *DijkstraFloodfill) SearchRight(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row
	newCol := currentNode.m_col + 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_ONE
	newSuccessor := d.m_mapNodes.node(newRow, newCol)
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, Right, givenCost)
}

func (d *DijkstraFloodfill) SearchUpRight(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row - 1
	newCol := currentNode.m_col + 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_SQRT_2
	newSuccessor := d.m_mapNodes.node(newRow, newCol)
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, UpRight, givenCost)
}

func (d *DijkstraFloodfill) SearchUp(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row - 1
	newCol := currentNode.m_col
	givenCost := currentNode.m_givenCost + FIXED_POINT_ONE
	newSuccessor := d.m_mapNodes.node(newRow, newCol)
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, Up, givenCost)
}

func (d *DijkstraFloodfill) SearchUpLeft(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row - 1
	newCol := currentNode.m_col - 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_SQRT_2
	newSuccessor := d.m_mapNodes.node(newRow, newCol)
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, UpLeft, givenCost)
}

func (d *DijkstraFloodfill) SearchLeft(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row
	newCol := currentNode.m_col - 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_ONE
	newSuccessor := d.m_mapNodes.node(newRow, newCol)
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, Left, givenCost)
}

func (d *DijkstraFloodfill) SearchDownLeft(currentNode *DijkstraPathfindingNode) {
	newRow := currentNode.m_row + 1
	newCol := currentNode.m_col - 1
	givenCost := currentNode.m_givenCost + FIXED_POINT_SQRT_2
	newSuccessor := d.m_mapNodes.node(newRow, newCol)
	d.PushNewNode(newSuccessor, currentNode, currentNode.m_directionFromStart, DownLeft, givenCost)
}
