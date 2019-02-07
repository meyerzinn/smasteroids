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

const winFooterMessageText = "Player 1: Press [Boost] to continue."

type winScene struct {
	titleMessage *text.Text

	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
}

func (s *winScene) Render(win *pixelgl.Window) {
	if Players[0].Boost.GetInput(win) {
		TransitionTo(NewMainscreenScene())
		return
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
	win.Clear(colornames.Black)
	win.SetMatrix(pixel.IM)
	// show the game title
	bounds := s.titleMessage.Bounds()
	matrix := pixel.IM.Moved(win.Bounds().Min.Add(pixel.V(win.Bounds().W()/2, win.Bounds().H()*2/3)).Sub(bounds.Center()))
	s.titleMessage.Draw(win, matrix)
	// show the footer message
	if s.footerActive.Load().(bool) {
		bounds = s.footerMessage.Bounds()
		matrix = pixel.IM.Moved(win.Bounds().Min.Add(pixel.V(win.Bounds().W()/2, win.Bounds().H()*1/3)).Sub(bounds.Center()))
		s.footerMessage.Draw(win, matrix)
	}
}

func (s *winScene) Destroy() {
	// stop the footer blinking ticker
	s.footerBlinkTicker.Stop()
}

func NewWin() Scene {
	titleMessage := text.New(pixel.V(0, 0), assets.FontTitle)
	_, _ = titleMessage.WriteString("You Winn!")
	footerMessage := text.New(pixel.ZV, assets.FontInterface)
	_, _ = footerMessage.WriteString(winFooterMessageText)
	footerBlinkTicker := time.NewTicker(time.Second)
	var footerActive atomic.Value
	footerActive.Store(true)
	return &winScene{
		titleMessage:      titleMessage,
		footerMessage:     footerMessage,
		footerBlinkTicker: footerBlinkTicker,
		footerActive:      footerActive,
	}
}
