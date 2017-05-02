package imconv

import (
	"bytes"
	"image"
	"image/color"

	"github.com/Sirupsen/logrus"
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

	r, g, b, _ := c.RGBA()

	// For some reason, Color.RGBA duplicates the lower 8 bits and shifts them by 8
	if r, g, b = r>>8, g>>8, b>>8; r > 255 || g > 255 || b > 255 {
		logrus.Warnln("Some pixel color values are longer than 8 bits:", r, g, b)
	}

	bb.Write(convert(uint8(g)))
	bb.Write(convert(uint8(r)))
	bb.Write(convert(uint8(b)))

	return bb.Bytes()
}
