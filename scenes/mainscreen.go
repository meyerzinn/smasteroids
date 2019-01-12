package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/meyerzinn/smastroids/assets"
	"golang.org/x/image/colornames"
	"sync/atomic"
	"time"
)

const footerMessageText = "Press [SPACE] to start."

type mainscreenScene struct {
	titleMessage *text.Text

	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
}

func (s *mainscreenScene) Render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
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
	canvas.Clear(colornames.Black)
	// show the game title
	bounds := s.titleMessage.Bounds()
	//matrix := pixel.IM.Moved(canvas.Bounds().Center().ScaledXY(pixel.V(.5, 2.0/3.0)).Sub(bounds.Center()))
	matrix := pixel.IM.Moved(canvas.Bounds().Min.Add(pixel.V(canvas.Bounds().W()/2, canvas.Bounds().H() * 2/3)).Sub(bounds.Center()))
	s.titleMessage.Draw(canvas, matrix)
	// show the footer message
	if s.footerActive.Load().(bool) {
		bounds = s.footerMessage.Bounds()
		matrix = pixel.IM.Moved(canvas.Bounds().Min.Add(pixel.V(canvas.Bounds().W()/2, canvas.Bounds().H() * 1/3)).Sub(bounds.Center()))
		s.footerMessage.Draw(canvas, matrix)
	}
}

func (s *mainscreenScene) Destroy() {
	// stop the footer blinking ticker
	s.footerBlinkTicker.Stop()
}

func Start() Scene {
	titleMessage := text.New(pixel.V(0, 0), assets.FontTitle)
	_, _ = titleMessage.WriteString("SMastroids")
	footerMessage := text.New(pixel.ZV, assets.FontInterface)
	_, _ = footerMessage.WriteString(footerMessageText)
	footerBlinkTicker := time.NewTicker(time.Second)
	var footerActive atomic.Value
	footerActive.Store(true)
	return &mainscreenScene{
		titleMessage:      titleMessage,
		footerMessage:     footerMessage,
		footerBlinkTicker: footerBlinkTicker,
		footerActive:      footerActive,
	}
}
