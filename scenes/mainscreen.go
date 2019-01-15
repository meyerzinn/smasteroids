package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/meyerzinn/smasteroids/assets"
	"golang.org/x/image/colornames"
	"sync/atomic"
	"time"
)

const footerMessageText = "Press [SPACE] to start."

type MainscreenScene struct {
	titleMessage *text.Text

	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
	canvas            *pixelgl.Canvas
}

func (s *MainscreenScene) Render(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeySpace) {
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
		matrix = pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*1/3)).Sub(bounds.Center()))
		s.footerMessage.Draw(s.canvas, matrix)
	}
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
	return &MainscreenScene{
		titleMessage:      titleMessage,
		footerMessage:     footerMessage,
		footerBlinkTicker: footerBlinkTicker,
		footerActive:      footerActive,
		canvas:            pixelgl.NewCanvas(CanvasBounds),
	}
}
