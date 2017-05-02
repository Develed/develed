package main

import (
	"encoding/base64"
	"errors"
	"image"
	_ "image/png"
	"strings"
)

func testImage() (image.Image, error) {
	rd := base64.NewDecoder(base64.StdEncoding, strings.NewReader(cImageData))
	m, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}

	if (m.Bounds().Size() != image.Point{X: 2, Y: 2}) {
		return nil, errors.New("mismatched image size")
	}

	return m, nil
}

const cImageData = `
iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAHWlU
WHRDb21tZW50AAAAAABDcmVhdGVkIHdpdGggR0lNUGQuZQcAAAAWSURBVAjXY5ZeF2bFt5Jlp9zHWz/c
ASMMBflQQRXXAAAAAElFTkSuQmCC
`

var cRawImageData = []byte{
	137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 2, 0,
	0, 0, 2, 8, 2, 0, 0, 0, 253, 212, 154, 115, 0, 0, 0, 9, 112, 72, 89, 115, 0,
	0, 11, 19, 0, 0, 11, 19, 1, 0, 154, 156, 24, 0, 0, 0, 29, 105, 84, 88, 116,
	67, 111, 109, 109, 101, 110, 116, 0, 0, 0, 0, 0, 67, 114, 101, 97, 116, 101,
	100, 32, 119, 105, 116, 104, 32, 71, 73, 77, 80, 100, 46, 101, 7, 0, 0, 0,
	22, 73, 68, 65, 84, 8, 215, 99, 150, 94, 23, 102, 197, 183, 146, 101, 167,
	220, 199, 91, 63, 220, 1, 35, 12, 5, 249, 80, 65, 21, 215, 0, 0, 0, 0, 73,
	69, 78, 68, 174, 66, 96, 130,
}

var cRawData = []byte{
	240, 128, 240, 128, 240, 240, 240, 128,
	128, 128, 128, 240, 240, 128, 240, 240,
	128, 240, 128, 240, 128, 240, 240, 128,

	128, 240, 240, 128, 128, 240, 128, 240,
	128, 240, 128, 128, 128, 240, 240, 240,
	240, 240, 128, 240, 128, 240, 128, 128,

	128, 240, 128, 240, 240, 240, 128, 240,
	240, 128, 240, 128, 240, 240, 240, 128,
	128, 128, 128, 240, 240, 128, 240, 240,

	240, 240, 128, 128, 240, 240, 128, 128,
	240, 240, 128, 240, 128, 240, 128, 128,
	128, 240, 128, 128, 128, 240, 240, 240,
}
