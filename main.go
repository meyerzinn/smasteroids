package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/meyerzinn/smastroids/scenes"
	_ "image/png"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "smasteroids",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	defer win.Destroy()
	if err != nil {
		panic(err)
	}
	scenes.Current = scenes.LoadingScene()

	for !win.Closed() {
		scenes.Current.Render(win)
		if win.Pressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
