package jpsplus

import (
	"fmt"
)

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

type DistantJumpPoints [8]int

func (d DistantJumpPoints) get(i int) int {
	return d[i]
}

func (d DistantJumpPoints) set(i int, v int) {
	d[i] = v
}

type DistantJumpPointMap [][]*DistantJumpPoints

var DefaultDistantJumpPoint *DistantJumpPointMap

func (*DistantJumpPointMap) init() {
	if nil != DefaultJumpPoint {
		height := DefaultJumpPoint.height()
		width := DefaultJumpPoint.width()
		map_data := make(DistantJumpPointMap, height)
		for r := 0; r < height; r++ {
			map_data[r] = make([]*DistantJumpPoints, width)
			for c := 0; c < width; c++ {
				map_data[r][c] = new(DistantJumpPoints)
			}
		}
		DefaultDistantJumpPoint = &map_data
	} else {
		fmt.Println("init DefaultDistantJumpPoint faild : not DefaultJumpPoint")
	}
}

func (d DistantJumpPointMap) width() int {
	return len(d[0])
}

func (d DistantJumpPointMap) height() int {
	return len(d)
}

func (d DistantJumpPointMap) setJumpdistanceLeft(r int, c int, distance int) {
	d[r][c].set(Left, distance)
}

func (d DistantJumpPointMap) setJumpdistanceRight(r int, c int, distance int) {
	d[r][c].set(Right, distance)
}

func (d DistantJumpPointMap) setJumpdistanceUp(r int, c int, distance int) {
	d[r][c].set(Up, distance)
}

func (d DistantJumpPointMap) setJumpdistanceDown(r int, c int, distance int) {
	d[r][c].set(Down, distance)
}

func (d DistantJumpPointMap) setJumpdistanceUpLeft(r int, c int, distance int) {
	d[r][c].set(UpLeft, distance)
}

func (d DistantJumpPointMap) setJumpdistanceUpRight(r int, c int, distance int) {
	d[r][c].set(UpRight, distance)
}

func (d DistantJumpPointMap) setJumpdistanceDownLeft(r int, c int, distance int) {
	d[r][c].set(DownLeft, distance)
}

func (d DistantJumpPointMap) setJumpdistanceDownRight(r int, c int, distance int) {
	d[r][c].set(DownRight, distance)
}

func (d DistantJumpPointMap) getJumpdistanceLeft(r int, c int) int {
	return d[r][c].get(Left)
}

func (d DistantJumpPointMap) getJumpdistanceRight(r int, c int) int {
	return d[r][c].get(Right)
}

func (d DistantJumpPointMap) getJumpdistanceUp(r int, c int) int {
	return d[r][c].get(Up)
}

func (d DistantJumpPointMap) getJumpdistanceDown(r int, c int) int {
	return d[r][c].get(Down)
}

func (d DistantJumpPointMap) getJumpdistanceUpLeft(r int, c int) int {
	return d[r][c].get(UpLeft)
}

func (d DistantJumpPointMap) getJumpdistanceUpRight(r int, c int) int {
	return d[r][c].get(UpRight)
}

func (d DistantJumpPointMap) getJumpdistanceDownLeft(r int, c int) int {
	return d[r][c].get(DownLeft)
}

func (d DistantJumpPointMap) getJumpdistanceDownRight(r int, c int) int {
	return d[r][c].get(DownRight)
}

func (d DistantJumpPointMap) getJumpdistance(r int, c int, dir int) int {
	return d[r][c].get(dir)
}

func (d DistantJumpPointMap) get(r int, c int) *DistantJumpPoints {
	return d[r][c]
}

func (d *DistantJumpPointMap) CalculateDistantJumpPointMapLeft() {
	var countMovingLeft int
	var jumpPointLastSeen bool
	width := d.width()
	height := d.height()
	for r := 0; r < height; r++ {
		countMovingLeft = -1
		jumpPointLastSeen = false
		for c := 0; c < width; c++ {
			if IsWall(r, c) {
				countMovingLeft = -1
				jumpPointLastSeen = false
				d.setJumpdistanceLeft(r, c, 0)
			} else {
				countMovingLeft += 1
				if jumpPointLastSeen {
					d.setJumpdistanceLeft(r, c, countMovingLeft)
				} else {
					// Wall last seen
					d.setJumpdistanceLeft(r, c, -countMovingLeft)
				}
				if (DefaultJumpPoint.node(r, c) & MovingLeft) > 0 {
					countMovingLeft = 0
					jumpPointLastSeen = true
				}
			}
		}
	}
}

func (d *DistantJumpPointMap) CalculateDistantJumpPointMapRight() {
	var countMovingRight int
	var jumpPointLastSeen bool
	width := d.width()
	height := d.height()
	for r := 0; r < height; r++ {
		countMovingRight = -1
		jumpPointLastSeen = false
		for c := width - 1; c >= 0; c-- {
			if IsWall(r, c) {
				countMovingRight = -1
				jumpPointLastSeen = false
				d.setJumpdistanceRight(r, c, 0)
			} else {
				countMovingRight += 1
				if jumpPointLastSeen {
					d.setJumpdistanceRight(r, c, countMovingRight)
				} else {
					// Wall last see
					d.setJumpdistanceRight(r, c, -countMovingRight)
				}
				if (DefaultJumpPoint.node(r, c) & MovingRight) > 0 {
					countMovingRight = 0
					jumpPointLastSeen = true
				}
			}
		}
	}
}

func (d *DistantJumpPointMap) CalculateDistantJumpPointMapUp() {
	var countMovingUp int
	var jumpPointLastSeen bool
	width := d.width()
	height := d.height()
	for c := 0; c < width; c++ {
		countMovingUp = -1
		jumpPointLastSeen = false
		for r := 0; r < height; r++ {
			if IsWall(r, c) {
				countMovingUp = -1
				jumpPointLastSeen = false
				// p.m_distantJumpPointMap[r][c].jumpDistance[Up] = 0
				d.setJumpdistanceUp(r, c, 0)
			} else {
				countMovingUp += 1
				if jumpPointLastSeen {
					d.setJumpdistanceUp(r, c, countMovingUp)
				} else { // Wall last seen
					d.setJumpdistanceUp(r, c, -countMovingUp)
				}
				if (DefaultJumpPoint.node(r, c) & MovingUp) > 0 {
					countMovingUp = 0
					jumpPointLastSeen = true
				}
			}
		}
	}
}

func (d *DistantJumpPointMap) CalculateDistantJumpPointMapDown() {
	var countMovingDown int
	var jumpPointLastSeen bool
	width := d.width()
	height := d.height()
	for c := 0; c < width; c++ {
		countMovingDown = -1
		jumpPointLastSeen = false
		for r := height - 1; r >= 0; r-- {
			if IsWall(r, c) {
				countMovingDown = -1
				jumpPointLastSeen = false
				d.setJumpdistanceDown(r, c, 0)
			} else {
				countMovingDown += 1
				if jumpPointLastSeen {
					d.setJumpdistanceDown(r, c, countMovingDown)
				} else { // Wall last seen
					d.setJumpdistanceDown(r, c, -countMovingDown)
				}
				if (DefaultJumpPoint.node(r, c) & MovingDown) > 0 {
					countMovingDown = 0
					jumpPointLastSeen = true
				}
			}
		}
	}
}

func (d *DistantJumpPointMap) CalculateDistantJumpPointMapUpLeftandUpRight() {
	width := d.width()
	height := d.height()
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			if IsEmpty(r, c) {
				//UpLeft
				if r == 0 || c == 0 || (IsWall(r-1, c) || IsWall(r, c-1) || IsWall(r-1, c-1)) {
					// Wall one away
					d.setJumpdistanceUpLeft(r, c, 0)
				} else if IsEmpty(r-1, c) && IsEmpty(r, c-1) &&
					(d.getJumpdistanceUp(r-1, c-1) > 0 ||
						d.getJumpdistanceLeft(r-1, c-1) > 0) {
					// Diagonal one away
					d.setJumpdistanceUpLeft(r, c, 1)
				} else {
					// Increment from last
					jumpDistance := d.getJumpdistanceUpLeft(r-1, c-1)

					if jumpDistance > 0 {
						d.setJumpdistanceUpLeft(r, c, jumpDistance+1)
					} else {
						//if( jumpDistance <= 0 )
						d.setJumpdistanceUpLeft(r, c, jumpDistance-1)
					}
				}

				//UpRight
				if r == 0 || c == width-1 || (IsWall(r-1, c) || IsWall(r, c+1) || IsWall(r-1, c+1)) {
					// Wall one away
					d.setJumpdistanceUpRight(r, c, 0)
				} else if IsEmpty(r-1, c) && IsEmpty(r, c+1) &&
					(d.getJumpdistanceUp(r-1, c+1) > 0 ||
						d.getJumpdistanceRight(r-1, c+1) > 0) {
					// Diagonal one away
					// p.m_distantJumpPointMap[r][c].jumpDistance[UpRight] = 1
					d.setJumpdistanceUpRight(r, c, 1)
				} else {
					// Increment from last
					// jumpDistance := p.m_distantJumpPointMap[r-1][c+1].jumpDistance[UpRight]
					jumpDistance := d.getJumpdistanceUpRight(r-1, c+1)
					if jumpDistance > 0 {
						// p.m_distantJumpPointMap[r][c].jumpDistance[UpRight] = 1 + jumpDistance
						d.setJumpdistanceUpRight(r, c, jumpDistance+1)
					} else {
						//if( jumpDistance <= 0 )
						// p.m_distantJumpPointMap[r][c].jumpDistance[UpRight] = -1 + jumpDistance
						d.setJumpdistanceUpRight(r, c, jumpDistance-1)
					}
				}
			}
		}
	}
}
func (d *DistantJumpPointMap) CalculateDistantJumpPointMapDownLeftandDownRight() {
	width := d.width()
	height := d.height()
	for r := height - 1; r >= 0; r-- {
		for c := 0; c < width; c++ {
			if IsEmpty(r, c) {
				if r == height-1 || c == 0 ||
					(IsWall(r+1, c) || IsWall(r, c-1) || IsWall(r+1, c-1)) {
					// Wall one away
					d.setJumpdistanceDownLeft(r, c, 0)
				} else if IsEmpty(r+1, c) && IsEmpty(r, c-1) &&
					(d.getJumpdistanceDown(r+1, c-1) > 0 ||
						d.getJumpdistanceLeft(r+1, c-1) > 0) {
					// Diagonal one away
					d.setJumpdistanceDownLeft(r, c, 1)
				} else {
					// Increment from last
					// jumpDistance := p.m_distantJumpPointMap[r+1][c-1].jumpDistance[DownLeft]
					jumpDistance := d.getJumpdistanceDownLeft(r+1, c-1)
					if jumpDistance > 0 {
						// p.m_distantJumpPointMap[r][c].jumpDistance[DownLeft] = 1 + jumpDistance
						d.setJumpdistanceDownLeft(r, c, jumpDistance+1)
					} else { //if( jumpDistance <= 0 )
						d.setJumpdistanceDownLeft(r, c, jumpDistance-1)
					}
				}

				if r == height-1 || c == width-1 || (IsWall(r+1, c) || IsWall(r, c+1) || IsWall(r+1, c+1)) {
					// Wall one away
					d.setJumpdistanceDownRight(r, c, 0)
				} else if IsEmpty(r+1, c) && IsEmpty(r, c+1) &&
					(d.getJumpdistanceDown(r+1, c+1) > 0 ||
						d.getJumpdistanceRight(r+1, c+1) > 0) {
					// Diagonal one away
					d.setJumpdistanceDownRight(r, c, 1)
				} else {
					// Increment from last
					jumpDistance := d.getJumpdistanceDownRight(r+1, c+1)

					if jumpDistance > 0 {
						d.setJumpdistanceDownRight(r, c, jumpDistance+1)
					} else { //if( jumpDistance <= 0 )
						d.setJumpdistanceDownRight(r, c, jumpDistance-1)

					}
				}
			}
		}
	}
}

func (d *DistantJumpPointMap) CalculateDistantJumpPointMap() {
	// Calculate distant jump points (Left and Right)
	d.CalculateDistantJumpPointMapLeft()
	d.CalculateDistantJumpPointMapRight()

	// Calculate distant jump points (Up and Down)
	d.CalculateDistantJumpPointMapUp()
	d.CalculateDistantJumpPointMapDown()

	// Calculate distant jump points (Diagonally UpLeft and UpRight)
	d.CalculateDistantJumpPointMapUpLeftandUpRight()
	// Calculate distant jump points (Diagonally DownLeft and DownRight)
	d.CalculateDistantJumpPointMapDownLeftandDownRight()

}
