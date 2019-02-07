package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

var (
	// current represents the currently active scene which is rendered.
	current Scene

	// GameBounds represents the dimensions of the game (used for physics).

)

// Scene represents a logical division of the game (like a slide in PowerPoint).
type Scene interface {
	// Render is called repeatedly as long as this Scene is the currently active scene (the value of current).
	Render(win *pixelgl.Window)
}

// DestroyableScene represents a scene that must be explicitly disposed of.
type DestroyableScene interface {
	Scene
	Destroy()
}

// DrawCanvas draws a canvas to the window, properly scaling and moving the canvas to fit within the window.
//
// The aspect ratio of the canvas is maintained when scaled.
func DrawCanvas(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
		math.Min(
			win.Bounds().W()/canvas.Bounds().W(),
			win.Bounds().H()/canvas.Bounds().H(),
		),
	).Moved(win.Bounds().Center()))
	canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
}

// TransitionTo changes the currently active scene, destroying the previous scene if necessary.
func TransitionTo(o Scene) {
	if current != nil {
		if d, ok := o.(DestroyableScene); ok {
			d.Destroy()
		}
	}
	current = o
}

// Render calls the render method of the currently active scene if it exists.
func Render(win *pixelgl.Window) {
	if current != nil {
		current.Render(win)
	}
}
