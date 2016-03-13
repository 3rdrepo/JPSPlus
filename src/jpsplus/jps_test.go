package jpsplus

import (
	"fmt"
	"testing"
	"time"
)

func TestJPSplus(*testing.T) {

	// fmt.Println(DefaultBoolMap.toString())
	timePreprocessMapStart := time.Now().UnixNano()
	PreprocessMap("../../map/map3kx3k.png")
	timePreprocessMapEnd := time.Now().UnixNano()
	PrintTime("timePreprocessMap", timePreprocessMapStart, timePreprocessMapEnd)

	jps := NewJPSPlus()

	s := LocJPS{0, 0}
	g := LocJPS{2999, 2999}
	timeGetPathStart := time.Now().UnixNano()

	_, ok := jps.GetPath(s, g)

	timeGetPathEnd := time.Now().UnixNano()

	PrintTime("GetPath", timeGetPathStart, timeGetPathEnd)

	fmt.Printf("ok = %v\n ", ok)
}

func str_map(path []LocJPS) string {
	var result string
	for r, data := range *DefaultBoolMap {
		for c, cell := range data {
			ismap := true
			for _, node := range path {
				if node.X == c && node.Y == r {
					result += "o"
					ismap = false
					break
				}
			}
			if ismap {
				if cell {
					result += "."
				} else {
					result += "#"
				}
			}
		}
		result += "\n"
	}
	return result
}
