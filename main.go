package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/meyerzinn/smastroids/scenes"
	"golang.org/x/image/colornames"
	_ "image/png"
	"math"
	"time"
)

func run() {
	primaryMonitor := pixelgl.PrimaryMonitor()
	width, height := primaryMonitor.Size()
	cfg := pixelgl.WindowConfig{
		Title:   "SMasteroids",
		Bounds:  pixel.R(0, 0, width, height),
		VSync:   true,
		Monitor: primaryMonitor,
	}
	win, err := pixelgl.NewWindow(cfg)
	defer win.Destroy()
	if err != nil {
		panic(err)
	}
	//win.SetMatrix(pixel.IM.Scaled(win.Bounds().Center(), width/1024.0))
	canvas := pixelgl.NewCanvas(pixel.R(-1920/2, -1080/2, 1920/2, 1080/2))
	scenes.Current = scenes.Start()
	tickDuration := time.Duration(math.Floor((1.0/primaryMonitor.RefreshRate())*math.Pow10(9))) * time.Nanosecond
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()
	for !win.Closed() {
		<-ticker.C // wait for next tick
		win.Clear(colornames.Black)
		scenes.Current.Render(win, canvas)
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV, math.Min(
			win.Bounds().W()/canvas.Bounds().W(),
			win.Bounds().H()/canvas.Bounds().H(),
		)).Moved(win.Bounds().Center()))
		canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
		if win.Pressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
