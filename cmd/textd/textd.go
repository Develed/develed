package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/develed/develed/bitmapfont"
	"github.com/develed/develed/config"
	srv "github.com/develed/develed/services"
	"google.golang.org/grpc"
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

var cRenderImgChannel = make(chan RenderCtx, 1)
var cRenderTextChannel = make(chan RenderCtx, 1)
var cSyncChannel = make(chan bool)
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

	cSyncChannel <- true
	cRenderTextChannel <- RenderCtx{text_img, charWidth, conf.Textd.TextScrollTime * time.Millisecond, "scroll"}

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

func renderLoop(dr_srv *server) {
	ctx := RenderCtx{nil, cFrameWidth, 0, "fix"}

	for {
		select {
		case ctx = <-cRenderImgChannel:
			log.Debug("Text Render channel")
		default:
			// Message from a channel lets render it
			if ctx.img != nil {
				blitFrame(dr_srv, ctx.img, image.Rect(0, 0, cFrameWidth, cFrameHigh))
			}
		}
	}
}

func generazioneImmagini(dr_srv *server) {
	// var err error
	// txt_color := color.RGBA{255, 0, 0, 255}
	// txt_bg := color.RGBA{0, 0, 0, 255}

	//cxt := RenderCtx{nil, cFrameWidth, 0, "fix"}

	var clockTickElapse time.Duration = 5 * time.Second
	fmt.Println(clockTickElapse)

	for {
		select {
		case <-cSyncChannel:
			clockTickElapse = conf.Textd.TextStayTime * time.Second
		case cxt := <-cRenderTextChannel:
			if cxt.img != nil {
				cRenderImgChannel <- cxt
			}
		default:
			// text_img := image.NewRGBA(image.Rect(0, 0, cFrameWidth, cFrameHigh))
			// draw.Draw(text_img, text_img.Bounds(), &image.Uniform{color.RGBA{0, 255, 0, 255}}, image.ZP, draw.Src)
			// blitFrame(dr_srv, text_img, image.Rect(0, 0, cFrameWidth, cFrameHigh))
			clock(dr_srv) //ultima cosa aggiunta, non so se va bene
		}
	}

}

func clock(dr_srv *server) {
	var err error
	var loc *time.Location
	txt_color := color.RGBA{255, 0, 0, 255}
	txt_bg := color.RGBA{0, 0, 0, 255}

	//set timezone,
	loc, err = time.LoadLocation("Europe/Rome")
	if err != nil {
		log.Error("Unable go get time clock..")
		panic(err)
	}
	var flag bool = true
	var show_date int = 30
	now := time.Now().In(loc)
	time_str := ""
	if now.Unix()%int64(conf.Textd.DateStayTime) == 0 && (show_date <= 0) {
		time_str = now.Format("02.01.06")
		show_date = 30
	} else {
		flag = !flag
		time_str = now.Format("15:04")
		if flag {
			time_str = now.Format("15 04")
		}
		show_date--
	}

	err = bitmapfont.Init(conf.Textd.FontPath, conf.Textd.DatetimeFont, conf.BitmapFonts)
	if err != nil {
		panic(err)
	}
	text_img, charWidth, err := bitmapfont.Render(time_str, txt_color, txt_bg, 1, 0)
	ctx := RenderCtx{text_img, charWidth, 200 * time.Millisecond, "center"}
	if ctx.img != nil {
		text_img = image.NewRGBA(ctx.img.Bounds())
	}

	blitFrame(dr_srv, text_img, image.Rect(0, 0, cFrameWidth, cFrameHigh))
}

func main() {
	fmt.Println("CIAOOO	")
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

	go renderLoop(drawing_srv)
	go generazioneImmagini(drawing_srv)

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
