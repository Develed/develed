package main

import (
	"bytes"
	"encoding/gob"
	"reflect"
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
