package jpsplus

import (
	"fmt"
	"time"
)

const (
	k  = 1000
	us = k
	ms = us * k
	s  = ms * k
)

func PreprocessMap(filename string) {
	// precomputeMap := newPrecomputeMap()
	fmt.Printf("BoolMap Preprocessing\n")
	tbs := time.Now().UnixNano()
	ok := GetMapFromImage(filename)
	tbe := time.Now().UnixNano()
	PrintTime("BoolMap", tbs, tbe)
	fmt.Printf("width = %v  heigth = %v\n", DefaultBoolMap.width(), DefaultBoolMap.height())
	if ok {
		CalculateMap()
		fmt.Printf("TODO:Writing to file '%s'\n", filename)
	} else {
		fmt.Println("open mapfile faild")
	}
	// fmt.Printf("precomputeMap = %#v\n", precomputeMap)
	// precomputeMap.SaveMap(filename)

}

func CalculateMap() {
	fmt.Printf("JumpPoint Preprocessing\n")
	tjs := time.Now().UnixNano()
	DefaultJumpPoint.init()
	DefaultJumpPoint.CalculateJumpPointMap()
	tje := time.Now().UnixNano()
	PrintTime("JumpPoint", tjs, tje)

	fmt.Printf("DistantJumpPoint Preprocessing\n")
	tds := time.Now().UnixNano()
	DefaultDistantJumpPoint.init()
	DefaultDistantJumpPoint.CalculateDistantJumpPointMap()
	tde := time.Now().UnixNano()
	PrintTime("DistantJumpPoint", tds, tde)

	fmt.Printf("GoalBounding Preprocessing\n")
	tgs := time.Now().UnixNano()
	DefautGoalBounds.init()
	DefautGoalBounds.CalculateGoalBounding()
	tge := time.Now().UnixNano()
	PrintTime("GoalBounding", tgs, tge)

	fmt.Printf("JumpDistancesAndGoalBounds Preprocessing\n")
	tjgs := time.Now().UnixNano()
	initDefaultJumpDistancesAndGoalBounds()
	tjge := time.Now().UnixNano()
	PrintTime("JumpDistancesAndGoalBounds", tjgs, tjge)

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

func PrintTime(str string, start int64, end int64) {
	time := end - start
	if time > s {
		fmt.Printf("%s  %v s\n", str, float64(time)/s)
	} else {
		if time > ms {
			fmt.Printf("%s  %v ms\n", str, float64(time)/ms)
		} else {
			if time > us {
				fmt.Printf("%s  %v us\n", str, float64(time)/us)
			} else {
				fmt.Printf("%s  %v ns\n", str, time)
			}
		}
	}
}
