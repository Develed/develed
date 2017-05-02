package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
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

type RawImageSource struct {
	r io.Reader
}

func NewRawImageSource(r io.Reader) *RawImageSource {
	return &RawImageSource{
		r: r,
	}
}

func (ris *RawImageSource) Read() (image.Image, error) {
	var size uint64
	if err := binary.Read(ris.r, binary.LittleEndian, &size); err != nil {
		return nil, err
	}

	data := make([]byte, size)
	if _, err := io.ReadFull(ris.r, data); err != nil {
		return nil, err
	}

	m, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return m, nil
}
