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

type RenderCtx struct {
	img        image.Image
	charWidth  int
	scrollTime time.Duration
}

var cRenderChannel = make(chan RenderCtx, 1)
var cFrameWidth int = 39
var cFrameHigh int = 9

func (s *server) Write(ctx context.Context, req *srv.TextRequest) (*srv.TextResponse, error) {
	var err error
	err = bitmapfont.Init(conf.Textd.FontPath, req.Font, conf.BitmapFonts)
	if err != nil {
		return &srv.TextResponse{
			Code:   -1,
			Status: err.Error(),
		}, nil
	}

	txt_color := color.RGBA{255, 0, 0, 255}
	txt_bg := color.RGBA{0, 0, 0, 255}
	text_img, charWidth, err := bitmapfont.Render(req.Text, txt_color, txt_bg, 1, 0)
	if err != nil {
		return &srv.TextResponse{
			Code:   -1,
			Status: err.Error(),
		}, nil
	}

	cRenderChannel <- RenderCtx{text_img, charWidth, 300 * time.Millisecond}
	log.Debugf("Color: %v Bg: %v\n", req.FontColor, req.FontBg)

	return &srv.TextResponse{
		Code:   0,
		Status: "Ok",
	}, nil
}

func renderLoop(dr_srv *server) {
	var err error
	frameCtx := RenderCtx{
		nil,
		cFrameWidth,
		500 * time.Millisecond,
	}

	text_img := image.NewRGBA(image.Rect(0, 0, cFrameWidth, cFrameHigh))
	//draw.Draw(text_img, text_img.Bounds(), &image.Uniform{color.RGBA{0, 255, 0, 255}}, image.ZP, draw.Src)

	for {
		select {
		case ctx := <-cRenderChannel:
			log.Debug("Render channel")
			frameCtx.charWidth = ctx.charWidth
			frameCtx.scrollTime = ctx.scrollTime

			text_img = image.NewRGBA(ctx.img.Bounds())
			draw.Draw(text_img, ctx.img.Bounds(), ctx.img, image.ZP, draw.Src)
		default:
			for frame_idx := 0; ; frame_idx++ {
				time.Sleep(frameCtx.scrollTime)
				if text_img == nil {
					log.Debug("..")
					break
				}

				frame := image.NewRGBA(image.Rect(0, 0, cFrameWidth, cFrameHigh))

				// Fill frame only with max char lenght
				maxCharNum := (cFrameWidth / frameCtx.charWidth)
				xStart := (maxCharNum - frame_idx) * frameCtx.charWidth
				xStop := maxCharNum * frameCtx.charWidth

				draw.Draw(frame, image.Rect(xStart, 0, xStop, cFrameHigh), text_img, image.ZP, draw.Src)
				buf := &bytes.Buffer{}
				png.Encode(buf, frame)

				_, err = dr_srv.sink.Draw(context.Background(), &srv.DrawRequest{
					Data: buf.Bytes(),
				})

				if (frame_idx+1)*frameCtx.charWidth >= (text_img.Bounds().Dx() + cFrameWidth) {
					break
				}

				if err != nil {
					break
				}
			}
		}
	}
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
	drawing_srv := &server{sink: srv.NewImageSinkClient(conn)}
	srv.RegisterTextdServer(s, drawing_srv)
	reflection.Register(s)

	go renderLoop(drawing_srv)

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
