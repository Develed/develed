package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage %s: <inputPath> <output path>", os.Args[0])
		os.Exit(-1)
	}

	inputPath := string(os.Args[1])
	outputPath := string(os.Args[2])

	cDebug := false
	if os.Getenv("TXTDEBUG") == "1" {
		fmt.Println("Debug Mode..")
		cDebug = true
	}

	for {
		c, _ := os.OpenFile(inputPath, os.O_RDONLY, 0600)
		defer c.Close()

		reader := bufio.NewReader(c)
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
			fmt.Println(err)
			continue
		}

		fmt.Println(cfg)

		var font FontMgr
		cfgFont := "font6x8"
		if cfg["font"] == "" {
			fmt.Println("No font specify use default..font6x8")
		} else {
			cfgFont = cfg["font"]
		}

		fontImage := font.Init(cfgFont)

		var r int = 0
		var g int = 0
		var b int = 0
		var a int = 255

		if cfg["bg_color"] == "" {
			fmt.Println("No font specify, use default [0,0,0,255]")
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

		// Export to output file
		f, _ := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0600)
		defer f.Close()

		if cDebug {
			png.Encode(f, img)

		} else {
			buf := new(bytes.Buffer)
			png.Encode(buf, img)

			encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
			fmt.Println(encoded)

			f.Write([]byte(encoded))
		}

	}
}
