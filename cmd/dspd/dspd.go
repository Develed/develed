package main

import (
	"flag"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/develed/develed/config"
	srv "github.com/develed/develed/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	debug = flag.Bool("debug", false, "output to screen instead of /dev/dsp")
	cfg   = flag.String("config", "/etc/develed.toml", "configuration file")
)

type ImageSink interface {
	srv.ImageSinkServer
	Run() error
}

func main() {
	var sink ImageSink
	var err error

	flag.Parse()

	conf, err := config.Load(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	if !*debug {
		sink, err = NewDeviceSink("/dev/sscdev0")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		sink, err = NewTermSink()
		if err != nil {
			log.Fatalln(err)
		}
	}

	sock, err := net.Listen("tcp", conf.DSPD.GRPCServerAddress)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	srv.RegisterImageSinkServer(s, sink)
	reflection.Register(s)

	go func() {
		if err := s.Serve(sock); err != nil {
			log.Fatalln(err)
		}
	}()

	if err := sink.Run(); err != nil {
		log.Fatalln(err)
	}
}
