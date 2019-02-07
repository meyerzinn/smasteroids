package scenes

import (
	"fmt"
	"github.com/20zinnm/smasteroids/assets"
	"github.com/20zinnm/smasteroids/smasteroids"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"math/rand"
	"sync/atomic"
	"time"
)

type deathMessageScene struct {
	text  *text.Text
	level int

	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
}

func (s *deathMessageScene) Render(win *pixelgl.Window) {
	if Players[0].Boost.GetInput(win) {
		TransitionTo(NewTitleScene(s.level))
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

	win.Clear(colornames.Black)
	win.SetMatrix(pixel.IM)
	bounds := s.text.Bounds()
	matrix := pixel.IM.Moved(win.Bounds().Min.Add(pixel.V(win.Bounds().W()/2, win.Bounds().H()*2/3)).Sub(bounds.Center()))
	s.text.Draw(win, matrix)
	if s.footerActive.Load().(bool) {
		bounds := s.footerMessage.Bounds()
		matrix = pixel.IM.Moved(win.Bounds().Min.Add(pixel.V(win.Bounds().W()/2, win.Bounds().H()*1/3)).Sub(bounds.Center()))
		s.footerMessage.Draw(win, matrix)
	}
}

func (s *deathMessageScene) Destroy() {
	// stop the footer blinking ticker
	s.footerBlinkTicker.Stop()
}

func Death(index int) Scene {
	var name, quote string
	for n, quotes := range smasteroids.DeathMessages { // take advantage of Go's guaranteed random ranging
		name = n
		quote = quotes[rand.Intn(len(quotes))]
		break
	}
	lines := []string{`"` + quote + `"`, "- " + name}
	txt := text.New(pixel.ZV, assets.FontInterface)
	for _, line := range lines {
		txt.Dot.X -= txt.BoundsOf(line).W()
		_, _ = fmt.Fprintln(txt, line)
	}

	footerMessage := text.New(pixel.ZV, assets.FontInterface)
	_, _ = footerMessage.WriteString("Player 1: Press [Boost] to continue.")
	footerBlinkTicker := time.NewTicker(time.Second)
	var footerActive atomic.Value
	footerActive.Store(true)
	return &deathMessageScene{
		text:              txt,
		level:             index,
		footerMessage:     footerMessage,
		footerBlinkTicker: footerBlinkTicker,
		footerActive:      footerActive,
	}
}
