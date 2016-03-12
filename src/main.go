package main

import (
	"./jpsplus"
	"fmt"
	"time"
)

func main() {

	// fmt.Printf("width = %v  heigth = %v\n", DefaultBoolMap.width(), DefaultBoolMap.height())
	// fmt.Println(DefaultBoolMap.toString())
	timePreprocessMapStart := time.Now().UnixNano()

	jpsplus.PreprocessMap("../map/map3kx3k.png")

	timePreprocessMapEnd := time.Now().UnixNano()
	fmt.Printf("timePreprocessMap  %v ms\n", (timePreprocessMapEnd-timePreprocessMapStart)/1000000.0)

	jps := jpsplus.NewJPSPlus()
	s := jpsplus.LocJPS{0, 0}
	g := jpsplus.LocJPS{2999, 2999}

	timeGetPathStart := time.Now().UnixNano()

	_, ok := jps.GetPath(s, g)

	timeGetPathEnd := time.Now().UnixNano()
	jpsplus.PrintTime("GetPath", timeGetPathStart, timeGetPathEnd)
	fmt.Printf("ok = %v\n ", ok)
}
