package assets

import (
	"bytes"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"image"
	_ "image/png"
)

var (
	FontTitle     *text.Atlas
	FontInterface *text.Atlas
	FontSubtitle  *text.Atlas
	FontLabel     *text.Atlas
)

var (
	Icon pixel.Picture
)

func init() {
	assets := packr.NewBox("./data")

	// load fonts
	ps2pRaw, err := assets.Find("fonts/PressStart2P.ttf")
	if err != nil {
		panic(errors.Wrap(err, "loading font PressStart2P"))
	}
	title, err := loadTTF(ps2pRaw, 52)
	if err != nil {
		panic(errors.Wrap(err, "initializing title font face"))
	}
	FontTitle = text.NewAtlas(title, text.ASCII)

	inter, err := loadTTF(ps2pRaw, 18)
	if err != nil {
		panic(errors.Wrap(err, "initializing interface font face"))
	}
	FontInterface = text.NewAtlas(inter, text.ASCII)

	subtitle, err := loadTTF(ps2pRaw, 36)
	if err != nil {
		panic(errors.Wrap(err, "initializing subtitle font face"))
	}
	FontSubtitle = text.NewAtlas(subtitle, text.ASCII)

	label, err := loadTTF(ps2pRaw, 12)
	if err != nil {
		panic(errors.Wrap(err, "initializing label font face"))
	}
	FontLabel = text.NewAtlas(label, text.ASCII)

	Icon, err = loadPicture(assets, "icon.png")
	if err != nil {
		panic(errors.Wrap(err, "loading icon"))
	}
}

func loadTTF(raw []byte, size float64) (font.Face, error) {
	f, err := truetype.Parse(raw)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(f, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func loadPicture(box packr.Box, path string) (pixel.Picture, error) {
	imageBytes, err := box.Find(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not load image asset")
	}
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
