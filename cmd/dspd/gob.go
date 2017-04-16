package main

import (
	"encoding/gob"
	"image"
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
