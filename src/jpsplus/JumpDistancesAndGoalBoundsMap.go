package jpsplus

type JumpDistancesAndGoalBounds struct {
	blockedDirectionBitfield int
	jumpDistance             *DistantJumpPoints
	// bounds                   *GoalBounds
}
type JumpDistancesAndGoalBoundsMap [][]*JumpDistancesAndGoalBounds

var DefaultJumpDistancesAndGoalBounds *JumpDistancesAndGoalBoundsMap

func initDefaultJumpDistancesAndGoalBounds() {
	width := DefaultDistantJumpPoint.width()
	height := DefaultDistantJumpPoint.height()

	j := make(JumpDistancesAndGoalBoundsMap, height)
	for pos := 0; pos < height; pos++ {
		j[pos] = make([]*JumpDistancesAndGoalBounds, width)
	}
	DefaultJumpDistancesAndGoalBounds = &j
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			jmp := new(JumpDistancesAndGoalBounds)
			// jmp.blockedDirectionBitfield = DefaultDijkstra.m_mapNodes.nodeBlocked(r, c)
			jmp.blockedDirectionBitfield = DefaultDistantJumpPoint.getBlocked(r, c)
			// jmp.bounds = DefautGoalBounds.get(r, c)
			jmp.jumpDistance = DefaultDistantJumpPoint.get(r, c)

			DefaultJumpDistancesAndGoalBounds.set(r, c, jmp)
		}
	}
}

func (j JumpDistancesAndGoalBoundsMap) set(r int, c int, node *JumpDistancesAndGoalBounds) {
	j[r][c] = node
}

func (j JumpDistancesAndGoalBoundsMap) get(r int, c int) *JumpDistancesAndGoalBounds {
	return j[r][c]
}

func (j JumpDistancesAndGoalBoundsMap) getBlocked(r int, c int) int {
	return j[r][c].blockedDirectionBitfield
}

func (j JumpDistancesAndGoalBoundsMap) width() int {
	return len(j[0])
}

func (j JumpDistancesAndGoalBoundsMap) height() int {
	return len(j)
}
