package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
)

var (
	debug = flag.Bool("debug", false, "output to screen instead of /dev/dsp")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [opts...] [PIPE]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	var sink ImageSink
	var in io.Reader
	var err error

	if flag.Parse(); flag.NArg() > 0 {
		fifo, err := os.OpenFile(flag.Arg(0), os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatalln(err)
		}
		defer fifo.Close()
		in = fifo
	} else {
		in = os.Stdin
	}

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

	src := NewRawImageSource(in)

loop:
	for {
		img, err := src.Read()
		if err != nil {
			if err == io.EOF {
				break loop
			} else {
				log.Errorln(err)
				continue
			}
		}

		if err := sink.Write(img); err != nil {
			log.Errorln(err)
		}
	}
}
