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

const (
	Map10  = "../map/map10x10.png"
	Map100 = "../map/map100x100.png"
	Map3k  = "../map/map3kx3k.png"
)

var JPSmap = initJPSplus()
var Bmap = initBoolMap()

func initBoolMap() *BoolMap {
	start := time.Now().UnixNano()
	bmap := GetMapFromImage(Map100)
	end := time.Now().UnixNano()
	printTime("open file", start, end)
	return bmap
}

func initJPSplus() *JumpMap {
	start := time.Now().UnixNano()
	p := PreprocessMap(Bmap)
	end := time.Now().UnixNano()
	printTime("PreprocessMap", start, end)
	return p
}

func TestJpsplus(t *testing.T) {
	start := time.Now().UnixNano()
	// path, ok := findPath(9999, 9999, 0, 0)
	path, ok := JPSmap.GetPath(0, 0, 99, 0)

	end := time.Now().UnixNano()
	printTime("findPath", start, end)
	if ok {
		fmt.Println(path)
	} else {
		fmt.Println("not path !")
	}
}

func findPath(sx int, sy int, gx int, gy int) (path map[int]*Loc, isFind bool) {
	sRow, sCol := logicToTile(sx, sy)
	gRow, gCol := logicToTile(gx, gy)
	sLogic := newLoc(sx, sy)
	gLogic := newLoc(gx, gy)
	if sRow == gRow && sCol == gCol {
		path = make(map[int]*Loc)
		path[0] = sLogic
		path[1] = gLogic
		isFind = true
	} else {
		start := time.Now().UnixNano()
		tilePath, ok := JPSmap.GetPath(sRow, sCol, gRow, gCol)
		end := time.Now().UnixNano()
		printTime("GetPath", start, end)
		if ok {
			logicPath := tilePathToLogic(tilePath, sLogic, gLogic)
			isFind = true
			if 2 == len(logicPath) {
				path = logicPath
			} else {
				start := time.Now().UnixNano()
				path = smoothPath(logicPath)
				end := time.Now().UnixNano()
				printTime("smoothPath", start, end)
			}
		}
	}
	return
}

func str_path(path map[int]*Loc) (strMap string) {
	for r := 0; r < MapHeight; r++ {
		for c := 0; c < MapWidth; c++ {
			noPath := true
			for _, loc := range path {
				row, col := logicToTile(loc.X, loc.Y)
				if row == r && col == c {
					strMap += "o"
					noPath = false
					break
				}
			}
			if noPath {
				if Bmap[r][c] {
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

type Loc struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func newLoc(x int, y int) *Loc {
	l := new(Loc)
	l.X = x
	l.Y = y
	return l
}

func tilePathToLogic(tilePath map[int]*LocJPS, sLogic *Loc, gLogic *Loc) (locPath map[int]*Loc) {
	locPath = make(map[int]*Loc)
	locPath[0] = gLogic
	max := len(tilePath) - 1
	locPath[max] = sLogic
	for i := 1; i < max; i++ {
		x, y := tileToLogic(tilePath[i].row, tilePath[i].col)
		locPath[i] = newLoc(x+TileWidth/2, y+TileHeight/2)
	}
	return
}

func smoothPath(logicPath map[int]*Loc) (path map[int]*Loc) {
	path = make(map[int]*Loc)
	max := len(logicPath) - 1
	pathPos := 0
	for j := max; j > 1; j-- {
		// fmt.Printf("logicPath[j] =%v  %#v\n", j, logicPath[j])
		path[pathPos] = logicPath[j]
		pathPos++
		for i := 0; i < j-1; i++ {
			if lineCrossTile(logicPath[j], logicPath[i]) {
				// fmt.Printf("lazhi  =%v  %v\n", j, i)
				j = i + 1
				break
			}
		}
	}
	// fmt.Println(path)
	// fmt.Println(pathPos)
	path[pathPos] = logicPath[1]
	path[pathPos+1] = logicPath[0]
	return path
}

func lineCrossTile(start *Loc, goal *Loc) bool {
	if abs(goal.X-start.X) > abs(goal.Y-start.Y) {
		return lineCrossTileCol(start, goal)
	} else {
		return lineCrossTileRow(start, goal)
	}
}

func lineCrossTileCol(start *Loc, goal *Loc) bool {
	incX := TileWidth
	posX := start.X/incX*incX + incX
	endX := goal.X/incX*incX + incX
	prexC := -1
	if goal.X < start.X {
		incX = -incX
		posX = posX + incX
		endX = endX + incX
		prexC = -prexC
	}
	incY := (goal.Y - start.Y) * incX / (goal.X - start.X)
	posY := start.Y + (posX-start.X)*incY/incX
	for ; posX != endX; posX, posY = posX+incX, posY+incY {
		row, col := logicToTile(posX, posY)
		if posY%TileHeight != 0 {
			if !Bmap[row][col+prexC] {
				return false
			}
		}
		if !Bmap[row][col] {
			return false
		}
	}
	return true
}

func lineCrossTileRow(start *Loc, goal *Loc) bool {
	incY := TileHeight
	posY := start.Y/incY*incY + incY
	endY := goal.Y/incY*incY + incY
	prexR := -1
	if goal.Y < start.Y {
		incY = -incY
		posY = posY + incY
		endY = endY + incY
		prexR = -prexR
	}
	incX := (goal.X - start.X) * incY / (goal.Y - start.Y)
	posX := start.X + (posY-start.Y)*incX/incY
	for ; posY != endY; posX, posY = posX+incX, posY+incY {
		row, col := logicToTile(posX, posY)
		if posX%TileWidth != 0 {
			if !Bmap[row+prexR][col] {
				return false
			}
		}
		if !Bmap[row][col] {
			return false
		}
	}
	return true
}
