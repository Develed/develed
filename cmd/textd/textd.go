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
	efxType    string
}

var cRenderTextChannel = make(chan RenderCtx, 1)
var cRenderClockChannel = make(chan RenderCtx, 1)
var cFrameWidth int = 39
var cFrameHigh int = 9
var cScollText string = "scroll"
var cFixText string = "fix"
var cCenterText string = "center"
var cBlinkText string = "blink"

func (s *server) Write(ctx context.Context, req *srv.TextRequest) (*srv.TextResponse, error) {
	var err error
	err = bitmapfont.Init(conf.Textd.FontPath, req.Font, conf.BitmapFonts)
	if err != nil {
		return &srv.TextResponse{
			Code:   -1,
			Status: err.Error(),
		}, nil
	}

	log.Debugf("Color: %v Bg: %v", req.FontColor, req.FontBg)
	txt_color := color.RGBA{255, 0, 0, 255}
	txt_bg := color.RGBA{0, 0, 0, 255}
	text_img, charWidth, err := bitmapfont.Render(req.Text, txt_color, txt_bg, 1, 0)
	if err != nil {
		return &srv.TextResponse{
			Code:   -1,
			Status: err.Error(),
		}, nil
	}

	cRenderTextChannel <- RenderCtx{text_img, charWidth, 200 * time.Millisecond}

	return &srv.TextResponse{
		Code:   0,
		Status: "Ok",
	}, nil
}

func blitFrame(dr_srv *server, img image.Image, draw_rect image.Rectangle) error {
	frame := image.NewRGBA(image.Rect(0, 0, cFrameWidth, cFrameHigh))
	if img != nil {
		draw.Draw(frame, draw_rect, img, image.ZP, draw.Src)
	}
	buf := &bytes.Buffer{}
	png.Encode(buf, frame)
	_, err := dr_srv.sink.Draw(context.Background(), &srv.DrawRequest{
		Data: buf.Bytes(),
	})
	if err != nil {
		return err
	}
	return nil
}

func textRenderEfx(dr_srv *server, img image.Image, ctx RenderCtx) error {
	var err error
	ctx := RenderCtx{nil, cFrameWidth, 100 * time.Millisecond}
	text_img := image.NewRGBA(image.Rect(0, 0, cFrameWidth, cFrameHigh))
	//draw.Draw(text_img, text_img.Bounds(), &image.Uniform{color.RGBA{0, 255, 0, 255}}, image.ZP, draw.Src)

	for {
		select {
		case ctx = <-cRenderTextChannel:
			log.Debug("Text Render channel")
			log.Debug(ctx)
		case ctx = <-cRenderClockChannel:
			log.Debug("Clock Render channel")
			log.Debug(ctx)
		default:
			// Message from a channel lets render arrived image
			if ctx.img != nil {
				text_img = image.NewRGBA(ctx.img.Bounds())
				draw.Draw(text_img, ctx.img.Bounds(), ctx.img, image.ZP, draw.Src)
			}

			for frame_idx := 0; ; frame_idx++ {
				// Scrolling time..
				time.Sleep(ctx.scrollTime)

				// Fill frame only with max char lenght
				frame := image.NewRGBA(image.Rect(0, 0, cFrameWidth, cFrameHigh))
				draw.Draw(frame, image.Rect(cFrameWidth-frame_idx, 0, cFrameWidth, cFrameHigh), text_img, image.ZP, draw.Src)
				buf := &bytes.Buffer{}
				png.Encode(buf, frame)

				_, err = dr_srv.sink.Draw(context.Background(), &srv.DrawRequest{
					Data: buf.Bytes(),
				})
				if err != nil {
					log.Error(err.Error())
					break
				}

				if frame_idx >= (text_img.Bounds().Dx() + cFrameWidth) {
					log.Debug("End frame wrap..")
					break
				}
			}
		}
	}

}

func clockLoop() {
	var err error
	txt_color := color.RGBA{255, 0, 0, 255}
	txt_bg := color.RGBA{0, 0, 0, 255}

	err = bitmapfont.Init(conf.Textd.FontPath, "", conf.BitmapFonts)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			text_img, charWidth, err := bitmapfont.Render("12:00", txt_color, txt_bg, 1, 0)
			if err != nil {
				log.Error("Unable to render time clock [%v]", err.Error())
			} else {
				cRenderClockChannel <- RenderCtx{text_img, charWidth, 200 * time.Millisecond}
				log.Debug("Clock..")

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
	go clockLoop()

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
