package imconv

import (
	"bytes"
	"image"
	"image/color"
)

const (
	cBit1 = 0xf0
	cBit0 = 0x80
)

// FromImage converts an Image to a byte representation suitable for Develed.
func FromImage(img image.Image) []byte {
	sz := img.Bounds().Size()
	bb := bytes.NewBuffer(make([]byte, 0, sz.X*sz.Y*24))

	for y, invert := 0, false; y < sz.Y; y, invert = y+1, !invert {
		for x := 0; x < sz.X; x++ {
			var col color.Color
			if invert {
				col = img.At(sz.X-x-1, y)
			} else {
				col = img.At(x, y)
			}
			bb.Write(colorToPixelData(col))
		}
	}
	return bb.Bytes()
}

func colorToPixelData(c color.Color) []byte {
	var bb bytes.Buffer

	convert := func(c uint8) []byte {
		data := make([]byte, 8)
		for i := 0; i < 8; i++ {
			if ((c >> uint(i)) & 0x01) == 1 {
				data[7-i] = cBit1
			} else {
				data[7-i] = cBit0
			}
		}
		return data
	}

	conv := NormalizeColor(c)

	bb.Write(convert(uint8(conv.G)))
	bb.Write(convert(uint8(conv.R)))
	bb.Write(convert(uint8(conv.B)))

	return bb.Bytes()
}

func NormalizeColor(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()

	// For some reason, Color.RGBA duplicates the lower 8 bits and shifts them by 8
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
