package jpsplus

import (
	"fmt"
)

type Map interface {
	IsEmpty(int, int) bool
	Width() int
	Height() int
}

func PreprocessMap(m Map) {
	boolMap := newBoolMap(m)
	height := boolMap.height()
	width := boolMap.width()
	jumpPoint := newJumpPoint(width, height)
	jumpPoint.CalculateJumpPointMap(boolMap)
	jumpMap := NewJumpMap(width, height)
	jumpMap.CalculateDistantJumpPointMap(boolMap, jumpPoint)
	jumpMap.CalculateBlock()
	DefaultJumpMap = jumpMap

	jps := NewJPSPlus()
	path, ok := jps.GetPath(483, 694, 562, 819)
	fmt.Println(path, ok)
	// for _, d := range *boolMap {
	// 	for _, c := range d {
	// 		fmt.Println(c)
	// 	}
	// }
}
