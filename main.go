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
	"sync"
)

func main() {
	inputImagePath, sideBlocks, err := getFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	originalFile, err := os.Open(inputImagePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer originalFile.Close()

	originalImage, _, err := image.Decode(originalFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	newImage := createNewImage(originalImage, sideBlocks)

	name := getNewImageName(originalFile.Name(), sideBlocks)
	if _, err := os.Stat("./outputs"); os.IsNotExist(err) {
		err := os.Mkdir("./outputs", os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Chmod("./outputs", 0777)
	}
	f, err := os.Create("./outputs/" + name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = png.Encode(f, newImage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Pixelated!")
	fmt.Println("Here is pixelated image: " + f.Name())
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

func getNewRect(originalRect image.Rectangle, sideBlocks int) (width int, height int, rect image.Rectangle) {
	width = (originalRect.Dx() / sideBlocks) * sideBlocks
	height = (originalRect.Dy() / sideBlocks) * sideBlocks

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	rect = image.Rectangle{Min: upLeft, Max: lowRight}
	return
}

type SectionInfo struct {
	sectionNumber int
	startWidth    int
	endWidth      int
	startHeight   int
	endHeight     int
}

func createNewImage(originalImage image.Image, sideBlocks int) *image.RGBA {
	width, height, rect := getNewRect(originalImage.Bounds(), sideBlocks)
	newImage := image.NewRGBA(rect)
	var wg sync.WaitGroup

	for section := 0; section < sideBlocks*sideBlocks; section++ {
		sectionInfo := SectionInfo{
			sectionNumber: section,
			startWidth:    (width / sideBlocks) * (section / sideBlocks),
			endWidth:      (width / sideBlocks) * (section/sideBlocks + 1),
			startHeight:   (height / sideBlocks) * (section % sideBlocks),
			endHeight:     (height / sideBlocks) * (section%sideBlocks + 1),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			createBlock(originalImage, newImage, sectionInfo)
		}()
	}
	wg.Wait()

	return newImage
}

func get8bitColor(c uint32) uint32 {
	return c / 257
}

func createBlock(originalImage image.Image, newImage *image.RGBA, sectionInfo SectionInfo) *image.RGBA {
	var sumR, sumG, sumB, sumA, cnt uint32
	for x := sectionInfo.startWidth; x < sectionInfo.endWidth; x++ {
		for y := sectionInfo.startHeight; y < sectionInfo.endHeight; y++ {
			r, g, b, a := originalImage.At(x, y).RGBA()
			sumR += get8bitColor(r)
			sumG += get8bitColor(g)
			sumB += get8bitColor(b)
			sumA += get8bitColor(a)
			cnt++
		}
	}
	aveR := uint8(sumR / cnt)
	aveG := uint8(sumG / cnt)
	aveB := uint8(sumB / cnt)
	aveA := uint8(sumA / cnt)

	for x := sectionInfo.startWidth; x < sectionInfo.endWidth; x++ {
		for y := sectionInfo.startHeight; y < sectionInfo.endHeight; y++ {
			newImage.Set(x, y, color.RGBA{R: aveR, G: aveG, B: aveB, A: aveA})
		}
	}

	return newImage
}

func getNewImageName(originalName string, side int) string {
	nameS := strings.Split(originalName, "/")
	name := nameS[len(nameS)-1]
	ext := path.Ext(name)

	return strconv.Itoa(side*side) + "blocks_" + name[0:len(name)-len(ext)] + ".png"
}
