package assets

import (
	"bytes"
	"encoding/json"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype/truetype"
	"github.com/meyerzinn/smastroids/game"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"image"
	_ "image/png"
	"path"
)

var (
	Teachers      map[string]game.TeacherData
	DeathMessages map[string][]string
	Levels        []game.Level

	FontTitle     *text.Atlas
	FontInterface *text.Atlas
	FontSubtitle  *text.Atlas
	FontLabel     *text.Atlas
)

func init() {
	assets := packr.NewBox("./data")

	// load teachers and associated images
	teachersRaw, err := assets.Find("teachers.json")
	if err != nil {
		panic(errors.Wrap(err, "loading teachers.json"))
	}
	var names map[string]string
	err = json.Unmarshal(teachersRaw, &names)
	if err != nil {
		panic(errors.Wrap(err, "parsing teachers from JSON"))
	}
	Teachers = make(map[string]game.TeacherData, len(names))
	for last, full := range names {
		filepath := path.Join("images", last+".png")
		picture, err := loadPicture(assets, filepath)
		if err != nil {
			panic(errors.Wrapf(err, "loading image for teacher %s (file %s)", last, filepath))
		}
		Teachers[last] = game.TeacherData{Name: full, Image: picture}
	}

	// load death messages
	deathMessagesRaw, err := assets.Find("death_messages.json")
	if err != nil {
		panic(errors.Wrap(err, "loading death messages"))
	}
	DeathMessages = make(map[string][]string)
	err = json.Unmarshal(deathMessagesRaw, &DeathMessages)
	if err != nil {
		panic(errors.Wrap(err, "parsing death messages from JSON"))
	}

	// load levels
	levelsRaw, err := assets.Find("levels.json")
	if err != nil {
		panic(errors.Wrap(err, "loading levels"))
	}
	err = json.Unmarshal(levelsRaw, &Levels)
	if err != nil {
		panic(errors.Wrap(err, "parsing levels from JSON"))
	}
	for i, l := range Levels {
		l.Number = i + 1
		Levels[i] = l
	}

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
