package main

import (
	"fmt"
	"image"
	"io"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/plorefice/develed/imconv"
)

func blitImage(img image.Image, w io.Writer) error {
	var err error

	data := imconv.FromImage(img)

	for total, last := 0, 0; total < len(data); total += last {
		if last, err = w.Write(data[total:]); err != nil {
			return err
		}
	}
	return nil
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

	dsp, err := os.OpenFile("prova.dat", os.O_WRONLY, os.ModeCharDevice)
	if err != nil {
		log.Fatalln(err)
	}

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
