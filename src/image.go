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
		fmt.Printf("err = %v\n", err.Error())
		return nil
	}
	defer f.Close()
	img, _, _ := image.Decode(f)
	return img
}

func parseImage(img image.Image) []bool {
	max := uint32(65536 - 1) // 2^16-1

	bounds := img.Bounds()
	map_data := make([]bool, bounds.Max.X*bounds.Max.Y)
	// fmt.Printf("map_data = %v\n", map_data)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			if r == max && g == max && b == max && a == max {
				map_data[x+bounds.Max.Y*y] = true
				//fmt.Printf(".")
			} else {
				map_data[x+bounds.Max.Y*y] = false
				//fmt.Printf("#")
			}
		}
		//fmt.Println("")
	}
	return map_data
}

func GetMapFromImage(filename string) ([]bool, int, int) {
	img := openImage(filename)
	if img == nil {
		return []bool{}, 0, 0
	}
	return parseImage(img), img.Bounds().Max.X, img.Bounds().Max.Y
}
