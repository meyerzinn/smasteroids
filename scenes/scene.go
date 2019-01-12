package scenes

import (
	"github.com/faiface/pixel/pixelgl"
)

var (
	// Current represents the currently active scene which is rendered.
	Current Scene
)

// Scene represents a logical division of the game (like a slide in PowerPoint).
type Scene interface {
	// Render is called repeatedly as long as this Scene is the currently active scene (the value of Current).
	Render(win *pixelgl.Window, canvas *pixelgl.Canvas)
}
