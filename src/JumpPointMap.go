package jpsplus

const (
	MovingDown = 1 << iota
	MovingRight
	MovingUp
	MovingLeft
)

type JumpPoint [MapHeight][MapWidth]int

func (j *JumpPoint) get(r int, c int) int {
	return j[r][c]
}

func (j *JumpPoint) CalculateJumpPointMap(b *BoolMap) {
	for r, data := range *b {
		for c, bl := range data {
			if bl {
				if b.IsJumpPoint(r, c, 1, 0) {
					j[r][c] |= MovingDown
				}
				if b.IsJumpPoint(r, c, -1, 0) {
					j[r][c] |= MovingUp
				}
				if b.IsJumpPoint(r, c, 0, 1) {
					j[r][c] |= MovingRight
				}
				if b.IsJumpPoint(r, c, 0, -1) {
					j[r][c] |= MovingLeft
				}
			}
		}
	}
}
