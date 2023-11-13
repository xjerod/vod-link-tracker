package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"time"

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

var NonAlphanumericRegex = regexp.MustCompile(`^[^a-zA-Z0-9 ]+`)

type Link struct {
	Word     string
	Filename string
}

type Post struct {
	Name string
	Date time.Time
}

type Content struct {
	Name  string
	Date  time.Time
	Links []string
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func main() {
	skipCrop := flag.Bool("nocrop", false, "skip the crop step")
	var vodID string
	flag.StringVar(&vodID, "vodid", "", "the id of the vod, looks for the files in data/frames/<id>/*")

	flag.Parse()

	client := gosseract.NewClient()
	defer client.Close()

	rawDirPath := fmt.Sprintf("%s/%s", RawDir, vodID)
	fmt.Println("Raw dir: ", rawDirPath)

	croppedDirPath := fmt.Sprintf("%s/%s", CroppedDir, vodID)
	fmt.Println("Cropped dir: ", croppedDirPath)

	entries, err := os.ReadDir(rawDirPath)
	check(err)

	files := []string{}
	for _, file := range entries {
		// TODO: quick hack
		if !file.IsDir() && file.Name() != ".gitkeep" {
			files = append(files, file.Name())
		}
	}

	cropInput := make(chan string, len(files)+1000)
	var wg sync.WaitGroup
	if skipCrop != nil && !*skipCrop {
		fmt.Printf("Files to crop %v\n", files)
		for i := 0; i < 200; i++ {
			wg.Add(1)
			go worker(cropInput, &wg, cropImage)
		}

		for _, name := range files {
			cropInput <- name
		}
	}

	close(cropInput)

	wg.Wait()

	linkMap := map[string][]string{}
	for _, file := range files {
		fmt.Println("Parsing ", file)
		err := client.SetImage(fmt.Sprintf("%s/%s", croppedDirPath, file))
		check(err)

		boxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
		check(err)
		for _, box := range boxes {
			if ClampByConfidence && box.Confidence < ClampPoint {
				continue
			}

			word := NonAlphanumericRegex.ReplaceAllString(box.Word, "")
			fmt.Println("DEBUG", file, box.Confidence, word)
			linkMap[word] = append(linkMap[word], file)
		}
	}

	tmplFile := "./azanlinks/archetypes/links.md"
	tmplName := path.Base(tmplFile)
	tmpl := template.Must(template.New(tmplName).ParseFiles(tmplFile))

	dateString := "2023-10-29"
	now, _ := time.Parse("2006-01-02", dateString)
	links := Content{
		Name:  "UAW VICTORY!!! DAGESTAN ANTISEMITISM! GAZA WARCRIMES CONTINUE.",
		Date:  now,
		Links: []string{},
	}

	for name := range linkMap {
		links.Links = append(links.Links, name)
	}
	fmt.Println("LINKS:")
	fmt.Println(links)

	file, err := os.Create(fmt.Sprintf("./azanlinks/content/posts/%s.md", dateString))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, links)
	if err != nil {
		panic(err)
	}
}

func worker(input chan string, wg *sync.WaitGroup, fn func(string)) {
	defer wg.Done()

	for file := range input {
		fn(file)
	}
}

func cropImage(name string) {
	file, err := os.Open(fmt.Sprintf("%s/%s/%s", RawDir, "1964280059", name))
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

	outpath := fmt.Sprintf("%s/%s/%s", CroppedDir, "1964280059", name)
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
