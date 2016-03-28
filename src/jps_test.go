package jpsplus

import (
	"fmt"
	"testing"
)

const (
	Test_N = 10 * 10000
)

type M [][]bool

func (m M) Width() int {
	return 4
}

func (m M) Height() int {
	return 5
}

func (m M) IsEmpty(r int, c int) bool {
	return m[r][c]
}

func TestJPSplus(*testing.T) {
	m := M{
		{true, true, true, true},
		{true, false, true, true},
		{true, true, false, true},
		{true, false, false, true},
		{true, true, true, true},
	}
	PreprocessMap(m)
	jps := NewJPSPlus()
	ch := make(chan int, 100)
	for i := 0; i < 100; i++ {
		p, _ := jps.GetPath(0, 0, 4, 3)
		fmt.Println(*p)
		ch <- 0
		fmt.Println(len(ch))
	}
}
