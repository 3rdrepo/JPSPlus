package jpsplus

import (
	"fmt"
	"image"
	"os"
)

import _ "image/png"

func openImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	img, _, _ := image.Decode(f)
	return img
}

func parseImage(img image.Image) *BoolMap {
	max := uint32(65536 - 1)
	bounds := img.Bounds()
	fmt.Printf("width = %v, height = %v\n", bounds.Max.Y, bounds.Max.X)
	map_data := new(BoolMap)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r == max && g == max && b == max && a == max {
				map_data[y][x] = true
				//fmt.Printf(".")
			} else {
				map_data[y][x] = false
				//fmt.Printf("#")
			}
		}
	}
	return map_data
}

func GetMapFromImage(filename string) (res *BoolMap) {
	img := openImage(filename)
	if img != nil {
		res = parseImage(img)
	}
	return
}
