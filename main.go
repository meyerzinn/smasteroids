package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gitlab.com/meyerzinn/smasteroids/scenes"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
)

//go:generate packr

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func run() {
	primaryMonitor := pixelgl.PrimaryMonitor()
	width, height := primaryMonitor.Size()
	cfg := pixelgl.WindowConfig{
		Title:  "SMasteroids",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
		//Resizable: true,
		Monitor: primaryMonitor,
	}
	win, err := pixelgl.NewWindow(cfg)
	defer win.Destroy()
	if err != nil {
		panic(err)
	}
	scenes.CanvasBounds = win.Bounds().Moved(win.Bounds().Center().Scaled(-1))
	CenterWindow(win)
	//win.SetMatrix(pixel.IM.Scaled(win.Bounds().Center(), width/1024.0))
	scenes.Current = scenes.Start()
	tickDuration := time.Duration(math.Floor((1.0/primaryMonitor.RefreshRate())*math.Pow10(9))) * time.Nanosecond
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()
	for !win.Closed() {
		win.Clear(colornames.Black)
		scenes.Current.Render(win)
		if win.Pressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		win.Update()
		<-ticker.C // wait for next tick
	}
}

func CenterWindow(win *pixelgl.Window) {
	x, y := pixelgl.PrimaryMonitor().Size()
	width, height := win.Bounds().Size().XY()
	win.SetPos(
		pixel.V(
			x/2-width/2,
			y/2-height/2,
		),
	)
}

func main() {
	pixelgl.Run(run)
}
