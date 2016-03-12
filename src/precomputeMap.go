package jpsplus

import (
	"fmt"
)

func PreprocessMap(filename string) {
	fmt.Printf("TODO:Writing to file '%s'\n", filename)
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

	// w := DefautGoalBounds.width()
	// h := DefautGoalBounds.height()
	// for r := 0; r < h; r++ {
	// 	for c := 0; c < w; c++ {
	// 		fmt.Println("DefaultJumpPoint")
	// 		fmt.Printf("%v\n", DefaultJumpPoint.get(r, c))
	// 		fmt.Println("")

	// 		fmt.Println("DefautGoalBounds")
	// 		fmt.Printf("%v\n", DefautGoalBounds.get(r, c))
	// 		fmt.Println("")

	// 		fmt.Println("DefaultJumpDistancesAndGoalBounds")
	// 		fmt.Printf("%v\n", *(DefaultJumpDistancesAndGoalBounds.get(r, c).bounds))
	// 		fmt.Println("")
	// 	}
	// }

	// fmt.Printf("%#v\n", DefaultDistantJumpPoint)
	// fmt.Println("")
	// fmt.Printf("%#v\n", DefautGoalBounds)
}
