package main

import (
	"encoding/base64"
	"errors"
	"image"
	_ "image/png"
	"strings"
)

func testImage() (image.Image, error) {
	rd := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imgdata))
	m, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}

	if (m.Bounds().Size() != image.Point{X: 2, Y: 2}) {
		return nil, errors.New("mismatched image size")
	}

	return m, nil
}

const imgdata = `
iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAHWlU
WHRDb21tZW50AAAAAABDcmVhdGVkIHdpdGggR0lNUGQuZQcAAAAWSURBVAjXY5ZeF2bFt5Jlp9zHWz/c
ASMMBflQQRXXAAAAAElFTkSuQmCC
`

var rawdata = []byte{
	240, 128, 240, 128, 240, 240, 240, 128,
	128, 240, 128, 240, 128, 240, 240, 128,
	128, 128, 128, 240, 240, 128, 240, 240,

	128, 240, 240, 128, 128, 240, 128, 240,
	240, 240, 128, 240, 128, 240, 128, 128,
	128, 240, 128, 128, 128, 240, 240, 240,

	240, 240, 128, 128, 240, 240, 128, 128,
	128, 240, 128, 128, 128, 240, 240, 240,
	240, 240, 128, 240, 128, 240, 128, 128,

	128, 240, 128, 240, 240, 240, 128, 240,
	128, 128, 128, 240, 240, 128, 240, 240,
	240, 128, 240, 128, 240, 240, 240, 128,
}
