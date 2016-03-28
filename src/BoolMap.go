package jpsplus

type BoolMap [][]bool

func newBoolMap(m Map) *BoolMap {
	row := m.Height()
	col := m.Width()
	b := make(BoolMap, row)
	for r := 0; r < row; r++ {
		b[r] = make([]bool, col)
		for c := 0; c < col; c++ {
			b[r][c] = m.IsEmpty(r, c)
		}
	}
	return &b
}

func (b BoolMap) width() int {
	return len(b[0])
}
func (b BoolMap) height() int {
	return len(b)
}

func (b BoolMap) IsEmpty(r int, c int) bool {
	if r < 0 {
		return false
	} else {
		if c < 0 {
			return false
		} else {
			if r >= b.height() {
				return false
			} else {
				if c >= b.width() {
					return false
				} else {
					return b[r][c]
				}
			}
		}
	}
}

func (b BoolMap) IsJumpPoint(r int, c int, rowDir int, colDir int) bool {

	return ((b.IsEmpty(r+colDir, c+rowDir) &&
		!b.IsEmpty(r-rowDir+colDir, c-colDir+rowDir)) ||
		(b.IsEmpty(r-colDir, c-rowDir) &&
			!b.IsEmpty(r-rowDir-colDir, c-colDir-rowDir)))
}
