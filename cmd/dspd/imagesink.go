package main

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"syscall"
	"unsafe"

	"github.com/develed/develed/imconv"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	cResetDuration = 60 // [us]
	cBytePerUSec   = 5
)

var resetCmd [cResetDuration * cBytePerUSec]byte

// ImageSink describes the output stage of the rendering pipeline.
type ImageSink interface {
	Write(img image.Image) error
	Close() error
}

// DummySink is a /dev/null-style ImageSink, useful for debugging purposes.
type DummySink struct{}

// Write does nothing and always returns without errors.
func (ds *DummySink) Write(img image.Image) error {
	return nil
}

// Close does nothing and always returns without errors.
func (ds *DummySink) Close() error {
	return nil
}

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
	sampleFormatIoctl := 0xc0045005
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

func writeFull(w io.Writer, data []byte) (err error) {
	for total, last := 0, 0; total < len(data); total += last {
		if last, err = w.Write(data[total:]); err != nil {
			return
		}
	}
	return
}

func sendResetCmd(w io.Writer) error {
	return writeFull(w, resetCmd[:])
}

func blitImage(img image.Image, w io.Writer) error {
	return writeFull(w, imconv.FromImage(img))
}

// Write serializes img using the appropriate parameters and sends the
// resulting byte stream to the underlying device.
// It also issues a reset sequence to correctly show the image on screen.
func (ds *DeviceSink) Write(img image.Image) error {
	if err := blitImage(img, ds.dev); err != nil {
		return err
	}
	if err := sendResetCmd(ds.dev); err != nil {
		return err
	}
	return nil
}

// Close closes the underlying device. Any Write called after a Close will
// return an error.
func (ds *DeviceSink) Close() error {
	return ds.dev.Close()
}

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

// Write renders img to stdout by using Truecolor ANSI escape codes.
// If the terminal does not support Truecolor, the image won't be rendered
// correctly.
func (ts *TermSink) Write(img image.Image) error {
	sz := img.Bounds().Size()
	for y := 0; y < sz.Y; y++ {
		for x := 0; x < sz.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if _, err := fmt.Printf("\033[48;2;%d;%d;%dm ", r, g, b); err != nil {
				return err
			}
		}
		if _, err := fmt.Print("\033[0m\n"); err != nil {
			return err
		}
	}
	return nil
}

// Close does nothing and always returns without errors.
func (ts *TermSink) Close() error {
	return nil
}
