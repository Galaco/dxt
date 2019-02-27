package main

import (
	"bytes"
	"github.com/galaco/dxt"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	r := image.Rect(0, 0, 512, 512)

	img := dxt.NewDxt5(r)

	file, err := os.Open("test-dxt5.dds")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}

	// Why not start at 0?
	err = img.Decompress(buf.Bytes(), true)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := os.Create("out.png")

	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}
