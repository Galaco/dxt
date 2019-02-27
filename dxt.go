package dxt

import (
	"image"
	"image/color"
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

type Header struct {
	Id                uint32
	Size              uint32
	Flags             uint32
	Height            uint32
	Width             uint32
	PitchOrLinearSize uint32
	Depth             uint32
	MipMapCount       uint32
	Reserved1         [11]uint32
	Spf               struct {
		Size        uint32
		Flags       uint32
		FourCC      uint32
		RGBBitCount uint32
		RBitMask    uint32
		GBitMask    uint32
		BBitMask    uint32
		ABitMask    uint32
	}
	Caps      uint32
	Caps2     uint32
	Caps3     uint32
	Caps4     uint32
	Reserved2 uint32
}
