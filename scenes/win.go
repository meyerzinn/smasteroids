package scenes

import (
	"github.com/20zinnm/smasteroids/assets"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"sync/atomic"
	"time"
)

const winFooterMessageText = "Press [Boost] to start again."

type WinScene struct {
	titleMessage *text.Text

	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
	canvas            *pixelgl.Canvas
}

func (s *WinScene) Render(win *pixelgl.Window) {
	if Players[0].Boost.GetInput(win) {
		s.Destroy()
		current = Play()
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

func (s *WinScene) Destroy() {
	// stop the footer blinking ticker
	s.footerBlinkTicker.Stop()
}

func Win() Scene {
	titleMessage := text.New(pixel.V(0, 0), assets.FontTitle)
	_, _ = titleMessage.WriteString("You Winn!")
	footerMessage := text.New(pixel.ZV, assets.FontInterface)
	_, _ = footerMessage.WriteString(winFooterMessageText)
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
