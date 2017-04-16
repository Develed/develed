package imconv

import (
	"bytes"
	"image/color"
	"reflect"
	"testing"
)

func TestColorToPixelData(t *testing.T) {
	tests := []struct {
		pixel color.Color
		data  []byte
	}{
		{
			pixel: color.RGBA{R: 255, G: 255, B: 255},
			data:  bytes.Repeat([]byte{cBit1}, 24),
		},
		{
			pixel: color.RGBA{R: 0, G: 0, B: 0},
			data:  bytes.Repeat([]byte{cBit0}, 24),
		},
		{
			pixel: color.RGBA{R: 255, G: 0, B: 0},
			data: append(bytes.Repeat([]byte{cBit0}, 16),
				bytes.Repeat([]byte{cBit1}, 8)...),
		},
		{
			pixel: color.RGBA{R: 0, G: 255, B: 0},
			data: append(bytes.Repeat([]byte{cBit1}, 8),
				bytes.Repeat([]byte{cBit0}, 16)...),
		},
		{
			pixel: color.RGBA{R: 0, G: 0, B: 255},
			data: append(append(bytes.Repeat([]byte{cBit0}, 8),
				bytes.Repeat([]byte{cBit1}, 8)...),
				bytes.Repeat([]byte{cBit0}, 8)...),
		},
		{
			pixel: color.RGBA{R: 113, G: 15, B: 15},
			data: []byte{
				cBit0, cBit0, cBit0, cBit0, cBit1, cBit1, cBit1, cBit1,
				cBit0, cBit0, cBit0, cBit0, cBit1, cBit1, cBit1, cBit1,
				cBit0, cBit1, cBit1, cBit1, cBit0, cBit0, cBit0, cBit1,
			},
		},
	}

	for _, test := range tests {
		data := colorToPixelData(test.pixel)
		if !reflect.DeepEqual(data, test.data) {
			t.Fatalf("Failed on %v with %v", test.pixel, data)
		}
	}
}
