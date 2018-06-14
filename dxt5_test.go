package dxt

import (
	"testing"
	"image"
	"os"
	"bytes"
	"image/png"
)

func TestNewDxt5(t *testing.T) {
	r := image.Rect(0, 0, 256, 256)

	img := NewDxt5(r)

	file,err := os.Open("../dxtpreview/test-dxt5.dds")
	if err != nil {
		t.Error(err)
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(file)

	img.Decompress(buf.Bytes())
	out, _ := os.Create("out.png")
	png.Encode(out, img)
}
