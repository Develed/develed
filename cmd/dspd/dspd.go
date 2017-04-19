package main

import (
	"fmt"
	"image"
	"io"
	"os"
	"path"
	"syscall"
	"unsafe"

	log "github.com/Sirupsen/logrus"
	"github.com/plorefice/develed/imconv"
)

const (
	cResetDuration = 60 // [us]
	cBytePerUSec   = 5

	cSampleFormatIoctl = 0xc0045005
)

var resetCmd [cResetDuration * cBytePerUSec]byte

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
	if err := writeFull(w, imconv.FromImage(img)); err != nil {
		return err
	}
	return sendResetCmd(w)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s FIFO\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	// Actually read-only, write flag required to avoid blocking on open()
	fifo, err := os.OpenFile(os.Args[1], os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatalln(err)
	}
	defer fifo.Close()

	dsp, err := os.OpenFile("/dev/dsp", os.O_WRONLY, os.ModeCharDevice)
	if err != nil {
		log.Fatalln(err)
	}
	defer dsp.Close()

	// configure /dev/dsp sample format
	sampleSize := 0x00002000 /* AFMT_S32_BE */
	syscall.Syscall(
		syscall.SYS_IOCTL,
		dsp.Fd(),
		cSampleFormatIoctl,
		uintptr(unsafe.Pointer(&sampleSize)))

	src := NewGobImageSource(fifo)
	for {
		img, err := src.Read()
		if err != nil {
			log.Errorln(err)
			continue
		}

		if err := blitImage(img, dsp); err != nil {
			log.Errorln(err)
		}
	}
}
