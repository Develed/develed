package main

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/develed/develed/imconv"
	srv "github.com/develed/develed/services"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"

	tm "github.com/buger/goterm"
)

// TermSink redirects any image written to it to the terminal's stdout.
// It requires a Truecolor-capable terminal in order to render images correctly.
type TermSink struct{}

// NewTermSink creates a new TermSink. It returns an error if the sdout of the
// calling process is not a terminal.
func NewTermSink() (*TermSink, error) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return nil, errors.New("TermSink requires stdout to be a terminal")
	}
	return &TermSink{}, nil
}

func (ts *TermSink) Draw(ctx context.Context, req *srv.DrawRequest) (*srv.DrawResponse, error) {
	img, _, err := image.Decode(bytes.NewReader(req.Data))
	if err != nil {
		return nil, err
	}

	sz := img.Bounds().Size()
	tm.Clear()
	tm.MoveCursor(0, 0)
	for y := 0; y < sz.Y; y++ {
		for x := 0; x < sz.X; x++ {
			col := imconv.NormalizeColor(img.At(x, y))
			tm.Printf("\033[48;2;%d;%d;%dm ", col.R, col.G, col.B)
		}
		tm.Print("\033[0m\n")
	}
	tm.Flush()
	return &srv.DrawResponse{Code: 0, Status: "OK"}, nil
}
