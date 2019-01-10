package main

import (
	"bytes"
	"encoding/json"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	"path"
)

var atlas = text.NewAtlas(
	basicfont.Face7x13,
	text.ASCII,
)

var teachers map[string]*TeacherData
var deathMessages map[string][]string
var levels []*Level

type TeacherData struct {
	Name  string
	Image pixel.Picture
}

type Level struct {
	Name       string
	Difficulty int
	Teachers   []string
}

func loadAssetsOrPanic() {
	assets := packr.NewBox("./assets")
	// load teachers and associated images
	teachersRaw, err := assets.MustBytes("teachers.json")
	panicOnError(err)
	var names map[string]string
	err = json.Unmarshal(teachersRaw, &names)
	panicOnError(err)
	teachers = make(map[string]*TeacherData, len(names))
	for last, full := range names {
		picture, err := loadPicture(assets, path.Join("images", last+".png"))
		panicOnError(err)
		teachers[last] = &TeacherData{Name: full, Image: picture}
	}
	// load death messages
	deathMessagesRaw, err := assets.MustBytes("death_messages.json")
	panicOnError(err)
	deathMessages = make(map[string][]string)
	err = json.Unmarshal(deathMessagesRaw, &deathMessages)
	panicOnError(err)
	// load levels
	levelsRaw, err := assets.MustBytes("levels.json")
	panicOnError(err)
	err = json.Unmarshal(levelsRaw, &levels)
	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadPicture(box packr.Box, path string) (pixel.Picture, error) {
	imageBytes, err := box.MustBytes(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not load image asset")
	}
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

var window *pixelgl.Window

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "smasteroids",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	window = win
	loadAssetsOrPanic()
	window.Clear(colornames.Black)
	doLoadingScreen()
}

func doLoadingScreen() {
	done := make(chan struct{})
	go func() {
		loadAssetsOrPanic()
		close(done)
	}()
	loadingMessage := text.New(pixel.V(100, 500), atlas)
	i := 0
	j := 0
	for update() {
		select {
		case <-done:
			return
		default:
			loadingMessage.Clear()
			_, _ = loadingMessage.WriteString("Loading")
			if i%60 == 0 {
				j = 3
			}
			if i%40 == 0 {
				j = 2
			} else
			if i%20 == 0 {
				j = 1
			}
			for k := 0; k < j; k++ {
				_, _ = loadingMessage.WriteString(".")
			}
			i++
			loadingMessage.Draw(window, pixel.IM)
		}
	}
}

func update() bool {
	if !window.Closed() {
		window.Update()
		return true
	} else {
		return false
	}
}

func main() {
	pixelgl.Run(run)
}
