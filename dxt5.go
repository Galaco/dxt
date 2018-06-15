package dxt

import (
	"image/color"
	"image"
	"encoding/binary"
	"bytes"
	"github.com/galaco/dxt/common"
	"errors"
)

// Dxt5
// Dxt5 Image fulfills the standard golang image interface.
// It also fulfils a slightly more specialised Dxt interface in the package.
type Dxt5 struct {
	Header Header
	Pix []uint8
	Stride int
	Rect image.Rectangle
}

// ColorModel
// Returns the color Model for the image (always RGBA)
// Fulfills the requirements for the image interface
func (p *Dxt5) ColorModel() color.Model { return color.RGBAModel }

// Bounds
// Returns image boundaries
// Fulfills the requirements for the image interface
func (p *Dxt5) Bounds() image.Rectangle { return p.Rect }

// At
// Returns generic Color data for a single pixel at location x,y
// Fulfills the requirements for the image interface
func (p *Dxt5) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

// RGBAAt
// Returns colour.RGBA information for a single pixel at location x,y
// Fulfills the requirements for the image interface
func (p *Dxt5) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	return color.RGBA{p.Pix[i+0], p.Pix[i+1], p.Pix[i+2], p.Pix[i+3]}
}

// PixOffset
// Returns the offset into image data of an x,y coordinate
// Fulfills the requirements for the image interface
func (p *Dxt5) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

// Set
// Set the Color at a given x,y coordinate
// Fulfills the requirements for the image interface
func (p *Dxt5) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	p.Pix[i+0] = c1.R
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.B
	p.Pix[i+3] = c1.A
}

// Decompress
// Decompresses and populates the image from packed dxt5 data
func (p *Dxt5) Decompress(packed []byte, withHeader bool) error {
	var rgba []color.RGBA
	var err error
	if withHeader == true {
		var header Header
		binary.Read(bytes.NewBuffer(packed[:128]), binary.LittleEndian, &header)
		if header.Id != 0x20534444 {
			return errors.New("dds format identifier mismatch")
		}
		p.Header = header
		rgba,err = decompressDxt5(packed[128:], p.Rect.Dx(), p.Rect.Dy())
	} else {
		rgba,err = decompressDxt5(packed, p.Rect.Dx(), p.Rect.Dy())
	}

	if err != nil {
		return err
	}
	for i,c := range rgba {
		i *= 4
		p.Pix[i] = c.R
		p.Pix[i+1] = c.G
		p.Pix[i+2] = c.B
		p.Pix[i+3] = c.A
	}

	return nil
}

// NewDxt5
// Create a new Dxt5 image
func NewDxt5(r image.Rectangle) *Dxt5 {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 4*w*h)
	return &Dxt5{
		Header: Header{},
		Pix: buf,
		Stride: 4 * w,
		Rect: r,
	}
}


// decompressDxt5
// Decompress a Dxt5 compressed slice of bytes.
// Decompresses block by block
// Width and Height are required, as this information is impossible to derive with
// 100% accuracy (e.g. 256x1024 cannot be distinguished from 512x512) from raw alone
func decompressDxt5(packed []byte, width int, height int) ([]color.RGBA,error) {
	unpacked := make([]color.RGBA, width * height)

	blockCountX := int((width + 3) / blockSize)
	blockCountY := int((height + 3) / blockSize)

	offset := 0
	for j := 0; j < blockCountY; j++ {
		for i := 0; i < blockCountX; i++ {
			decompressDxt5Block(packed[offset+(i * 16):], i * blockSize, j * blockSize, width, unpacked)
		}
		offset += blockCountX * 16
	}

	return unpacked,nil
}

// decompressDxt5Block
// decompress a single dxt5 compressed block.
// A single decompressed block is 4x4 pixels located at x,y location in the resultant image
func decompressDxt5Block(packed []byte, offsetX int, offsetY int, width int, unpacked []color.RGBA) {
	var alpha0, alpha1 uint8
	binary.Read(bytes.NewBuffer(packed[:1]), binary.LittleEndian, &alpha0)
	binary.Read(bytes.NewBuffer(packed[1:2]), binary.LittleEndian, &alpha1)

	var bits [6]uint8
	binary.Read(bytes.NewBuffer(packed[2:8]), binary.LittleEndian, &bits)

	alphaCode1 := uint32(bits[2] | (bits[3] << 8) | (bits[4] << 16) | (bits[5] << 24))
	alphaCode2 := uint16(bits[0] | (bits[1] << 8))

	// Construct colours to transform between
	var c0, c1 uint16
	binary.Read(bytes.NewBuffer(packed[8:10]), binary.LittleEndian, &c0)
	binary.Read(bytes.NewBuffer(packed[10:12]), binary.LittleEndian, &c1)

	colour0 := common.Rgb565toargb8888(c0)
	colour1 := common.Rgb565toargb8888(c1)

	var code uint32
	binary.Read(bytes.NewBuffer(packed[12:16]), binary.LittleEndian, &code)

	for j := 0; j < blockSize; j++ {
		for i := 0; i < blockSize; i++ {
			alphaCodeIndex := uint(3*(4*j+i))
			var alphaCode int

			if alphaCodeIndex <= 12 {
				alphaCode = int((alphaCode2 >> alphaCodeIndex) & 0x07)
			} else if alphaCodeIndex == 15 {
				alphaCode = int((uint32(alphaCode2) >> 15) | ((alphaCode1 << 1) & 0x06))
			} else {
				// alphaCodeIndex >= 18 && alphaCodeIndex <= 45
				alphaCode = int((alphaCode1 >> (alphaCodeIndex - 16)) & 0x07)
			}

			var finalAlpha uint8
			if alphaCode == 0 {
				finalAlpha = alpha0
			} else if alphaCode == 1 {
				finalAlpha = alpha1
			} else {
				if alpha0 > alpha1 {
					finalAlpha = ((8-uint8(alphaCode))*alpha0 + (uint8(alphaCode)-1)*alpha1)/7
				} else {
					if alphaCode == 6 {
						finalAlpha = 0
					} else if alphaCode == 7 {
						finalAlpha = 255
					} else {
						finalAlpha = ((6-uint8(alphaCode))*alpha0 + (uint8(alphaCode)-1)*alpha1)/5
					}
				}
			}

			colorCode := (code >> uint32(2*(4*j+i))) & 0x03

			var finalColour color.RGBA
			switch colorCode {
			case 0:
				finalColour = colour0
			case 1:
				finalColour = colour1
			case 2:
				finalColour = color.RGBA{
					R: (2*colour0.R+colour1.R)/3,
					G: (2*colour0.G+colour1.G)/3,
					B: (2*colour0.B+colour1.B)/3,
				}
			case 3:
				finalColour = color.RGBA{
					R: (colour0.R+2*colour1.R)/3,
					G: (colour0.G+2*colour1.G)/3,
					B: (colour0.B+2*colour1.B)/3,
				}
			}

			if finalAlpha != 255 {
				a := 0
				a -=2
			}

			// Set alpha
			finalColour.A = finalAlpha

			if offsetX + i < width {
				unpacked[(offsetY + j) * width + (offsetX + i)] = finalColour
			}
		}
	}
}