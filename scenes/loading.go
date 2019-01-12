package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/meyerzinn/smastroids/assets"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"time"
)

var (
	Current Scene
)

func next(scene Scene) {
	Current.Destroy()
	Current = scene
}

type Scene interface {
	Render(win *pixelgl.Window)
	Destroy()
}

type sceneFunc func()

func (fn sceneFunc) Render() {
	fn()
}

type loadingScene struct {
	ticker         *time.Ticker
	dots           string
	doneLoading    chan struct{}
	loadingMessage *text.Text
}

func (s *loadingScene) Render(win *pixelgl.Window) {
	select {
	case <-s.doneLoading:
		next(TitleScene(win))
		return
	default:
		win.Clear(colornames.Black)
		select {
		case <-s.ticker.C:
			s.dots += "."
			if len(s.dots) > 3 {
				s.dots = ""
			}
		default:
		}
		s.loadingMessage.Clear()
		_, _ = s.loadingMessage.WriteString("Loading" + s.dots)
		bounds := s.loadingMessage.Bounds()
		pos := pixel.IM.Moved(win.Bounds().Center().Sub(bounds.Center()))
		s.loadingMessage.Draw(win, pos)
	}
}

func (s *loadingScene) Destroy() {
	s.ticker.Stop()
}

func LoadingScene() Scene {
	doneLoading := make(chan struct{})
	go func() {
		assets.Init()
		time.Sleep(5 * time.Second)
		close(doneLoading)
	}()
	var atlas = text.NewAtlas(
		basicfont.Face7x13,
		text.ASCII,
	)
	loadingMessage := text.New(pixel.V(0, 0), atlas)
	dotTicker := time.NewTicker(time.Second)
	return &loadingScene{
		ticker:         dotTicker,
		doneLoading:    doneLoading,
		loadingMessage: loadingMessage,
	}
}
