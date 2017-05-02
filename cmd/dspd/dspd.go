package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
)

var (
	debug = flag.Bool("debug", false, "output to screen instead of /dev/dsp")
)

func main() {
	var sink ImageSink
	var err error

	if flag.Parse(); flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s FIFO\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	// Actually read-only, write flag required to avoid blocking on open()
	fifo, err := os.OpenFile(flag.Arg(0), os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatalln(err)
	}
	defer fifo.Close()

	if !*debug {
		sink, err = NewDeviceSink("/dev/dsp")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		sink, err = NewTermSink()
		if err != nil {
			log.Fatalln(err)
		}
	}

	src := NewRawImageSource(fifo)
	for {
		img, err := src.Read()
		if err != nil {
			log.Errorln(err)
			continue
		}

		if err := sink.Write(img); err != nil {
			log.Errorln(err)
		}
	}
}
