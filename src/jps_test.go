package jpsplus

import (
	"fmt"
	"testing"
)

func TestJPSplus(*testing.T) {
	// mapData := make([]bool, 0, 10000)
	// thePath := make([]Point, 0, 10000)

	mapData, width, height := GetMapFromImage("map100_100P.png")
	fmt.Printf("width = %v height = %v\n", width, height)
	// fmt.Printf("mapData = %v\n", mapData)
	PreprocessMap(mapData, width, height, "mapPreprocessedFilename")
	// reference := PrepareForSearch(mapData, width, height, "mapPreprocessedFilename")
	// thePath, ok := GetPath(reference, Point{1, 1}, Point{10, 10})
}
