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

type Jump struct {
	distant []int
	blocked int
}

func newJump() *Jump {
	j := new(Jump)
	j.distant = make([]int, 8)
	return j
}

type JumpMap [][]*Jump

var DefaultJumpMap *JumpMap

func NewJumpMap(width int, height int) *JumpMap {
	j := make(JumpMap, height)
	for r := 0; r < height; r++ {
		j[r] = make([]*Jump, width)
		for c := 0; c < width; c++ {
			j[r][c] = newJump()
		}
	}
	return &j
}

func (j JumpMap) setJumpdistance(r int, c int, dir int, v int) {
	j[r][c].distant[dir] = v
}

func (j JumpMap) getJumpdistance(r int, c int, dir int) int {
	return j[r][c].distant[dir]
}

func (j JumpMap) setBlocked(r int, c int, b int) {
	j[r][c].blocked = b
}

func (j JumpMap) Jump(r int, c int) *Jump {
	return j[r][c]
}

func (j *JumpMap) CalculateDistantJumpPointMapLeft(b *BoolMap, jp *JumpPoint) {
	var countMovingLeft int
	var jumpPointLastSeen bool
	width := b.width()
	height := b.height()
	for r := 0; r < height; r++ {
		countMovingLeft = -1
		jumpPointLastSeen = false
		for c := 0; c < width; c++ {
			if !b.IsEmpty(r, c) {
				countMovingLeft = -1
				jumpPointLastSeen = false
				j.setJumpdistance(r, c, Left, 0)
			} else {
				countMovingLeft += 1
				if jumpPointLastSeen {
					j.setJumpdistance(r, c, Left, countMovingLeft)
				} else {
					j.setJumpdistance(r, c, Left, -countMovingLeft)
				}
				if (jp.get(r, c) & MovingLeft) > 0 {
					countMovingLeft = 0
					jumpPointLastSeen = true
				}
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMapRight(b *BoolMap, jp *JumpPoint) {
	var countMovingRight int
	var jumpPointLastSeen bool
	width := b.width()
	height := b.height()
	for r := 0; r < height; r++ {
		countMovingRight = -1
		jumpPointLastSeen = false
		for c := width - 1; c >= 0; c-- {
			if !b.IsEmpty(r, c) {
				countMovingRight = -1
				jumpPointLastSeen = false
				j.setJumpdistance(r, c, Right, 0)
			} else {
				countMovingRight += 1
				if jumpPointLastSeen {
					j.setJumpdistance(r, c, Right, countMovingRight)
				} else {
					j.setJumpdistance(r, c, Right, -countMovingRight)
				}
				if (jp.get(r, c) & MovingRight) > 0 {
					countMovingRight = 0
					jumpPointLastSeen = true
				}
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMapUp(b *BoolMap, jp *JumpPoint) {
	var countMovingUp int
	var jumpPointLastSeen bool
	width := b.width()
	height := b.height()
	for c := 0; c < width; c++ {
		countMovingUp = -1
		jumpPointLastSeen = false
		for r := 0; r < height; r++ {
			if !b.IsEmpty(r, c) {
				countMovingUp = -1
				jumpPointLastSeen = false
				// p.m_distantJumpPointMap[r][c].jumpDistance[Up] = 0
				j.setJumpdistance(r, c, Up, 0)
			} else {
				countMovingUp += 1
				if jumpPointLastSeen {
					j.setJumpdistance(r, c, Up, countMovingUp)
				} else {
					j.setJumpdistance(r, c, Up, -countMovingUp)
				}
				if (jp.get(r, c) & MovingUp) > 0 {
					countMovingUp = 0
					jumpPointLastSeen = true
				}
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMapDown(b *BoolMap, jp *JumpPoint) {
	var countMovingDown int
	var jumpPointLastSeen bool
	width := b.width()
	height := b.height()
	for c := 0; c < width; c++ {
		countMovingDown = -1
		jumpPointLastSeen = false
		for r := height - 1; r >= 0; r-- {
			if !b.IsEmpty(r, c) {
				countMovingDown = -1
				jumpPointLastSeen = false
				j.setJumpdistance(r, c, Down, 0)
			} else {
				countMovingDown += 1
				if jumpPointLastSeen {
					j.setJumpdistance(r, c, Down, countMovingDown)
				} else {
					j.setJumpdistance(r, c, Down, -countMovingDown)
				}
				if (jp.get(r, c) & MovingDown) > 0 {
					countMovingDown = 0
					jumpPointLastSeen = true
				}
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMapUpLeftandUpRight(b *BoolMap) {
	width := b.width()
	height := b.height()
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if b.IsEmpty(r, c) {
				//UpLeft
				if r == 0 || c == 0 || (!b.IsEmpty(r-1, c) || !b.IsEmpty(r, c-1) || !b.IsEmpty(r-1, c-1)) {
					j.setJumpdistance(r, c, UpLeft, 0)
				} else if b.IsEmpty(r-1, c) && b.IsEmpty(r, c-1) &&
					(j.getJumpdistance(r-1, c-1, Up) > 0 ||
						j.getJumpdistance(r-1, c-1, Left) > 0) {
					j.setJumpdistance(r, c, UpLeft, 1)
				} else {
					jumpDistance := j.getJumpdistance(r-1, c-1, UpLeft)

					if jumpDistance > 0 {
						j.setJumpdistance(r, c, UpLeft, jumpDistance+1)
					} else {
						j.setJumpdistance(r, c, UpLeft, jumpDistance-1)
					}
				}

				//UpRight
				if r == 0 || c == width-1 || (!b.IsEmpty(r-1, c) || !b.IsEmpty(r, c+1) || !b.IsEmpty(r-1, c+1)) {
					j.setJumpdistance(r, c, UpRight, 0)
				} else if b.IsEmpty(r-1, c) && b.IsEmpty(r, c+1) &&
					(j.getJumpdistance(r-1, c+1, Up) > 0 ||
						j.getJumpdistance(r-1, c+1, Right) > 0) {
					j.setJumpdistance(r, c, UpRight, 1)
				} else {
					jumpDistance := j.getJumpdistance(r-1, c+1, UpRight)
					if jumpDistance > 0 {
						j.setJumpdistance(r, c, UpRight, jumpDistance+1)
					} else {
						j.setJumpdistance(r, c, UpRight, jumpDistance-1)
					}
				}
			}
		}
	}
}
func (j *JumpMap) CalculateDistantJumpPointMapDownLeftandDownRight(b *BoolMap) {
	width := b.width()
	height := b.height()
	for r := height - 1; r >= 0; r-- {
		for c := 0; c < width; c++ {
			if b.IsEmpty(r, c) {
				if r == height-1 || c == 0 ||
					(!b.IsEmpty(r+1, c) || !b.IsEmpty(r, c-1) || !b.IsEmpty(r+1, c-1)) {
					j.setJumpdistance(r, c, DownLeft, 0)
				} else if b.IsEmpty(r+1, c) && b.IsEmpty(r, c-1) &&
					(j.getJumpdistance(r+1, c-1, Down) > 0 ||
						j.getJumpdistance(r+1, c-1, Left) > 0) {
					j.setJumpdistance(r, c, DownLeft, 1)
				} else {
					// jumpDistance := p.m_distantJumpPointMap[r+1][c-1].jumpDistance[DownLeft]
					jumpDistance := j.getJumpdistance(r+1, c-1, DownLeft)
					if jumpDistance > 0 {
						// p.m_distantJumpPointMap[r][c].jumpDistance[DownLeft] = 1 + jumpDistance
						j.setJumpdistance(r, c, DownLeft, jumpDistance+1)
					} else {
						j.setJumpdistance(r, c, DownLeft, jumpDistance-1)
					}
				}

				if r == height-1 || c == width-1 || (!b.IsEmpty(r+1, c) || !b.IsEmpty(r, c+1) || !b.IsEmpty(r+1, c+1)) {
					j.setJumpdistance(r, c, DownRight, 0)
				} else if b.IsEmpty(r+1, c) && b.IsEmpty(r, c+1) &&
					(j.getJumpdistance(r+1, c+1, Down) > 0 ||
						j.getJumpdistance(r+1, c+1, Right) > 0) {
					j.setJumpdistance(r, c, DownRight, 1)
				} else {
					jumpDistance := j.getJumpdistance(r+1, c+1, DownRight)

					if jumpDistance > 0 {
						j.setJumpdistance(r, c, DownRight, jumpDistance+1)
					} else {
						j.setJumpdistance(r, c, DownRight, jumpDistance-1)
					}
				}
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMap(b *BoolMap, jp *JumpPoint) {
	// Calculate distant jump points (Left and Right)
	j.CalculateDistantJumpPointMapLeft(b, jp)
	j.CalculateDistantJumpPointMapRight(b, jp)

	// Calculate distant jump points (Up and Down)
	j.CalculateDistantJumpPointMapUp(b, jp)
	j.CalculateDistantJumpPointMapDown(b, jp)

	// Calculate distant jump points (Diagonally UpLeft and UpRight)
	j.CalculateDistantJumpPointMapUpLeftandUpRight(b)
	// Calculate distant jump points (Diagonally DownLeft and DownRight)
	j.CalculateDistantJumpPointMapDownLeftandDownRight(b)
}

func (j *JumpMap) CalculateBlock() {
	for r, data := range *j {
		for c, jump := range data {
			var blocked int
			for i := 0; i < 8; i++ {
				if 0 == jump.distant[i] {
					blocked |= (1 << uint(i))
				}
			}
			j.setBlocked(r, c, blocked)
		}
	}
}
