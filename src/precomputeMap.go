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
	m_jumpPointMap                  [][]int
	m_distantJumpPointMap           []*DistantJumpPoints
	m_goalBoundsMap                 []*GoalBounds
	m_jumpDistancesAndGoalBoundsMap []*JumpDistancesAndGoalBounds
}

func newPrecomputeMap(width int, height int, mAp []bool) *PrecomputeMap {
	p := new(PrecomputeMap)
	p.m_mapCreated = false
	p.m_width = width
	p.m_height = height
	p.m_map = mAp
	return p
}

func (p *PrecomputeMap) CalculateMap() []*DistantJumpPoints {
	p.m_mapCreated = true
	p.CalculateJumpPointMap()
	p.CalculateDistantJumpPointMap()
	p.CalculateGoalBounding()
	return p.m_distantJumpPointMap
}

func (p *PrecomputeMap) CalculateJumpPointMap() {
	p.m_jumpPointMap = make([p.m_height][p.m_width]int)
	for r := 0; r < p.m_height; r++ {
		for c := 0; c < p.m_width; c++ {
			if m_map[c+(r*m_width)] {
				if p.IsJumpPoint(r, c, 1, 0) {
					m_jumpPointMap[r][c] |= MovingDown
				}
				if p.IsJumpPoint(r, c, -1, 0) {
					m_jumpPointMap[r][c] |= MovingUp
				}
				if p.IsJumpPoint(r, c, 0, 1) {
					m_jumpPointMap[r][c] |= MovingRight
				}
				if p.IsJumpPoint(r, c, 0, -1) {
					m_jumpPointMap[r][c] |= MovingLeft
				}
			}
		}
	}
}

func (p *PrecomputeMap) IsJumpPoint(r int, c int, rowDir int, colDir int) bool {
	return p.IsEmpty(r-rowDir, c-colDir) && // Parent not a wall (not necessary)
		((p.IsEmpty(r+colDir, c+rowDir) && // 1st forced neighbor
			p.IsWall(r-rowDir+colDir, c-colDir+rowDir)) || // 1st forced neighbor (continued)
			(p.IsEmpty(r-colDir, c-rowDir) && // 2nd forced neighbor
				p.IsWall(r-rowDir-colDir, c-colDir-rowDir))) // 2nd forced neighbor (continued)
}

func (p *PrecomputeMap) IsEmpty(r int, c int) bool {
	var colBoundsCheck uint = c
	var rowBoundsCheck uint = r
	if colBoundsCheck < uint.(p.m_width) && rowBoundsCheck < uint(m_height) {
		return p.m_map[c+(r*p.m_width)]
	} else {
		return false
	}
}

func (p *PrecomputeMap) IsWall(r int, c int) bool {
	var colBoundsCheck uint = c
	var rowBoundsCheck uint = r
	if colBoundsCheck < uint.(p.m_width) && rowBoundsCheck < uint(m_height) {
		return !p.m_map[c+(r*p.m_width)]
	} else {
		return true
	}
}
