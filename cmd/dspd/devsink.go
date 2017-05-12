package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"syscall"
	"unsafe"

	"github.com/develed/develed/imconv"
	srv "github.com/develed/develed/services"
	"golang.org/x/net/context"
)

const (
	cResetDuration = 60 // [us]
	cBytePerUSec   = 5
)

var resetCmd [cResetDuration * cBytePerUSec]byte

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

	// configure sample format
	sampleFormatIoctl := uint32(0xc0045005)
	sampleSize := 0x00002000 /* AFMT_S32_BE */
	syscall.Syscall(
		syscall.SYS_IOCTL,
		out.Fd(),
		uintptr(sampleFormatIoctl),
		uintptr(unsafe.Pointer(&sampleSize)))

	return &DeviceSink{
		dev: out,
	}, nil
}

func (ds *DeviceSink) Draw(ctx context.Context, req *srv.DrawRequest) (*srv.DrawResponse, error) {
	img, _, err := image.Decode(bytes.NewReader(req.Data))
	if err != nil {
		return nil, err
	}

	if err := blitImage(img, ds.dev); err != nil {
		return nil, err
	}
	if err := sendResetCmd(ds.dev); err != nil {
		return nil, err
	}

	return &srv.DrawResponse{Code: 0, Status: "OK"}, nil
}

func sendResetCmd(w io.Writer) error {
	return writeFull(w, resetCmd[:])
}

func blitImage(img image.Image, w io.Writer) error {
	return writeFull(w, imconv.FromImage(img))
}

func writeFull(w io.Writer, data []byte) (err error) {
	for total, last := 0, 0; total < len(data); total += last {
		if last, err = w.Write(data[total:]); err != nil {
			return
		}
	}
	return
}
