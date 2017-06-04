package bitmapfont

import (
	"image"
	"image/color"
	"os"

	log "github.com/Sirupsen/logrus"
)

type FontMgr struct {
	img   image.Image
	width int
	high  int
	col   int
	row   int
}

type FontInterface interface {
	Render(fontPath string, fontName string, text string) (image.Image, error)
}

var cTable = map[string][]int{
	// Font     name,   high width rows, columns
	"font5x7": {7, 5, 4, 24},
	//"font4x7": {4, 7},
	//"font6x6": {6, 6},
	//"font6x8": {6, 8},
	//"font6x7": {6, 7},
}

func (f *FontMgr) Render(fontPath string, fontName string, text string, char_space int, top_off int) (image.Image, error) {

	if fontName == "" {
		fontName = "font5x7"
	}

	reader, err := os.Open(fontPath + string(os.PathSeparator) + fontName + ".png")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer reader.Close()

	// Decode fonts table.
	cImg, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	f.img = cImg
	f.high = cTable[fontName][0]
	f.width = cTable[fontName][1]
	f.row = cTable[fontName][2]
	f.col = cTable[fontName][3]

	frame_width := len(text)*f.width + (len(text)-1)*(char_space)
	log.Debug("len ", frame_width)

	// Allocate frame
	img := image.NewRGBA(image.Rect(0, 0, frame_width, 9))
	col := color.RGBA{0, 0, 0, 255}
	nm := img.Bounds()
	for y := 0; y < nm.Dy(); y++ {
		for x := 0; x < nm.Dx(); x++ {
			img.Set(x, y, col)
		}
	}
	log.Debug(text)

	// Fill frame
	for n, key := range text {
		col := int(key-' ') % f.col
		row := int(key-' ') / f.col

		log.Debug("offset ", int(key-' '))
		log.Debug("Col ", col)
		log.Debug("Row ", row)

		for y := 0; y < f.high+top_off; y++ {
			for x := 0; x < f.width+char_space; x++ {
				if x >= f.width {
					img.Set(x+n*(f.width+char_space), y+top_off, color.RGBA{0, 0, 0, 255})
				} else {
					img.Set(x+n*(f.width+char_space), y+top_off, f.img.At(x+f.width*col, y+f.high*row))
				}
			}
		}
	}

	return img, nil
}
