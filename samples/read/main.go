package main

import (
	"log"
	"image"
	"os"
	"bytes"
	"github.com/galaco/dxt"
	"image/png"
)


func main() {
	r := image.Rect(0, 0, 512, 512)

	img := dxt.NewDxt5(r)

	file,err := os.Open("test-dxt5.dds")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(file)

	// Why not start at 0?
	img.Decompress(buf.Bytes(), true)
	out, _ := os.Create("out.png")


	err = png.Encode(out, img)
}
