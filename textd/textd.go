package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	debug = flag.Bool("debug", false, "enter debug mode")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [opts...] [INPUT] [OUTPUT]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	var err error

	in := os.Stdin
	out := os.Stdout

	flag.Parse()

	if flag.NArg() > 0 {
		in, err = os.OpenFile(flag.Arg(0), os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		defer in.Close()
	}

	if flag.NArg() > 1 {
		out, err = os.OpenFile(flag.Arg(1), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		defer out.Close()
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	reader := bufio.NewReader(in)

	for {
		var cfgLine []string
		var cfg = make(map[string]string, 10)
		line, err := reader.ReadString('\n')
		for err == nil {
			if !strings.HasPrefix(line, "#") {
				line = strings.TrimSpace(line)
				if len(line) > 0 {
					cfgLine = strings.Split(line, "=")
					cfg[cfgLine[0]] = cfgLine[1]
				}
			}
			line, err = reader.ReadString('\n')
		}

		if err != io.EOF {
			log.Errorln(err)
			continue
		}

		log.Debugln(cfg)

		var font FontMgr
		cfgFont := "font6x8"
		if cfg["font"] == "" {
			log.Debugln("No font specify use default..font6x8")
		} else {
			cfgFont = cfg["font"]
		}

		fontImage := font.Init(cfgFont)

		var r int = 0
		var g int = 0
		var b int = 0
		var a int = 255

		if cfg["bg_color"] == "" {
			log.Debugln("No font specify, use default [0,0,0,255]")
		} else {
			r, _ = strconv.Atoi(strings.Split(cfg["bg_color"], ",")[0])
			g, _ = strconv.Atoi(strings.Split(cfg["bg_color"], ",")[1])
			b, _ = strconv.Atoi(strings.Split(cfg["bg_color"], ",")[2])
			a, _ = strconv.Atoi(strings.Split(cfg["bg_color"], ",")[3])
		}

		// Allocate frame
		img := image.NewRGBA(image.Rect(0, 0, 39, 9))
		nm := img.Bounds()
		for y := 0; y < nm.Dy(); y++ {
			for x := 0; x < nm.Dx(); x++ {
				img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			}
		}

		// Fill frame
		for n, key := range cfg["text"] {
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

		if !*debug {
			binary.Write(out, binary.LittleEndian, uint64(buf.Len()))
		}
		out.Write(buf.Bytes())
	}
}
