package main

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/png"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/develed/develed/config"
	srv "github.com/develed/develed/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	debug = flag.Bool("debug", false, "enter debug mode")
	cfg   = flag.String("config", "/etc/develed.toml", "configuration file")
)

type server struct {
	sink srv.ImageSinkClient
}

func (s *server) Write(ctx context.Context, req *srv.TextRequest) (*srv.TextResponse, error) {
	var font FontMgr

	fontImage := font.Init(req.Font)

	// Allocate frame
	img := image.NewRGBA(image.Rect(0, 0, 39, 9))
	col := color.RGBA{0, 0, 0, 255}
	nm := img.Bounds()
	for y := 0; y < nm.Dy(); y++ {
		for x := 0; x < nm.Dx(); x++ {
			img.Set(x, y, col)
		}
	}

	// Fill frame
	for n, key := range req.Text {
		outx := n * font.Width()
		wf := font.Width()
		hf := font.High()
		col := font.Col(key)
		row := font.Row(key)

		for y := 0; y < hf; y++ {
			for x := 0; x < wf; x++ {
				img.Set(x+outx, y, fontImage.At(x+wf*col, y+hf*row))
			}
		}
	}

	buf := &bytes.Buffer{}
	png.Encode(buf, img)

	resp, err := s.sink.Draw(context.Background(), &srv.DrawRequest{
		Data: buf.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	return &srv.TextResponse{
		Code:   resp.Code,
		Status: resp.Status,
	}, nil
}

func main() {
	var err error

	flag.Parse()

	conf, err := config.Load(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	sock, err := net.Listen("tcp", conf.Textd.GRPCServerAddress)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := grpc.Dial(conf.DSPD.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	s := grpc.NewServer()
	srv.RegisterTextdServer(s, &server{
		sink: srv.NewImageSinkClient(conn),
	})
	reflection.Register(s)

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
