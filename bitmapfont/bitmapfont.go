package bitmapfont

import (
	"errors"
	"image"
	"image/color"
	"os"

	log "github.com/Sirupsen/logrus"
	conf "github.com/develed/develed/config"
)

var fontImageTable image.Image
var Config conf.BitmapFont

func Render(text string, char_space int, top_off int) (image.Image, error) {

	img_rect := fontImageTable.Bounds()
	fontcolums := img_rect.Dx() / Config.Width

	frame_width := len(text)*Config.Width + (len(text)-1)*(char_space)
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
		col := int(key-' ') % fontcolums
		row := int(key-' ') / fontcolums

		log.Debug("offset ", int(key-' '))
		log.Debug("Col ", col)
		log.Debug("Row ", row)

		for y := 0; y < Config.High+top_off; y++ {
			for x := 0; x < Config.Width+char_space; x++ {
				if x >= Config.Width {
					img.Set(x+n*(Config.Width+char_space), y+top_off, color.RGBA{0, 0, 0, 255})
				} else {
					img.Set(x+n*(Config.Width+char_space), y+top_off, fontImageTable.At(x+Config.Width*col, y+Config.High*row))
				}
			}
		}
	}

	return img, nil
}

func Init(path string, name string, cfg []conf.BitmapFont) error {
	if name == "" {
		name = "font5x7"
	}

	for _, s := range cfg {
		log.Debug("Cerco font ", s.Name, " ", name)
		if name == s.Name {
			Config = s
			log.Debug(path + string(os.PathSeparator) + Config.FileName)
			reader, err := os.Open(path + string(os.PathSeparator) + Config.FileName)
			if err != nil {
				return err
			}
			defer reader.Close()

			// Decode fonts table.
			fontImageTable, _, err = image.Decode(reader)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("Wrong BitmatFont name.\n")
}
