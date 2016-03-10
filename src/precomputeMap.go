package jpsplus

const (
	Down = iota
	DownRight
	Right
	UpRight
	Up
	UpLeft
	Left
	DownLeft
	All
)

const (
	MinRow = iota
	MaxRow
	MinCol
	MaxCol
)

const (
	MovingDown = 1 << iota
	MovingRight
	MovingUp
	MovingLeft
)

type DistantJumpPoints struct {
	jumpDistance [8]int
}

type GoalBounds struct {
	bounds [8][4]int
}

type JumpDistancesAndGoalBounds struct {
	blockedDirectionBitfield int
	jumpDistance             [8]int
	bounds                   [8][4]int
}

type PrecomputeMap struct {
	m_mapCreated                    bool
	m_width                         int
	m_height                        int
	m_map                           []bool
	m_jumpPointMap                  [][]int8
	m_distantJumpPointMap           [][]DistantJumpPoints
	m_goalBoundsMap                 [][]GoalBounds
	m_jumpDistancesAndGoalBoundsMap []*JumpDistancesAndGoalBounds
}

func newPrecomputeMap(width int, height int, mAp []bool) *PrecomputeMap {
	p := new(PrecomputeMap)
	p.m_mapCreated = false
	p.m_width = width
	p.m_height = height
	p.m_map = mAp
	p.m_jumpPointMap = make([]([]int8), height)
	p.m_distantJumpPointMap = make([]([]DistantJumpPoints), height)
	p.m_goalBoundsMap = make([]([]GoalBounds), height)
	for pos := 0; pos < height; pos++ {
		p.m_jumpPointMap[pos] = make([]int8, width)
		p.m_distantJumpPointMap[pos] = make([]DistantJumpPoints, width)
		p.m_goalBoundsMap[pos] = make([]GoalBounds, width)
	}

	return p
}

func (p *PrecomputeMap) CalculateMap() [][]DistantJumpPoints {
	p.m_mapCreated = true
	p.CalculateJumpPointMap()
	p.CalculateDistantJumpPointMap()
	p.CalculateGoalBounding()
	return p.m_distantJumpPointMap
}

func (p *PrecomputeMap) CalculateJumpPointMap() {
	for r := 0; r < p.m_height; r++ {
		for c := 0; c < p.m_width; c++ {
			if p.m_map[c+(r*p.m_width)] {
				if p.IsJumpPoint(r, c, 1, 0) {
					p.m_jumpPointMap[r][c] |= MovingDown
				}
				if p.IsJumpPoint(r, c, -1, 0) {
					p.m_jumpPointMap[r][c] |= MovingUp
				}
				if p.IsJumpPoint(r, c, 0, 1) {
					p.m_jumpPointMap[r][c] |= MovingRight
				}
				if p.IsJumpPoint(r, c, 0, -1) {
					p.m_jumpPointMap[r][c] |= MovingLeft
				}
			}
		}
	}
}

func (p *PrecomputeMap) IsJumpPoint(r int, c int, rowDir int, colDir int) bool {
	// return p.IsEmpty(r-rowDir, c-colDir) && // Parent not a wall (not necessary)
	return ((p.IsEmpty(r+colDir, c+rowDir) && // 1st forced neighbor
		p.IsWall(r-rowDir+colDir, c-colDir+rowDir)) || // 1st forced neighbor (continued)
		(p.IsEmpty(r-colDir, c-rowDir) && // 2nd forced neighbor
			p.IsWall(r-rowDir-colDir, c-colDir-rowDir))) // 2nd forced neighbor (continued)
}

func (p *PrecomputeMap) IsEmpty(r int, c int) bool {
	if 0 > r {
		return false
	} else {
		if 0 > c {
			return false
		} else {
			if (c < p.m_width) && (r < p.m_height) {
				return p.m_map[c+(r*p.m_width)]
			} else {
				return false
			}
		}
	}
}

func (p *PrecomputeMap) IsWall(r int, c int) bool {
	return !p.IsEmpty(r, c)
}

func (p *PrecomputeMap) CalculateDistantJumpPointMap() {
	// Calculate distant jump points (Left and Right)
	for r := 0; r < p.m_height; r++ {
		{
			countMovingLeft := -1
			jumpPointLastSeen := false
			for c := 0; c < p.m_width; c++ {
				if p.IsWall(r, c) {
					countMovingLeft = -1
					jumpPointLastSeen = false
					p.m_distantJumpPointMap[r][c].jumpDistance[Left] = 0
				} else {
					countMovingLeft++
					if jumpPointLastSeen {
						p.m_distantJumpPointMap[r][c].jumpDistance[Left] = countMovingLeft
					} else {
						// Wall last seen
						p.m_distantJumpPointMap[r][c].jumpDistance[Left] = -countMovingLeft
					}
					if (p.m_jumpPointMap[r][c] & MovingLeft) > 0 {
						countMovingLeft = 0
						jumpPointLastSeen = true
					}
				}

			}
		}

		{
			countMovingRight := -1
			jumpPointLastSeen := false
			for c := p.m_width - 1; c >= 0; c-- {
				if p.IsWall(r, c) {
					countMovingRight = -1
					jumpPointLastSeen = false
					p.m_distantJumpPointMap[r][c].jumpDistance[Right] = 0
					continue
				} else {
					countMovingRight++
					if jumpPointLastSeen {
						p.m_distantJumpPointMap[r][c].jumpDistance[Right] = countMovingRight
					} else {
						// Wall last see
						p.m_distantJumpPointMap[r][c].jumpDistance[Right] = -countMovingRight
					}
					if (p.m_jumpPointMap[r][c] & MovingRight) > 0 {
						countMovingRight = 0
						jumpPointLastSeen = true
					}
				}

			}
		}
	}

	// Calculate distant jump points (Up and Down)
	for c := 0; c < p.m_width; c++ {
		{
			countMovingUp := -1
			jumpPointLastSeen := false
			for r := 0; r < p.m_height; r++ {
				if p.IsWall(r, c) {
					countMovingUp = -1
					jumpPointLastSeen = false
					p.m_distantJumpPointMap[r][c].jumpDistance[Up] = 0
					continue
				}
				countMovingUp++
				if jumpPointLastSeen {
					p.m_distantJumpPointMap[r][c].jumpDistance[Up] = countMovingUp
				} else { // Wall last seen
					p.m_distantJumpPointMap[r][c].jumpDistance[Up] = -countMovingUp
				}
				if (p.m_jumpPointMap[r][c] & MovingUp) > 0 {
					countMovingUp = 0
					jumpPointLastSeen = true
				}
			}
		}

		{
			countMovingDown := -1
			jumpPointLastSeen := false
			for r := p.m_height - 1; r >= 0; r-- {
				if p.IsWall(r, c) {
					countMovingDown = -1
					jumpPointLastSeen = false
					p.m_distantJumpPointMap[r][c].jumpDistance[Down] = 0
					continue
				}
				countMovingDown++
				if jumpPointLastSeen {
					p.m_distantJumpPointMap[r][c].jumpDistance[Down] = countMovingDown
				} else { // Wall last seen
					p.m_distantJumpPointMap[r][c].jumpDistance[Down] = -countMovingDown
				}
				if (p.m_jumpPointMap[r][c] & MovingDown) > 0 {
					countMovingDown = 0
					jumpPointLastSeen = true
				}
			}
		}
	}

	// Calculate distant jump points (Diagonally UpLeft and UpRight)
	for r := 0; r < p.m_height; r++ {
		for c := 0; c < p.m_width; c++ {
			if p.IsEmpty(r, c) {
				if r == 0 || c == 0 || (p.IsWall(r-1, c) || p.IsWall(r, c-1) || p.IsWall(r-1, c-1)) {
					// Wall one away
					p.m_distantJumpPointMap[r][c].jumpDistance[UpLeft] = 0
				} else if p.IsEmpty(r-1, c) && p.IsEmpty(r, c-1) &&
					(p.m_distantJumpPointMap[r-1][c-1].jumpDistance[Up] > 0 ||
						p.m_distantJumpPointMap[r-1][c-1].jumpDistance[Left] > 0) {
					// Diagonal one away
					p.m_distantJumpPointMap[r][c].jumpDistance[UpLeft] = 1
				} else {
					// Increment from last
					jumpDistance := p.m_distantJumpPointMap[r-1][c-1].jumpDistance[UpLeft]

					if jumpDistance > 0 {
						p.m_distantJumpPointMap[r][c].jumpDistance[UpLeft] = 1 + jumpDistance
					} else {
						//if( jumpDistance <= 0 )
						p.m_distantJumpPointMap[r][c].jumpDistance[UpLeft] = -1 + jumpDistance
					}
				}

				if r == 0 || c == p.m_width-1 || (p.IsWall(r-1, c) || p.IsWall(r, c+1) || p.IsWall(r-1, c+1)) {
					// Wall one away
					p.m_distantJumpPointMap[r][c].jumpDistance[UpRight] = 0
				} else if p.IsEmpty(r-1, c) && p.IsEmpty(r, c+1) &&
					(p.m_distantJumpPointMap[r-1][c+1].jumpDistance[Up] > 0 ||
						p.m_distantJumpPointMap[r-1][c+1].jumpDistance[Right] > 0) {
					// Diagonal one away
					p.m_distantJumpPointMap[r][c].jumpDistance[UpRight] = 1
				} else {
					// Increment from last
					jumpDistance := p.m_distantJumpPointMap[r-1][c+1].jumpDistance[UpRight]

					if jumpDistance > 0 {
						p.m_distantJumpPointMap[r][c].jumpDistance[UpRight] = 1 + jumpDistance
					} else {
						//if( jumpDistance <= 0 )
						p.m_distantJumpPointMap[r][c].jumpDistance[UpRight] = -1 + jumpDistance
					}
				}
			}
		}
	}

	// Calculate distant jump points (Diagonally DownLeft and DownRight)
	for r := p.m_height - 1; r >= 0; r-- {
		for c := 0; c < p.m_width; c++ {
			if p.IsEmpty(r, c) {
				if r == p.m_height-1 || c == 0 ||
					(p.IsWall(r+1, c) || p.IsWall(r, c-1) || p.IsWall(r+1, c-1)) {
					// Wall one away
					p.m_distantJumpPointMap[r][c].jumpDistance[DownLeft] = 0
				} else if p.IsEmpty(r+1, c) && p.IsEmpty(r, c-1) &&
					(p.m_distantJumpPointMap[r+1][c-1].jumpDistance[Down] > 0 ||
						p.m_distantJumpPointMap[r+1][c-1].jumpDistance[Left] > 0) {
					// Diagonal one away
					p.m_distantJumpPointMap[r][c].jumpDistance[DownLeft] = 1
				} else {
					// Increment from last
					jumpDistance := p.m_distantJumpPointMap[r+1][c-1].jumpDistance[DownLeft]

					if jumpDistance > 0 {
						p.m_distantJumpPointMap[r][c].jumpDistance[DownLeft] = 1 + jumpDistance
					} else { //if( jumpDistance <= 0 )
						p.m_distantJumpPointMap[r][c].jumpDistance[DownLeft] = -1 + jumpDistance
					}
				}

				if r == p.m_height-1 || c == p.m_width-1 || (p.IsWall(r+1, c) || p.IsWall(r, c+1) || p.IsWall(r+1, c+1)) {
					// Wall one away
					p.m_distantJumpPointMap[r][c].jumpDistance[DownRight] = 0
				} else if p.IsEmpty(r+1, c) && p.IsEmpty(r, c+1) &&
					(p.m_distantJumpPointMap[r+1][c+1].jumpDistance[Down] > 0 ||
						p.m_distantJumpPointMap[r+1][c+1].jumpDistance[Right] > 0) {
					// Diagonal one away
					p.m_distantJumpPointMap[r][c].jumpDistance[DownRight] = 1
				} else {
					// Increment from last
					jumpDistance := p.m_distantJumpPointMap[r+1][c+1].jumpDistance[DownRight]

					if jumpDistance > 0 {
						p.m_distantJumpPointMap[r][c].jumpDistance[DownRight] = 1 + jumpDistance
					} else { //if( jumpDistance <= 0 )
						p.m_distantJumpPointMap[r][c].jumpDistance[DownRight] = -1 + jumpDistance
					}
				}
			}
		}
	}
}

func (p *PrecomputeMap) CalculateGoalBounding() {
	fmt.Printf("Goal Bounding Preprocessing\n")

	dijkstra := newDijkstraFloodfill(p.m_width, p.m_height, p.m_map, p.m_distantJumpPointMap)

	// InitArray(m_goalBoundsMap, m_width, m_height);
	for r := 0; r < m_height; r++ {
		for c := 0; c < m_width; c++ {
			for dir := 0; dir < 8; dir++ {
				p.m_goalBoundsMap[r][c].bounds[dir][MinRow] = p.m_height
				p.m_goalBoundsMap[r][c].bounds[dir][MaxRow] = 0
				p.m_goalBoundsMap[r][c].bounds[dir][MinCol] = p.m_width
				p.m_goalBoundsMap[r][c].bounds[dir][MaxCol] = 0
			}
		}
	}

	for startRow := 0; startRow < m_height; startRow++ {
		fmt.Printf("Row: %d\n", startRow)

		for startCol := 0; startCol < m_width; startCol++ {
			if p.IsEmpty(startRow, startCol) {
				dijkstra.Flood(startRow, startCol)
				currentIteration := dijkstra.GetCurrentInteration()

				for r := 0; r < m_height; r++ {
					for c := 0; c < m_width; c++ {
						if p.IsEmpty(r, c) {
							iteration := dijkstra.m_mapNodes[r][c].m_iteration
							status := dijkstra.m_mapNodes[r][c].m_listStatus
							dir := dijkstra.m_mapNodes[r][c].m_directionFromStart

							if iteration == currentIteration &&
								status == OnClosed &&
								dir >= 0 && dir <= 7 {
								row := dijkstra.m_mapNodes[r][c].m_row
								col := dijkstra.m_mapNodes[r][c].m_col

								if p.m_goalBoundsMap[startRow][startCol].bounds[dir][MinRow] > row {
									p.m_goalBoundsMap[startRow][startCol].bounds[dir][MinRow] = row
								}
								if p.m_goalBoundsMap[startRow][startCol].bounds[dir][MaxRow] < row {
									p.m_goalBoundsMap[startRow][startCol].bounds[dir][MaxRow] = row
								}
								if p.m_goalBoundsMap[startRow][startCol].bounds[dir][MinCol] > col {
									p.m_goalBoundsMap[startRow][startCol].bounds[dir][MinCol] = col
								}
								if p.m_goalBoundsMap[startRow][startCol].bounds[dir][MaxCol] < col {
									p.m_goalBoundsMap[startRow][startCol].bounds[dir][MaxCol] = col
								}
							}
						}

					}
				}
			}

		}
	}

}
