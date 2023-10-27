package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"sync"

	"github.com/oliamb/cutter"
	"github.com/otiai10/gosseract/v2"
)

const (
	RawDir     = "./data/frames/raw"
	CroppedDir = "./data/frames/cropped"

	CropHeight      = 40
	CropWidthFull   = 900
	CropWidthNarrow = 600
	CropAnchorX     = 125
	CropAnchorY     = 40

	ClampByConfidence = true
	ClampPoint        = 45.0
)

type Link struct {
	Word     string
	Filename string
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	entries, err := os.ReadDir(RawDir)
	check(err)

	files := []string{}
	for _, file := range entries {
		// TODO: quick hack
		if !file.IsDir() && file.Name() != ".gitkeep" {
			files = append(files, file.Name())
		}
	}
	fmt.Printf("Files to crop %v\n", files)

	var wg sync.WaitGroup
	for _, name := range files {
		wg.Add(1)
		go cropImage(name, &wg)
	}

	wg.Wait()

	linkMap := map[string][]string{}
	for _, file := range files {
		fmt.Println("Parsing ", file)
		err := client.SetImage(fmt.Sprintf("%s/%s", CroppedDir, file))
		check(err)

		boxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
		check(err)
		for _, box := range boxes {
			if ClampByConfidence && box.Confidence < ClampPoint {
				continue
			}

			word := box.Word
			fmt.Println("DEBUG", file, box.Confidence, word)
			linkMap[word] = append(linkMap[word], file)
		}
	}

	fmt.Println("LINKS:")
	for name, files := range linkMap {
		fmt.Println(name, files)
	}
}

func cropImage(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fmt.Sprintf("%s/%s", RawDir, name))
	check(err)
	defer file.Close()

	img, _, err := image.Decode(file)
	check(err)

	croppedImage, err := cutter.Crop(img, cutter.Config{
		Height: CropHeight,
		Width:  CropWidthNarrow,
		Mode:   cutter.TopLeft,
		Anchor: image.Point{CropAnchorX, CropAnchorY},
	})
	check(err)

	outpath := fmt.Sprintf("%s/%s", CroppedDir, name)
	croppedFile, err := os.Create(outpath)
	check(err)
	defer croppedFile.Close()

	switch filepath.Ext(name) {
	case ".png":
		err = png.Encode(croppedFile, croppedImage)
	case ".jpg":
		err = jpeg.Encode(croppedFile, croppedImage, &jpeg.Options{})
	default:
		err = errors.New("Unsupported format: " + filepath.Ext(outpath))
	}
	check(err)
	fmt.Println("Image saved to", outpath)
}
