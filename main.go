package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	inputImagePath, side, err := getFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	originalFile, err := os.Open(inputImagePath)
	defer originalFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	originalImage, _, err := image.Decode(originalFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	originalRect := originalImage.Bounds()
	width, height, rect := getNewRect(originalRect, side)

	colors := getColors(originalImage, width, height, side)

	newImage := createNewImage(colors, side, width, height, rect)

	name := getNewImageName(originalFile.Name(), side)
	f, err := os.Create("./outputs/" + name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	png.Encode(f, newImage)

	fmt.Println("Tiled!")
	fmt.Println("Here is tiled image: " + f.Name())
}

func getFlags() (inputPath string, sideNum int, err error) {
	inputImagePath := flag.String("i", "", "Path to original image")
	side := flag.Int("n", 4, "Number for converting nxn tile")

	flag.Parse()
	if *inputImagePath == "" {
		return "", *side, errors.New("[Error] Path to original image was not found")
	}

	return *inputImagePath, *side, nil
}

func getNewRect(originalRect image.Rectangle, side int) (width int, height int, rect image.Rectangle) {
	width = (originalRect.Dx() / side) * side
	height = (originalRect.Dy() / side) * side

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	rect = image.Rectangle{upLeft, lowRight}
	return
}

func getColors(originalImage image.Image, width int, height int, side int) (colors []color.RGBA) {
	const NumToDiv = 257

	for sec := 0; sec < side*side; sec++ {
		secWidth := (width / side) * (sec/side + 1)
		secHeight := (height / side) * (sec%side + 1)
		var sumR, sumG, sumB, sumA uint32
		for x := 0; x < secWidth; x++ {
			for y := 0; y < secHeight; y++ {
				r, g, b, a := originalImage.At(x, y).RGBA()
				sumR += r / NumToDiv
				sumG += g / NumToDiv
				sumB += b / NumToDiv
				sumA += a / NumToDiv
			}
		}
		numOfElem := uint32(secWidth * secHeight)
		aveR := sumR / numOfElem
		aveG := sumG / numOfElem
		aveB := sumB / numOfElem
		aveA := sumA / numOfElem

		colors = append(colors, color.RGBA{uint8(aveR), uint8(aveG), uint8(aveB), uint8(aveA)})
	}
	return
}

func createNewImage(colors []color.RGBA, side int, width int, height int, rect image.Rectangle) (newImage *image.RGBA) {
	newImage = image.NewRGBA(rect)

	for sec := 0; sec < side*side; sec++ {
		secStartWidth := (width / side) * (sec / side)
		secEndWidth := (width / side) * (sec/side + 1)
		secStartHeight := (height / side) * (sec % side)
		secEndHeight := (height / side) * (sec%side + 1)

		col := colors[sec]
		for x := secStartWidth; x < secEndWidth; x++ {
			for y := secStartHeight; y < secEndHeight; y++ {
				newImage.Set(x, y, col)
			}
		}
	}
	return
}

func getNewImageName(originalName string, side int) string {
	nameS := strings.Split(originalName, "/")
	name := nameS[len(nameS)-1]
	ext := path.Ext(name)

	return strconv.Itoa(side*side) + "tile_" + name[0:len(name)-len(ext)] + ".png"
}
