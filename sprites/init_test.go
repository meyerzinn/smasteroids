package sprites_test

import (
	"github.com/20zinnm/smasteroids/sprites"
	"github.com/faiface/pixel/pixelgl"
	"testing"
)

func TestInit(t *testing.T) {
	pixelgl.Run(func() {
		sprites.Init()
		if sprites.TextureShip == nil {
			t.Fail()
		}
		if sprites.TextureBullet == nil {
			t.Fail()
		}
	})
}
