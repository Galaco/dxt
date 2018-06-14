package common

import "image/color"

// Rgb565toargb8888
// Convert packed RGB565 color data to RGBA data.
// Note alpha is always initialized to 0xFF
func Rgb565toargb8888(packed uint16) color.RGBA {
	color := color.RGBA{}
	color.R = uint8((packed >> 11) & 0x1F)
	color.G = uint8((packed >> 5) & 0x3F)
	color.B = uint8((packed) & 0x1F)
	color.A = 0xFF

	color.R = (color.R << 3) | (color.R >> 2)
	color.G = (color.G << 2) | (color.G >> 4)
	color.B = (color.B << 3) | (color.B >> 2)

	return color
}