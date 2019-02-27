package common

import "image/color"

// Rgb565toargb8888
// Convert packed RGB565 color data to RGBA data.
// Note alpha is always initialized to 0xFF
func Rgb565toargb8888(packed uint16) color.RGBA {
	colour := color.RGBA{}

	//var temp uint32
	//temp = uint32(packed >> 11) * 255 + 16
	//color.R = uint8((temp/32 + temp)/32)
	//temp = (uint32(packed & 0x07E0) >> 5) * 255 + 32
	//color.G = uint8((temp/64 + temp)/64)
	//temp = uint32(packed & 0x001F) * 255 + 16
	//color.B = uint8((temp/32 + temp)/32)

	colour.R = uint8((packed >> 11) & 0x1F)
	colour.G = uint8((packed >> 5) & 0x3F)
	colour.B = uint8((packed) & 0x1F)
	colour.A = 0xFF

	colour.R = (colour.R << 3) | (colour.R >> 2)
	colour.G = (colour.G << 2) | (colour.G >> 4)
	colour.B = (colour.B << 3) | (colour.B >> 2)

	return colour
}
