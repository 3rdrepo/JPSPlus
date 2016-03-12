package jpsplus

import (
	"fmt"
	"image"
	"os"
)

import _ "image/png"

func openImage(filename string) image.Image {
	fmt.Println(filename)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("err = %v\n", err.Error())
		return nil
	}
	defer f.Close()
	img, _, _ := image.Decode(f)
	return img
}

func parseImage(img image.Image) {
	max := uint32(65536 - 1) // 2^16-1

	bounds := img.Bounds()
	DefaultBoolMap.init(bounds.Max.Y, bounds.Max.X)
	// fmt.Printf("bool map w %v h %v\n", DefaultBoolMap.width(), DefaultBoolMap.height())
	// fmt.Printf("map_data = %#v\n", bounds)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			if r == max && g == max && b == max && a == max {
				// fmt.Println(y, x)
				DefaultBoolMap.insertTrue(y, x)
				//fmt.Printf(".")
			} else {
				DefaultBoolMap.insertFalse(y, x)
				// map_data[y][x] = false
				//fmt.Printf("#")
			}
		}
		//fmt.Println("")
	}
}

func GetMapFromImage(filename string) bool {
	img := openImage(filename)
	if nil == img {
		return false
	}
	parseImage(img)
	if nil == DefaultBoolMap {
		return false
	} else {
		return true
	}
}
