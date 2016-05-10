package jpsplus

import (
// "fmt"
)

func PreprocessMap(m *BoolMap) *JumpMap {
	jumpPoint := new(JumpPoint)
	jumpPoint.CalculateJumpPointMap(m)
	jumpMap := new(JumpMap)
	jumpMap.CalculateDistantJumpPointMap(m, jumpPoint)
	jumpMap.CalculateBlock()
	return jumpMap
}

type BoolMap [MapHeight][MapWidth]bool

func (b *BoolMap) IsEmpty(r int, c int) bool {
	if r < 0 || c < 0 || r >= MapHeight || c >= MapWidth {
		return false
	} else {
		return b[r][c]
	}
}

func (b *BoolMap) IsJumpPoint(r int, c int, rowDir int, colDir int) bool {
	return (b.IsEmpty(r+colDir, c+rowDir) && !b.IsEmpty(r-rowDir+colDir, c-colDir+rowDir)) ||
		(b.IsEmpty(r-colDir, c-rowDir) && !b.IsEmpty(r-rowDir-colDir, c-colDir-rowDir))
}

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
	distant [8]int
	blocked int
}

type JumpMap [MapHeight][MapWidth]*Jump

func (j *JumpMap) CalculateDistantJumpPointMapLeft(b *BoolMap, jp *JumpPoint) {
	var countMovingLeft int
	var jumpPointLastSeen bool
	for r := 0; r < MapHeight; r++ {
		countMovingLeft = -1
		jumpPointLastSeen = false
		for c := 0; c < MapWidth; c++ {
			jump := new(Jump)
			j[r][c] = jump
			if b[r][c] {
				countMovingLeft += 1
				if jumpPointLastSeen {
					jump.distant[Left] = countMovingLeft
				} else {
					jump.distant[Left] = -countMovingLeft
				}
				if (jp[r][c] & MovingLeft) == MovingLeft {
					countMovingLeft = 0
					jumpPointLastSeen = true
				}
			} else {
				countMovingLeft = -1
				jumpPointLastSeen = false
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMapRight(b *BoolMap, jp *JumpPoint) {
	var countMovingRight int
	var jumpPointLastSeen bool
	for r := 0; r < MapHeight; r++ {
		countMovingRight = -1
		jumpPointLastSeen = false
		for c := MapWidth - 1; c >= 0; c-- {
			if b[r][c] {
				countMovingRight += 1
				if jumpPointLastSeen {
					j[r][c].distant[Right] = countMovingRight
				} else {
					j[r][c].distant[Right] = -countMovingRight
				}
				if (jp[r][c] & MovingRight) == MovingRight {
					countMovingRight = 0
					jumpPointLastSeen = true
				}
			} else {
				countMovingRight = -1
				jumpPointLastSeen = false
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMapUp(b *BoolMap, jp *JumpPoint) {
	var countMovingUp int
	var jumpPointLastSeen bool
	for c := 0; c < MapWidth; c++ {
		countMovingUp = -1
		jumpPointLastSeen = false
		for r := 0; r < MapHeight; r++ {
			if b[r][c] {
				countMovingUp += 1
				if jumpPointLastSeen {
					j[r][c].distant[Up] = countMovingUp
				} else {
					j[r][c].distant[Up] = -countMovingUp
				}
				if (jp[r][c] & MovingUp) == MovingUp {
					countMovingUp = 0
					jumpPointLastSeen = true
				}
			} else {
				countMovingUp = -1
				jumpPointLastSeen = false
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMapDown(b *BoolMap, jp *JumpPoint) {
	var countMovingDown int
	var jumpPointLastSeen bool
	for c := 0; c < MapWidth; c++ {
		countMovingDown = -1
		jumpPointLastSeen = false
		for r := MapHeight - 1; r >= 0; r-- {
			if b[r][c] {
				countMovingDown += 1
				if jumpPointLastSeen {
					j[r][c].distant[Down] = countMovingDown
				} else {
					j[r][c].distant[Down] = -countMovingDown
				}
				if (jp[r][c] & MovingDown) == MovingDown {
					countMovingDown = 0
					jumpPointLastSeen = true
				}
			} else {
				countMovingDown = -1
				jumpPointLastSeen = false
			}
		}
	}
}

func (j *JumpMap) CalculateDistantJumpPointMapUpLeftandUpRight(b *BoolMap) {
	for r := 0; r < MapHeight; r++ {
		for c := 0; c < MapWidth; c++ {
			if b[r][c] {
				//UpLeft
				prexR := r - 1
				prexC := c - 1
				if r > 0 && c > 0 && b[prexR][c] && b[r][c] && b[prexR][prexC] {
					if b[prexR][c] && b[r][prexC] &&
						(j[prexR][prexC].distant[Up] > 0 || j[prexR][prexC].distant[Left] > 0) {
						j[r][c].distant[UpLeft] = 1
					} else {
						jumpDistance := j[prexR][prexC].distant[UpLeft]
						if jumpDistance > 0 {
							j[prexR][prexC].distant[UpLeft] = jumpDistance + 1
						} else {
							j[prexR][prexC].distant[UpLeft] = jumpDistance - 1
						}
					}
				}
				//UpRight
				nextC := c + 1
				if r > 0 && c < MapWidth-1 && b[prexR][c] && b[r][nextC] && b[prexR][nextC] {
					if b[prexR][c] && b[r][nextC] &&
						(j[prexR][nextC].distant[Up] > 0 ||
							j[prexR][nextC].distant[Right] > 0) {
						j[r][c].distant[UpRight] = 1
					} else {
						jumpDistance := j[prexR][c+1].distant[UpRight]
						if jumpDistance > 0 {
							j[r][c].distant[UpRight] = jumpDistance + 1
						} else {
							j[r][c].distant[UpRight] = jumpDistance - 1
						}
					}
				}
			}
		}
	}
}
func (j *JumpMap) CalculateDistantJumpPointMapDownLeftandDownRight(b *BoolMap) {
	for r := MapHeight - 1; r >= 0; r-- {
		for c := 0; c < MapWidth; c++ {
			if b[r][c] {
				nextR := r + 1
				prexC := c - 1
				if r < MapHeight-1 && c > 0 && b[nextR][c] && b[r][prexC] && b[nextR][prexC] {
					if b[nextR][c] && b[r][prexC] &&
						(j[nextR][prexC].distant[Down] > 0 ||
							j[nextR][prexC].distant[Left] > 0) {
						j[r][c].distant[DownLeft] = 1
					} else {
						jumpDistance := j[nextR][prexC].distant[DownLeft]
						if jumpDistance > 0 {
							j[r][c].distant[DownLeft] = jumpDistance + 1
						} else {
							j[r][c].distant[DownLeft] = jumpDistance - 1
						}
					}
				}
				nextC := c + 1
				if r < MapHeight-1 && c < MapWidth-1 && b[nextR][c] && b[r][nextC] && b[nextR][nextC] {
					if b[nextR][c] && b[r][nextC] &&
						(j[nextR][nextC].distant[Down] > 0 ||
							j[nextR][nextC].distant[Right] > 0) {
						j[r][c].distant[DownRight] = 1
					} else {
						jumpDistance := j[nextR][nextC].distant[DownRight]
						if jumpDistance > 0 {
							j[r][c].distant[DownRight] = jumpDistance + 1
						} else {
							j[r][c].distant[DownRight] = jumpDistance - 1
						}
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
	var i uint
	for r := 0; r < MapHeight; r++ {
		for c := 0; c < MapWidth; c++ {
			jump := j[r][c]
			var blocked int
			for i = 0; i < 8; i++ {
				if 0 == jump.distant[i] {
					blocked |= (1 << i)
				}
			}
			j[r][c].blocked = blocked
		}
	}
}
