package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"gitlab.com/meyerzinn/smasteroids/assets"
	"gitlab.com/meyerzinn/smasteroids/smasteroids"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"sync/atomic"
	"time"
)

const footerMessageText = "Press [ENTER] to start."

type MainscreenScene struct {
	titleMessage      *text.Text
	versionMessage    *text.Text
	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
	canvas            *pixelgl.Canvas
}

func (s *MainscreenScene) Render(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyEnter) {
		s.Destroy()
		Current = Play()
	}

	// make the footer blink
	select {
	case <-s.footerBlinkTicker.C:
		s.footerActive.Store(false)
		time.AfterFunc(time.Second/10, func() {
			s.footerActive.Store(true)
		})
	default:
	}
	// clear the window
	s.canvas.Clear(colornames.Black)
	// show the game title
	bounds := s.titleMessage.Bounds()
	//matrix := pixel.IM.Moved(canvas.Bounds().Center().ScaledXY(pixel.V(.5, 2.0/3.0)).Sub(bounds.Center()))
	matrix := pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*2/3)).Sub(bounds.Center()))
	s.titleMessage.Draw(s.canvas, matrix)
	// show the footer message
	if s.footerActive.Load().(bool) {
		bounds = s.footerMessage.Bounds()
		matrix = pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*1/5)).Sub(bounds.Center()))
		s.footerMessage.Draw(s.canvas, matrix)
	}

	s.versionMessage.Draw(s.canvas, pixel.IM.Moved(CanvasBounds.Min).Moved(pixel.V(4, 4)))

	Draw(win, s.canvas)
}

func (s *MainscreenScene) Destroy() {
	// stop the footer blinking ticker
	s.footerBlinkTicker.Stop()
}

func Start() Scene {
	titleMessage := text.New(pixel.V(0, 0), assets.FontTitle)
	_, _ = titleMessage.WriteString("SMasteroids")
	footerMessage := text.New(pixel.ZV, assets.FontInterface)
	_, _ = footerMessage.WriteString(footerMessageText)
	footerBlinkTicker := time.NewTicker(time.Second)
	var footerActive atomic.Value
	footerActive.Store(true)

	versionMessage := text.New(pixel.ZV, text.NewAtlas(basicfont.Face7x13, text.ASCII))
	_, _ = versionMessage.WriteString("Version " + smasteroids.Version() + ".")
	return &MainscreenScene{
		titleMessage:      titleMessage,
		versionMessage:    versionMessage,
		footerMessage:     footerMessage,
		footerBlinkTicker: footerBlinkTicker,
		footerActive:      footerActive,
		canvas:            pixelgl.NewCanvas(CanvasBounds),
	}
}
