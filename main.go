package main

import (
	"github.com/20zinnm/smasteroids/scenes"
	"github.com/20zinnm/smasteroids/sprites"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
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
	var monitor = pixelgl.PrimaryMonitor()
	for _, m := range pixelgl.Monitors() {
		xp, _ := m.PhysicalSize()
		xo, _ := monitor.PhysicalSize()
		if xp > xo {
			monitor = m
		}
	}
	cfg := pixelgl.WindowConfig{
		Title:  "SMasteroids",
		Bounds: pixel.R(0, 0, 1920, 1080),
		VSync:  true,
		//Resizable: true,
		Monitor: monitor,
	}
	win, err := pixelgl.NewWindow(cfg)
	defer win.Destroy()
	if err != nil {
		panic(err)
	}
	win.SetCursorVisible(false)
	// initialize sprites
	sprites.Init()
	// set up canvas bounds
	scenes.CanvasBounds = win.Bounds().Moved(win.Bounds().Center().Scaled(-1))
	CenterWindow(win)
	scenes.TransitionTo(scenes.Start())
	tickDuration := time.Duration(math.Floor((1.0/60.0)*math.Pow10(9))) * time.Nanosecond
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()
	for !win.Closed() {
		for _, m := range pixelgl.Monitors() {
			xp, _ := m.PhysicalSize()
			xo, _ := monitor.PhysicalSize()
			if xp > xo {
				monitor = m
			}
		}
		win.SetMonitor(monitor)
		w, _ := monitor.Size()
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV, w/scenes.CanvasBounds.W()))
		win.Clear(colornames.Black)
		scenes.Render(win)
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
