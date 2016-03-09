package jpsplus

import (
	"testing"
)

func TestJPSplus(*testing.T) {
	// mapData := make([]bool, 0, 10000)
	// thePath := make([]Point, 0, 10000)

	mapData, width, height := LoadMap("mapFilename")
	PreprocessMap(mapData, width, height, "mapPreprocessedFilename")
	reference := PrepareForSearch(mapData, width, height, "mapPreprocessedFilename")
	thePath, ok := GetPath(reference, Point{1, 1}, Point{10, 10})
}
