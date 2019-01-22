package sprites

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	TextureShip   *pixelgl.Canvas
	TextureBullet *pixelgl.Canvas
)

var (
	shipVertices = []pixel.Vec{{-48, -48}, {0, 96}, {48, -48}}
	bulletRadius = 16.0
)

// Init should be called from the main thread before any sprite textures are accessed.
func Init() {
	TextureShip = pixelgl.NewCanvas(pixel.R(-48, -48, 48, 96))
	TextureBullet = pixelgl.NewCanvas(pixel.R(-16, -16, 16, 16))

	imd := imdraw.New(nil)
	// player sprite
	imd.Color = colornames.White
	for _, v := range shipVertices {
		imd.Push(v)
	}
	imd.Polygon(8)
	imd.Draw(TextureShip)
	imd.Reset()
	// bullet sprite
	imd.Color = colornames.White
	imd.Push(pixel.ZV)
	imd.Circle(bulletRadius, 4)
	imd.Draw(TextureBullet)
	imd.Reset()
}
