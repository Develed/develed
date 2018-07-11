package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"net"
	"net/http"
	"os"
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
	Time       time.Duration
}

var cRenderTextChannel = make(chan RenderCtx, 1)
var cRenderImgChannel = make(chan RenderCtx, 1)
var priorita bool = false
var cSyncChannel = make(chan bool)
var cFrameWidth int = 39
var cFrameHigh int = 9
var cScollText string = "scroll"
var cFixText string = "fix"
var cCenterText string = "center"
var cBlinkText string = "blink"

var currentTemp float64 = 0.0

const OWMUrl = "http://api.openweathermap.org/data/2.5/weather?q=Campi%20Bisenzio,Firenze&units=metric&appid="

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

	cRenderTextChannel <- RenderCtx{text_img, charWidth, conf.Textd.TextScrollTime * time.Millisecond, "scroll", conf.Textd.TextStayTime}

	return &srv.TextResponse{
		Code:   0,
		Status: "Ok",
	}, nil
}

func blitFrame(sink ImageSink, img image.Image, draw_rect image.Rectangle) error {
	frame := image.NewRGBA(image.Rect(0, 0, cFrameWidth, cFrameHigh))
	if img != nil {
		draw.Draw(frame, draw_rect, img, image.ZP, draw.Src)
	}
	buf := &bytes.Buffer{}
	png.Encode(buf, frame)
	_, err := sink.Draw(context.Background(), &srv.DrawRequest{
		Data: buf.Bytes(),
	})
	if err != nil {
		return err
	}
	return nil
}

func textRenderEfx(sink ImageSink, img image.Image, ctx RenderCtx) error {
	var err error
	switch ctx.efxType {
	case cScollText:
		for frame_idx := 0; ; frame_idx++ {
			// Scrolling time..
			time.Sleep(ctx.scrollTime)
			err = blitFrame(sink, img, image.Rect(cFrameWidth-frame_idx, 0, cFrameWidth, cFrameHigh))
			if err != nil {
				return err
			}
			if frame_idx >= (img.Bounds().Dx() + cFrameWidth) {
				log.Debug("End frame wrap..")
				return nil
			}
		}
	case cFixText:
		err = blitFrame(sink, img, image.Rect(0, 0, cFrameWidth, cFrameHigh))
		if err != nil {
			return err
		}
	case cCenterText:
		off := cFrameWidth - img.Bounds().Dx()
		if off > 0 {
			off = off / 2
		} else {
			off = 0
		}
		err = blitFrame(sink, img, image.Rect(off, 0, cFrameWidth-off, cFrameHigh))
		if err != nil {
			return err
		}
	}
	return nil
}

func renderLoop(sink ImageSink) {
	fmt.Println("render loop")
	ctx := RenderCtx{nil, cFrameWidth, 0, "fix", 0}
	for {
		select {
		case ctx = <-cRenderImgChannel:
			log.Debug("Text Render channel")

		default:
			// Message from a channel lets render it
			if ctx.img != nil {
				textRenderEfx(sink, ctx.img, ctx)
			}
		}
	}
}

type LoopFunc func() RenderCtx

func generazioneImmagini() {
	fmt.Println("generazione immagini")

	clock := clock
	var loop = []struct {
		f  LoopFunc
		tt time.Duration
	}{
		{binClock, 90 * time.Second},
		{clock, 5 * time.Second},
		{date, 2 * time.Second},
		{temperature, 10 * time.Second},
	}
	cont := 0
	ticker := time.NewTicker(1 * time.Second)
	appo := loop[cont]
	cRenderImgChannel <- appo.f()
	now := time.Now()

	for {
		select {
		case cxt := <-cRenderTextChannel:
			if cxt.img != nil {
				cRenderImgChannel <- cxt
				time.Sleep(cxt.Time)
			}
		case tempo := <-ticker.C:
			if (tempo.Sub(now) * time.Microsecond) < (loop[cont].tt * time.Microsecond) {
				appo := loop[cont]
				cRenderImgChannel <- appo.f()
			} else {
				cont = cont + 1
				if cont == len(loop) {
					cont = 0
				}
				now = time.Now()
			}
		}

	}
}

var clockFlag bool = true
var txtColor color.RGBA = color.RGBA{13, 25, 64, 128}
var txtBg color.RGBA = color.RGBA{0, 0, 0, 255}

func binClock() RenderCtx {
	var err error
	var loc *time.Location

	//set timezone,
	loc, err = time.LoadLocation("Europe/Rome")
	if err != nil {
		log.Error("Unable go get time clock..")
		panic(err)
	}

	time_str := time.Now().In(loc).Format("150405")

	err = bitmapfont.Init(conf.Textd.FontPath, conf.Textd.BitdatetimeFont, conf.BitmapFonts)
	if err != nil {
		panic(err)
	}

	text_img, charWidth, err := bitmapfont.Render(time_str, txtColor, txtBg, 1, 0)
	return RenderCtx{text_img, charWidth, 100 * time.Millisecond, "center", 500 * time.Millisecond}
}

func clock() RenderCtx {
	fmt.Println("clock")
	var err error
	var loc *time.Location

	//set timezone,
	loc, err = time.LoadLocation("Europe/Rome")
	if err != nil {
		log.Error("Unable go get time clock..")
		panic(err)
	}

	now := time.Now().In(loc)
	time_str := ""
	time_str = now.Format("15:04")
	if clockFlag {
		time_str = now.Format("15 04")
	}
	clockFlag = !clockFlag

	err = bitmapfont.Init(conf.Textd.FontPath, conf.Textd.DatetimeFont, conf.BitmapFonts)
	if err != nil {
		panic(err)
	}

	text_img, charWidth, err := bitmapfont.Render(time_str, txtColor, txtBg, 1, 0)
	return RenderCtx{text_img, charWidth, 100 * time.Millisecond, "center", 500 * time.Millisecond}
}

func date() RenderCtx {
	fmt.Println("clock")
	var err error
	var loc *time.Location

	//set timezone,
	loc, err = time.LoadLocation("Europe/Rome")
	if err != nil {
		log.Error("Unable go get time clock..")
		panic(err)
	}
	now := time.Now().In(loc)
	time_str := now.Format("02.01.06")

	err = bitmapfont.Init(conf.Textd.FontPath, conf.Textd.DatetimeFont, conf.BitmapFonts)
	if err != nil {
		panic(err)
	}

	text_img, charWidth, err := bitmapfont.Render(time_str, txtColor, txtBg, 1, 0)
	return RenderCtx{text_img, charWidth, 100 * time.Millisecond, "center", 500 * time.Millisecond}
}

func temperature() RenderCtx {
	var err error

	err = bitmapfont.Init(conf.Textd.FontPath, conf.Textd.TemperatureFont, conf.BitmapFonts)
	if err != nil {
		panic(err)
	}

	var sign = '+'
	if currentTemp < 0 {
		sign = '-'
	}

	txt := fmt.Sprintf("%c%2.1fC", sign, currentTemp)
	text_img, charWidth, err := bitmapfont.Render(txt, txtColor, txtBg, 1, 0)

	return RenderCtx{text_img, charWidth, 100 * time.Millisecond, "center", 500 * time.Millisecond}
}

type ImageSink interface {
	srv.ImageSinkServer
	Run() error
}

func main() {
	var sink ImageSink
	var err error

	flag.Parse()

	conf, err = config.Load(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	owm_token := os.Getenv("OWM_API_TOKEN")
	if owm_token == "" {
		owm_token = conf.Textd.OWMToken
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if !*debug {
		sink, err = NewDeviceSink("/dev/sscdev0")
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Println("debug")

		sink, err = NewTermSink()
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println("BO")
		}
	}

	sock, err := net.Listen("tcp", conf.Textd.GRPCServerAddress)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	drawing_srv := &server{sink: srv.NewImageSinkClient(nil)}
	srv.RegisterTextdServer(s, drawing_srv)
	reflection.Register(s)

	go renderLoop(sink)
	go generazioneImmagini()

	go func() {
		retrieveTemperature := func() {
			resp, err := http.Get(OWMUrl + owm_token)
			if err != nil {
				log.Errorln(err)
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorln(err)
				return
			}

			var jmap map[string]interface{}
			if err := json.Unmarshal(data, &jmap); err != nil {
				log.Errorln(err)
				return
			}

			currentTemp = jmap["main"].(map[string]interface{})["temp"].(float64)
		}

		retrieveTemperature()

		for {
			select {
			case <-time.After(5 * time.Minute):
				retrieveTemperature()
			}
		}
	}()

	if err := s.Serve(sock); err != nil {
		log.Fatalln(err)
	}
}
