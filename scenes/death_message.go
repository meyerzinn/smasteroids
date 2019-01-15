package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/meyerzinn/smasteroids/assets"
	"golang.org/x/image/colornames"
	"math/rand"
	"sync/atomic"
	"time"
)

type deathMessageScene struct {
	text   *text.Text
	canvas *pixelgl.Canvas

	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
}

func (s *deathMessageScene) Render(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeySpace) {
		s.Destroy()
		Current = Start()
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

	s.canvas.Clear(colornames.Black)
	bounds := s.text.Bounds()
	matrix := pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*2/3)).Sub(bounds.Center()))
	s.text.Draw(s.canvas, matrix)
	if s.footerActive.Load().(bool) {
		bounds := s.footerMessage.Bounds()
		matrix = pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*1/3)).Sub(bounds.Center()))
		s.footerMessage.Draw(s.canvas, matrix)
	}

	Draw(win, s.canvas)
}

func (s *deathMessageScene) Destroy() {
	// stop the footer blinking ticker
	s.footerBlinkTicker.Stop()
}

func Death() Scene {
	var teacher, quote string
	for t, quotes := range assets.DeathMessages { // take advantage of Go's guaranteed random ranging
		teacher = t
		quote = quotes[rand.Intn(len(quotes))]
		break
	}
	name := assets.Teachers[teacher]
	lines := []string{quote, "- " + name}
	txt := text.New(pixel.ZV, assets.FontInterface)
	for _, line := range lines {
		txt.Dot.X -= txt.BoundsOf(line).W()
		_, _ = fmt.Fprintln(txt, line)
	}

	footerMessage := text.New(pixel.ZV, assets.FontInterface)
	_, _ = footerMessage.WriteString("Press [SPACE] to continue.")
	footerBlinkTicker := time.NewTicker(time.Second)
	var footerActive atomic.Value
	footerActive.Store(true)
	return &deathMessageScene{
		text:              txt,
		canvas:            pixelgl.NewCanvas(CanvasBounds),
		footerMessage:     footerMessage,
		footerBlinkTicker: footerBlinkTicker,
		footerActive:      footerActive,
	}
}
