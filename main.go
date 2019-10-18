package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	fmt.Println("Starting program...")
	file, err := os.Open("images/image2.jpg")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	img, t, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Type of image:", t)

	rct := img.Bounds()
	fmt.Println("Width:", rct.Dx())
	fmt.Println("Height:", rct.Dy())
	r, g, b, a := img.At(1051, 1059).RGBA()
	fmt.Println("pixel:", r/257, g/257, b/257, a/257)

	var rr, gg, bb, aa uint32
	for i := 0; i < rct.Dx(); i++ {
		for j := 0; j < rct.Dy(); j++ {
			r, g, b, a = img.At(i, j).RGBA()
			rr += r / 257
			gg += g / 257
			bb += b / 257
			aa += a / 257
		}
	}
	m := uint32(rct.Dx() * rct.Dy())
	rr /= m
	gg /= m
	bb /= m
	aa /= m
	fmt.Println("Ave pixel", rr, gg, bb, aa)
}
