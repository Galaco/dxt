package dxt

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestNewDxt5(t *testing.T) {
	t.Skip()
	r := image.Rect(0, 0, 256, 256)

	img := NewDxt5(r)

	file, err := os.Open("../dxtpreview/test-dxt5.dds")
	if err != nil {
		t.Error(err)
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(file)
	if err != nil {
		t.Error(err)
	}

	err = img.Decompress(buf.Bytes(), false)
	if err != nil {
		t.Error(err)
	}
	out, _ := os.Create("out.png")
	err = png.Encode(out, img)
	if err != nil {
		t.Error(err)
	}
}
