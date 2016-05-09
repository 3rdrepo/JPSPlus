package jpsplus

import (
	"fmt"
	"testing"
	"time"
)

const (
	k  = 1000
	us = k
	ms = us * k
	s  = ms * k
)

func TestJpsplus(t *testing.T) {

	start := time.Now().UnixNano()
	map_data := GetMapFromImage("../map/map100x100.png")
	end := time.Now().UnixNano()
	printTime("open file", start, end)
	if map_data == nil {
		t.Errorf("Could not open map")
		return
	}
	// fmt.Println(map_data)
	start = time.Now().UnixNano()
	p := PreprocessMap(map_data)
	end = time.Now().UnixNano()
	printTime("PreprocessMap", start, end)

	// fmt.Println(p)
	start = time.Now().UnixNano()
	path, ok := p.GetPath(0, 0, 99, 99)
	end = time.Now().UnixNano()
	printTime("GetPath", start, end)
	if ok {
		fmt.Println(str_map(map_data, path))
		// fmt.Println(path)
	} else {
		fmt.Println("not path !")
	}
}

func str_map(data *BoolMap, path []LocJPS) (strMap string) {
	for r := 0; r < MapHeight; r++ {
		for c := 0; c < MapWidth; c++ {
			noPath := true
			for _, loc := range path {
				if loc.Row == r && loc.Col == c {
					strMap += "o"
					noPath = false
					break
				}
			}
			if noPath {
				if data[r][c] {
					strMap += "."
				} else {
					strMap += "#"
				}
			}
		}
		strMap += "\n"
	}
	return
}

func printTime(str string, start int64, end int64) {
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
