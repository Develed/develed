package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	bitmapfont "github.com/develed/develed/bitmapfont"
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

var cRenderTextChannel = make(chan RenderCtx, 1)
var cRenderClockChannel = make(chan RenderCtx, 1)
var cSyncChannel = make(chan bool)
var cFrameWidth int = 39
var cFrameHigh int = 9
var cScollText string = "scroll"
var cFixText string = "fix"
var cCenterText string = "center"
var cBlinkText string = "blink"

func renderLoop(dr_srv *server) {
	ctx := RenderCtx{nil, cFrameWidth, 0, "fix"}
	text_img := image.NewRGBA(image.Rect(0, 0, cFrameWidth, cFrameHigh))
	//draw.Draw(text_img, text_img.Bounds(), &image.Uniform{color.RGBA{0, 255, 0, 255}}, image.ZP, draw.Src)

	for {
		select {
		case ctx = <-cRenderTextChannel:
			log.Debug("Text Render channel")
		default:
			// Message from a channel lets render it
			if ctx.img != nil {
				text_img = image.NewRGBA(ctx.img.Bounds())
				draw.Draw(text_img, ctx.img.Bounds(), ctx.img, image.ZP, draw.Src)
			}
			err := textRenderEfx(dr_srv, text_img, ctx)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}

func generazioneImmagini() {
	var err error
	var loc *time.Location

	txt_color := color.RGBA{255, 0, 0, 255}
	txt_bg := color.RGBA{0, 0, 0, 255}

	var clockTickElapse time.Duration = 5 * time.Second
	fmt.Println("%e/n", clockTickElapse)

	for {
		select {
		case <-cSyncChannel:
			clockTickElapse = conf.Textd.TextStayTime * time.Second
		case <-time.After(clockTickElapse):
			if err != nil {
				panic(err)
			}
			text_img, charWidth, err := bitmapfont.Render("FRESCHEZZE", txt_color, txt_bg, 1, 0)
			if err != nil {
				log.Error("Unable to render time clock [%v]", err.Error())
			} else {
				cRenderTextChannel <- RenderCtx{text_img, charWidth, 200 * time.Millisecond, "center"}
			}

		}
	}

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
	go generazioneImmagini()

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
