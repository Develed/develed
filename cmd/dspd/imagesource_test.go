package main

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strings"
	"testing"
)

func TestGobImageSource(t *testing.T) {
	m, err := testImage()
	if err != nil {
		t.Fatal(err)
	}

	src := &bytes.Buffer{}
	gis := NewGobImageSource(src)

	enc := gob.NewEncoder(src)
	if err := enc.Encode(&m); err != nil {
		t.Fatal(err)
	}

	img, err := gis.Read()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(m, img) {
		t.Fatal("images do not match after encoding/decoding")
	}
}

func TestBase64ImageSource(t *testing.T) {
	m, err := testImage()
	if err != nil {
		t.Fatal(err)
	}

	dec := NewBase64ImageSource(strings.NewReader(cImageData))
	img, err := dec.Read()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(m, img) {
		t.Fatal("images do not match after encoding/decoding")
	}
}

func TestImageStream(t *testing.T) {
	m, err := testImage()
	if err != nil {
		t.Fatal(err)
	}

	testLen := 5
	dec := NewBase64ImageSource(strings.NewReader(strings.Repeat(cImageData, testLen)))
	for i := 0; i < testLen; i++ {
		img, err := dec.Read()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(m, img) {
			t.Fatal("images do not match after encoding/decoding")
		}
	}
}
