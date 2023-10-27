// #cgo LDFLAGS: -Wl,-ld_classic

package main

import (
	"log"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()
	if err := client.SetImage("./data/frames/cropped/0001.png"); err != nil {
		log.Fatal(err)
	}

	boxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		log.Fatal(err)
	}

	for _, box := range boxes {
		log.Println(box.Confidence, box.Word)
	}
}
