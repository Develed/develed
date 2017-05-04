package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	srv "github.com/develed/develed/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type TextMsg struct {
	Font    string  `json:"font"`
	Text    string  `json:"text"`
	BgColor []uint8 `json:"bg_color"`
}

var (
	debug = flag.Bool("debug", false, "enter debug mode")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [opts...] HOST:PORT\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	var err error

	if flag.Parse(); flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	conn, err := grpc.Dial(flag.Arg(0), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := srv.NewImageSinkClient(conn)
	dec := json.NewDecoder(os.Stdin)

loop:
	for {
		msg := TextMsg{
			Font:    "font6x8",
			BgColor: []uint8{0, 0, 0, 255},
		}

		if err := dec.Decode(&msg); err != nil {
			if err == io.EOF {
				break loop
			} else {
				log.Errorln(err)
				continue loop
			}
		}

		var font FontMgr
		fontImage := font.Init(msg.Font)

		// Allocate frame
		img := image.NewRGBA(image.Rect(0, 0, 39, 9))
		col := color.RGBA{msg.BgColor[0], msg.BgColor[1], msg.BgColor[2], msg.BgColor[3]}
		nm := img.Bounds()
		for y := 0; y < nm.Dy(); y++ {
			for x := 0; x < nm.Dx(); x++ {
				img.Set(x, y, col)
			}
		}

		// Fill frame
		for n, key := range msg.Text {
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

		resp, err := c.Draw(context.Background(), &srv.DrawRequest{
			Data: buf.Bytes(),
		})
		if err != nil {
			log.Fatalln(err)
			continue
		}

		log.Infoln(resp)
	}
}
