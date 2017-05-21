package bitmapfont

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"os"

	log "github.com/Sirupsen/logrus"
	conf "github.com/develed/develed/config"
)

var fontImageTable image.Image
var Config conf.BitmapFont

func Render(text string, text_color color.RGBA, text_bg color.RGBA, char_space int, top_off int) (image.Image, error) {

	fontcolums := fontImageTable.Bounds().Dx() / Config.Width
	frame_width := len(text)*Config.Width + (len(text)-1)*(char_space)
	log.Debug("Frame len in px:", frame_width)

	m := image.NewRGBA(image.Rect(0, 0, frame_width, 9))
	draw.Draw(m, m.Bounds(), &image.Uniform{text_bg}, image.ZP, draw.Src)

	src := &image.Uniform{text_color}

	for n, key := range text {
		col := int(key-' ') % fontcolums
		row := int(key-' ') / fontcolums

		draw.DrawMask(m, image.Rect(n*(Config.Width+char_space), 0,
			Config.Width+n*(Config.Width+char_space), Config.High),
			src, image.ZP, fontImageTable, image.Pt(col*Config.Width, row*Config.High), draw.Over)

		log.Debugf("key: %c off: %v c: %v r: %v", key, int(key-' '), col, row)
	}
	return m, nil
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
