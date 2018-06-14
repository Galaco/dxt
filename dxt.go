package dxt

import (
	"image/color"
	"image"
)


const blockSize = 4

type Dxt interface {
	ColorModel() color.Model
	Bounds() image.Rectangle
	At(x, y int) color.Color
	PixOffset(x, y int) int
	Set(x, y int, c color.Color)
	Decompress(packed []byte) error
}