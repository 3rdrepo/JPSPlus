package jpsplus

const (
	MovingDown = 1 << iota
	MovingRight
	MovingUp
	MovingLeft
)

type JumpPoint [][]int

func newJumpPoint(width int, height int) *JumpPoint {
	j := make(JumpPoint, height)
	for r := 0; r < height; r++ {
		j[r] = make([]int, width)
	}
	return &j
}

func (j JumpPoint) move(r int, c int, v int) {
	j[r][c] |= v
}

func (j JumpPoint) get(r int, c int) int {
	return j[r][c]
}

func (j *JumpPoint) CalculateJumpPointMap(b *BoolMap) {
	for r, data := range *b {
		for c, bl := range data {
			if bl {
				if b.IsJumpPoint(r, c, 1, 0) {
					j.move(r, c, MovingDown)
				}
				if b.IsJumpPoint(r, c, -1, 0) {
					j.move(r, c, MovingUp)
				}
				if b.IsJumpPoint(r, c, 0, 1) {
					j.move(r, c, MovingRight)
				}
				if b.IsJumpPoint(r, c, 0, -1) {
					j.move(r, c, MovingLeft)
				}
			}
		}
	}
}
