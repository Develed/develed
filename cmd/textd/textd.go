package main

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net"
	"time"

	bitmapfont "github.com/develed/develed/bitmapfont"

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

var conf *config.Global

type server struct {
	sink srv.ImageSinkClient
}

func (s *server) Write(ctx context.Context, req *srv.TextRequest) (*srv.TextResponse, error) {

	log.Debug(req.Font)
	err3 := bitmapfont.Init(conf.Textd.FontPath, req.Font, conf.BitmapFonts)
	if err3 != nil {
		log.Error(err3)
		return nil, err3
	}

	log.Debugf("Color: %v Bg: %v\n", req.FontColor, req.FontBg)

	txt_color := color.RGBA{255, 0, 0, 255}
	txt_bg := color.RGBA{0, 0, 0, 255}

	img, err2 := bitmapfont.Render(req.Text, txt_color, txt_bg, 1, 0)
	if err2 != nil {
		log.Error(err2)
		return nil, err2
	}

	var resp *srv.DrawResponse
	var err error
	for i := 0; ; i++ {

		if i*39 > (img.Bounds().Dx() + 39) {
			break
		}
		r := image.Pt(39*i, 0)
		m := image.NewRGBA(image.Rect(0, 0, 39, 9))
		draw.Draw(m, m.Bounds(), img, r, draw.Src)
		buf := &bytes.Buffer{}
		png.Encode(buf, m)

		resp, err = s.sink.Draw(context.Background(), &srv.DrawRequest{
			Data: buf.Bytes(),
		})
		if err != nil {
			return nil, err
		}
		time.Sleep(1 * time.Second)
	}

	return &srv.TextResponse{
		Code:   resp.Code,
		Status: resp.Status,
	}, nil
}

func main() {
	var err error

	flag.Parse()

	conf, err = config.Load(*cfg)
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
