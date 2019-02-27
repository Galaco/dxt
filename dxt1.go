package dxt

import (
	"bytes"
	"encoding/binary"
	"github.com/galaco/dxt/common"
	"image"
	"image/color"
)

type Dxt1 struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func (p *Dxt1) ColorModel() color.Model { return color.RGBAModel }

func (p *Dxt1) Bounds() image.Rectangle { return p.Rect }

func (p *Dxt1) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *Dxt1) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	return color.RGBA{R: p.Pix[i+0], G: p.Pix[i+1], B: p.Pix[i+2], A: p.Pix[i+3]}
}

func (p *Dxt1) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

func (p *Dxt1) Set(x, y int, c color.Color) {
	if !(image.Point{X: x, Y: y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	p.Pix[i+0] = c1.R
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.B
	p.Pix[i+3] = c1.A
}

func (p *Dxt1) Decompress(packed []byte) error {
	argb, err := decompressDxt1(packed, p.Rect.Dx(), p.Rect.Dy())
	if err != nil {
		return err
	}

	for i, c := range argb {
		i *= 4
		p.Pix[i] = c.R
		p.Pix[i+1] = c.G
		p.Pix[i+2] = c.B
		p.Pix[i+3] = c.A
	}

	return nil
}

func NewDxt1(r image.Rectangle) *Dxt1 {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 4*w*h)
	return &Dxt1{buf, 4 * w, r}
}

func decompressDxt1(packed []byte, width int, height int) ([]color.RGBA, error) {
	unpacked := make([]color.RGBA, width*height)

	blockCountX := (width + 3) / 4
	blockCountY := (height + 3) / 4

	offset := 0
	for j := 0; j < blockCountY; j++ {
		for i := 0; i < blockCountX; i++ {
			err := decompressDxt1Block(packed[offset+(i*8):], i*4, j*4, width, unpacked)
			if err != nil {
				return nil,err
			}
		}
		offset += blockCountX * 8
	}

	return unpacked, nil
}

func decompressDxt1Block(packed []byte, offsetX int, offsetY int, width int, unpacked []color.RGBA) error {
	// Construct colours to transform between
	var c0, c1 uint16
	err := binary.Read(bytes.NewBuffer(packed[:2]), binary.LittleEndian, &c0)
	if err != nil {
		return err
	}
	err = binary.Read(bytes.NewBuffer(packed[2:4]), binary.LittleEndian, &c1)
	if err != nil {
		return err
	}

	colour0 := common.Rgb565toargb8888(c0)
	colour1 := common.Rgb565toargb8888(c1)

	var code uint32
	err = binary.Read(bytes.NewBuffer(packed[4:]), binary.LittleEndian, &code)
	if err != nil {
		return err
	}

	// iterate through pixels
	for j := 0; j < blockSize; j++ {
		for i := 0; i < blockSize; i++ {
			var finalColour color.RGBA
			positionCode := (code >> uint32(2*(4*j+i))) & 0x03

			if c0 > c1 {
				switch positionCode {
				case 0:
					finalColour = colour0
				case 1:
					finalColour = colour1
				case 2:
					finalColour = color.RGBA{
						R: (2*colour0.R + colour1.R) / 3,
						G: (2*colour0.G + colour1.G) / 3,
						B: (2*colour0.B + colour1.B) / 3,
					}
				case 3:
					finalColour = color.RGBA{
						R: (colour0.R + 2*colour1.R) / 3,
						G: (colour0.G + 2*colour1.G) / 3,
						B: (colour0.B + 2*colour1.B) / 3,
					}
				}
			} else {
				switch positionCode {
				case 0:
					finalColour = colour0
				case 1:
					finalColour = colour1
				case 2:
					finalColour = color.RGBA{
						R: (colour0.R + colour1.R) / 2,
						G: (colour0.G + colour1.G) / 2,
						B: (colour0.B + colour1.B) / 2,
					}
				case 3:
					finalColour = color.RGBA{
						R: 0,
						G: 0,
						B: 0,
					}
				}
			}
			// Ensure no alpha
			finalColour.A = 0xFF

			if offsetX+i < width {
				unpacked[(offsetY+j)*width+(offsetX+i)] = finalColour
			}
		}
	}

	return nil
}
