package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	srv "github.com/develed/develed/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	debug = flag.Bool("debug", false, "output to screen instead of /dev/dsp")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [opts...] HOST:PORT\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	var sink srv.ImageSinkServer
	var err error

	if flag.Parse(); flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
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

	sock, err := net.Listen("tcp", flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	srv.RegisterImageSinkServer(s, sink)
	reflection.Register(s)

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
