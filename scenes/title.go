package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/meyerzinn/smastroids/assets"
	"golang.org/x/image/colornames"
	"time"
)

const footerMessageText = "Press any key to start."

type titleScene struct {
	titleMessage *text.Text

	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
}

func (s *titleScene) Render(win *pixelgl.Window) {
	select {
	case <-s.footerBlinkTicker.C:
		s.footerMessage.Clear()
		time.AfterFunc(time.Second/10, s.footerUnblink)
	default:

	}
	win.Clear(colornames.Black)
	bounds := s.titleMessage.Bounds()
	winbounds := win.Bounds()
	pos := pixel.IM.Moved(winbounds.Max.ScaledXY(pixel.V(.5, 2.0/3.0)).Sub(bounds.Center()))
	s.titleMessage.Draw(win, pos)

	bounds = s.footerMessage.Bounds()
	pos = pixel.IM.Moved(winbounds.Max.ScaledXY(pixel.V(.5, 1.0/3.0)).Sub(bounds.Center()))
	s.footerMessage.Draw(win, pos)
}

func (s *titleScene) Destroy() {
	return
}

func (s *titleScene) footerUnblink() {
	_, _ = s.footerMessage.WriteString(footerMessageText)
}

func TitleScene(win *pixelgl.Window) Scene {
	titleAtlas := text.NewAtlas(assets.FontTitle, text.ASCII)
	titleMessage := text.New(pixel.V(0, 0), titleAtlas)
	_, _ = titleMessage.WriteString("SMastroids")

	footerAtlas := text.NewAtlas(assets.FontInterface, text.ASCII)
	footerMessage := text.New(pixel.ZV, footerAtlas)
	_, _ = footerMessage.WriteString(footerMessageText)
	footerBlinkTicker := time.NewTicker(time.Second)
	return &titleScene{
		titleMessage:      titleMessage,
		footerMessage:     footerMessage,
		footerBlinkTicker: footerBlinkTicker,
	}
}
