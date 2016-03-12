package jpsplus

import (
	"fmt"
)

const (
	MovingDown = 1 << iota
	MovingRight
	MovingUp
	MovingLeft
)

type JumpPointMap [][]int

var DefaultJumpPoint *JumpPointMap

func (*JumpPointMap) init() {
	if nil != DefaultBoolMap {
		height := DefaultBoolMap.height()
		width := DefaultBoolMap.width()
		map_data := make(JumpPointMap, height)
		for pos := 0; pos < height; pos++ {
			map_data[pos] = make([]int, width)
		}
		DefaultJumpPoint = &map_data
	} else {
		fmt.Println("init DefaultJumpPoint faild : not initBoolMap")
	}
}

func (j JumpPointMap) width() int {
	return len(j[0])
}

func (j JumpPointMap) height() int {
	return len(j)
}

func (j JumpPointMap) moveDown(r int, c int) {
	j[r][c] |= MovingDown
}

func (j JumpPointMap) movingUp(r int, c int) {
	j[r][c] |= MovingUp
}

func (j JumpPointMap) movingRight(r int, c int) {
	j[r][c] |= MovingRight
}

func (j JumpPointMap) movingLeft(r int, c int) {
	j[r][c] |= MovingLeft
}

func (j JumpPointMap) get(r int, c int) int {
	return j[r][c]
}

func (j JumpPointMap) set(r int, c int, v int) {
	j[r][c] = v
}

func (j *JumpPointMap) CalculateJumpPointMap() {
	for r, data := range *DefaultBoolMap {
		for c, b := range data {
			if b {
				if IsJumpPoint(r, c, 1, 0) {
					j.moveDown(r, c)
				}
				if IsJumpPoint(r, c, -1, 0) {
					j.movingUp(r, c)
				}
				if IsJumpPoint(r, c, 0, 1) {
					j.movingRight(r, c)
				}
				if IsJumpPoint(r, c, 0, -1) {
					j.movingLeft(r, c)
				}
			}
		}
	}
}

func IsJumpPoint(r int, c int, rowDir int, colDir int) bool {
	// return p.IsEmpty(r-rowDir, c-colDir) && // Parent not a wall (not necessary)
	return ((IsEmpty(r+colDir, c+rowDir) && // 1st forced neighbor
		IsWall(r-rowDir+colDir, c-colDir+rowDir)) || // 1st forced neighbor (continued)
		(IsEmpty(r-colDir, c-rowDir) && // 2nd forced neighbor
			IsWall(r-rowDir-colDir, c-colDir-rowDir))) // 2nd forced neighbor (continued)
}
