package jpsplus

type BoolMap [MapHeight][MapWidth]bool

func (b BoolMap) IsEmpty(r int, c int) bool {
	if r < 0 || c < 0 || r >= MapHeight || c >= MapWidth {
		return false
	} else {
		return b[r][c]
	}
}

func (b BoolMap) IsJumpPoint(r int, c int, rowDir int, colDir int) bool {
	return ((b.IsEmpty(r+colDir, c+rowDir) &&
		!b.IsEmpty(r-rowDir+colDir, c-colDir+rowDir)) ||
		(b.IsEmpty(r-colDir, c-rowDir) &&
			!b.IsEmpty(r-rowDir-colDir, c-colDir-rowDir)))
}

func (b BoolMap) String() (strMap string) {
	for r := 0; r < MapHeight; r++ {
		for c := 0; c < MapWidth; c++ {
			if b[r][c] {
				strMap += "."
			} else {
				strMap += "#"
			}
		}
		strMap += "\n"
	}
	return
}
