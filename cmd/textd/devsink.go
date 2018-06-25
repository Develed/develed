package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/develed/develed/imconv"
	srv "github.com/develed/develed/services"
	"golang.org/x/net/context"
)

// Must be the same size as the driver's buffer
var blitBuff [16384]byte

// DeviceSink serializes images to a character device using the appropriate
// timings and parameters.
type DeviceSink struct {
	dev *os.File
}

// NewDeviceSink create a new DeviceSink using the device specified by devname.
// If no such device exists or the process lacks permissions to access it,
// an error is returned. Otherwise, the device is initialized with the correct
// parameters.
func NewDeviceSink(devname string) (*DeviceSink, error) {
	out, err := os.OpenFile(devname, os.O_WRONLY, os.ModeCharDevice)
	if err != nil {
		return nil, err
	}

	return &DeviceSink{
		dev: out,
	}, nil
}

func (ds *DeviceSink) Run() error {
	select {}
}

func (ds *DeviceSink) Draw(ctx context.Context, req *srv.DrawRequest) (*srv.DrawResponse, error) {
	img, _, err := image.Decode(bytes.NewReader(req.Data))
	if err != nil {
		return nil, err
	}

	if err := blitImage(img, ds.dev); err != nil {
		return nil, err
	}

	return &srv.DrawResponse{Code: 0, Status: "OK"}, nil
}

func blitImage(img image.Image, w io.Writer) error {
	copy(blitBuff[:], imconv.FromImage(img))

	// Write the buffer twice (required because the driver uses two buffers internally)
	if err := writeFull(w, blitBuff[:]); err != nil {
		return err
	}
	return writeFull(w, blitBuff[:])
}

func writeFull(w io.Writer, data []byte) (err error) {
	for total, last := 0, 0; total < len(data); total += last {
		fmt.Println(last)
		if last, err = w.Write(data[total:]); err != nil {
			return
		}
	}
	return
}
