package jpsplus

import (
	"fmt"
	"testing"
)

func TestJPSplus(*testing.T) {
	// mapData := make([]bool, 0, 10000)
	// thePath := make([]Point, 0, 10000)

	ok := GetMapFromImage("map10x10.png")
	if ok {
		fmt.Printf("width = %v  heigth = %v\n", DefaultBoolMap.width(), DefaultBoolMap.height())
		fmt.Println(DefaultBoolMap.toString())
		PreprocessMap("mapPreprocessedFilename")
	} else {
		fmt.Println("open mapfile faild")
	}

	// reference := PrepareForSearch(mapData, width, height, "mapPreprocessedFilename")
	// thePath, ok := GetPath(reference, Point{1, 1}, Point{10, 10})
}
