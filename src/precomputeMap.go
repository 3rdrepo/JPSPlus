package jpsplus

import (
	"fmt"
)

func PreprocessMap(filename string) {
	fmt.Printf("Writing to file '%s'\n", filename)
	// precomputeMap := newPrecomputeMap()
	// fmt.Printf("precomputeMap = %#v\n", precomputeMap)

	CalculateMap()
	// fmt.Printf("precomputeMap = %#v\n", precomputeMap)
	// precomputeMap.SaveMap(filename)
}

func CalculateMap() {
	DefaultJumpPoint.init()
	DefaultJumpPoint.CalculateJumpPointMap()

	DefaultDistantJumpPoint.init()
	DefaultDistantJumpPoint.CalculateDistantJumpPointMap()

	DefautGoalBounds.init()
	DefautGoalBounds.CalculateGoalBounding()
	initDefaultJumpDistancesAndGoalBounds()
	// fmt.Printf("%#v\n", DefaultDistantJumpPoint)
	// fmt.Println("")
	// fmt.Printf("%#v\n", DefautGoalBounds)
}
