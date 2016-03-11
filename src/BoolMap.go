package jpsplus

type BoolMap [][]bool

var DefaultBoolMap *BoolMap

func (*BoolMap) init(width int, height int) {
	map_data := make(BoolMap, height)
	for pos := 0; pos < height; pos++ {
		map_data[pos] = make([]bool, width)
	}
	DefaultBoolMap = &map_data
}

func (b BoolMap) insertTrue(row int, col int) {
	b[row][col] = true
}

func (b BoolMap) insertFalse(row int, col int) {
	b[row][col] = false
}

func (b BoolMap) isPassage(row int, col int) bool {
	if col < 0 {
		return false
	} else {
		if col >= b.width() {
			return false
		} else {
			if row < 0 {
				return false
			} else {
				if row >= b.height() {
					return false
				} else {
					return b[row][col]
				}
			}
		}
	}
}

func (b BoolMap) width() int {
	return len(b[0])
}
func (b BoolMap) height() int {
	return len(b)
}

func (b BoolMap) toString() string {
	var str string
	for _, row := range b {
		for _, cell := range row {
			if cell {
				str += "."
			} else {
				str += "#"
			}
		}
		str += "\n"
	}
	return str
}

func IsEmpty(r int, c int) bool {
	return DefaultBoolMap.isPassage(r, c)
}

func IsWall(r int, c int) bool {
	return !DefaultBoolMap.isPassage(r, c)
}
