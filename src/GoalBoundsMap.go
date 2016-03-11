package jpsplus

import (
	"fmt"
)

const (
	MinRow = iota
	MaxRow
	MinCol
	MaxCol
)

type GoalBounds [8][4]int

func (g GoalBounds) set(r int, c int, v int) {
	g[r][c] = v
}

func (g GoalBounds) get(r int, c int) int {
	return g[r][c]
}

type GoalBoundsMap [][]*GoalBounds

var DefautGoalBounds *GoalBoundsMap

func (g GoalBoundsMap) width() int {
	return len(g[0])
}

func (g GoalBoundsMap) height() int {
	return len(g)
}

func (g GoalBoundsMap) getMinRow(r int, c int, dir int) int {
	return g[r][c].get(dir, MinRow)
}

func (g GoalBoundsMap) getMaxRow(r int, c int, dir int) int {
	return g[r][c].get(dir, MaxRow)
}

func (g GoalBoundsMap) getMinCol(r int, c int, dir int) int {
	return g[r][c].get(dir, MinCol)
}

func (g GoalBoundsMap) getMaxCol(r int, c int, dir int) int {
	return g[r][c].get(dir, MaxCol)
}

func (g GoalBoundsMap) setMinRow(r int, c int, dir int, value int) {
	g[r][c].set(dir, MinRow, value)
}

func (g GoalBoundsMap) setMaxRow(r int, c int, dir int, value int) {
	g[r][c].set(dir, MaxRow, value)
}

func (g GoalBoundsMap) setMinCol(r int, c int, dir int, value int) {
	g[r][c].set(dir, MinCol, value)
}

func (g GoalBoundsMap) setMaxCol(r int, c int, dir int, value int) {
	g[r][c].set(dir, MaxCol, value)
}

func (g GoalBoundsMap) get(r int, c int) *GoalBounds {
	return g[r][c]
}

func (*GoalBoundsMap) init() {
	if nil != DefaultDistantJumpPoint {
		height := DefaultDistantJumpPoint.height()
		width := DefaultDistantJumpPoint.width()
		map_data := make(GoalBoundsMap, height)
		for pos := 0; pos < height; pos++ {
			map_data[pos] = make([]*GoalBounds, width)
		}
		for r := 0; r < height; r++ {
			for c := 0; c < width; c++ {
				map_data[r][c] = new(GoalBounds)
				for dir := 0; dir < 8; dir++ {
					// map_data[r][c].bounds[dir][MinRow] = height
					// map_data[r][c].bounds[dir][MaxRow] = 0
					// map_data[r][c].bounds[dir][MinCol] = width
					// map_data[r][c].bounds[dir][MaxCol] = 0
					map_data.setMinRow(r, c, dir, height)
					map_data.setMinCol(r, c, dir, width)
				}
			}
		}
		DefautGoalBounds = &map_data
	} else {
		fmt.Println("init DefautGoalBounds faild : not DistantJumpPointMap")
	}
}

func (g *GoalBoundsMap) CalculateGoalBounding() {
	fmt.Printf("Goal Bounding Preprocessing\n")

	DefaultDijkstra.init()

	// InitArray(m_goalBoundsMap, m_width, m_height);
	height := g.height()
	width := g.width()

	for startRow := 0; startRow < height; startRow++ {
		// fmt.Printf("Row: %d\n", startRow)

		for startCol := 0; startCol < width; startCol++ {
			if IsEmpty(startRow, startCol) {
				DefaultDijkstra.Flood(startRow, startCol)
				currentIteration := DefaultDijkstra.GetCurrentInteration()

				for r := 0; r < height; r++ {
					for c := 0; c < width; c++ {
						if IsEmpty(r, c) {
							iteration := DefaultDijkstra.m_mapNodes.node(r, c).m_iteration
							status := DefaultDijkstra.m_mapNodes.node(r, c).m_listStatus
							dir := DefaultDijkstra.m_mapNodes.node(r, c).m_directionFromStart

							if iteration == currentIteration && status == PathfindingNode_OnClosed && dir >= 0 && dir <= 7 {
								row := DefaultDijkstra.m_mapNodes.node(r, c).m_row
								col := DefaultDijkstra.m_mapNodes.node(r, c).m_col

								if g.getMinRow(startRow, startCol, dir) > row {
									g.setMinRow(startRow, startCol, dir, row)
								}
								if g.getMaxRow(startRow, startCol, dir) < row {
									g.setMaxRow(startRow, startCol, dir, row)

								}
								if g.getMinCol(startRow, startCol, dir) > col {
									g.setMinCol(startRow, startCol, dir, col)
								}
								if g.getMaxCol(startRow, startCol, dir) < col {
									g.setMaxCol(startRow, startCol, dir, col)
								}
							}
						}

					}
				}
			}

		}
	}

}
