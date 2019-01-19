package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	// current represents the currently active scene which is rendered.
	current Scene

	CanvasBounds = pixel.R(-1920/2, -1080/2, 1920/2, 1080/2)
)

// Scene represents a logical division of the game (like a slide in PowerPoint).
type Scene interface {
	// Render is called repeatedly as long as this Scene is the currently active scene (the value of current).
	Render(win *pixelgl.Window)
}

type DestroyableScene interface {
	Scene
	Destroy()
}

func Draw(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	win.SetMatrix(pixel.IM.Scaled(pixel.ZV, win.Bounds().W()/canvas.Bounds().W()).Moved(win.Bounds().Center()))
	canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
}

func TransitionTo(o Scene) {
	if current != nil {
		if d, ok := o.(DestroyableScene); ok {
			d.Destroy()
		}
	}
	current = o
}

func Render(win *pixelgl.Window) {
	if current != nil {
		current.Render(win)
	}
}
