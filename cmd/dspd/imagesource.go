package main

import (
	"bufio"
	"encoding/base64"
	"encoding/gob"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type ImageSource interface {
	Read() (image.Image, error)
}

type GobImageSource struct {
	dec *gob.Decoder
}

func NewGobImageSource(r io.Reader) *GobImageSource {
	gob.Register(&image.RGBA{})

	return &GobImageSource{
		dec: gob.NewDecoder(r),
	}
}

func (gis *GobImageSource) Read() (image.Image, error) {
	var img image.Image
	if err := gis.dec.Decode(&img); err != nil {
		return nil, err
	}
	return img, nil
}

type Base64ImageSource struct {
	dec *bufio.Reader
}

func NewBase64ImageSource(r io.Reader) *Base64ImageSource {
	return &Base64ImageSource{
		dec: bufio.NewReader(base64.NewDecoder(base64.StdEncoding, r)),
	}
}

func (bis *Base64ImageSource) Read() (image.Image, error) {
	m, _, err := image.Decode(bis.dec)
	if err != nil {
		return nil, err
	}
	return m, nil
}
