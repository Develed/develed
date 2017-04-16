package main

import (
	"bytes"
	_ "image/png"
	"reflect"
	"testing"
)

func TestBlitImage(t *testing.T) {
	m, err := testImage()
	if err != nil {
		t.Fatal(err)
	}

	dest := &bytes.Buffer{}
	if err := blitImage(m, dest); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dest.Bytes(), rawdata) {
		t.Fatal("mismatched data")
	}
}
