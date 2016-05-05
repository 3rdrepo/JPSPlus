package jpsplus

import (
	"fmt"
	"testing"
)

func TestJpsplus(t *testing.T) {
	map_data := GetMapFromImage("../map/map10x10.png")
	if map_data == nil {
		t.Errorf("Could not open map")
		return
	}
	fmt.Println(map_data)

	p := PreprocessMap(map_data)
	fmt.Println(p)
	// jps := NewJPSPlus()
	// path, ok := jps.GetPath(483, 694, 562, 819)
	// fmt.Println(path, ok)
}

func str_map(data *BoolMap) (strMap string) {
	for r := 0; r < MapHeight; r++ {
		for c := 0; c < MapWidth; c++ {
			if data[r][c] {
				strMap += "."
			} else {
				strMap += "#"
			}
		}
		strMap += "\n"
	}
	return
}
