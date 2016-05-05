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
