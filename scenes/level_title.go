package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"gitlab.com/meyerzinn/smasteroids/assets"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

const LevelTitleDelay = 4 * time.Second

type LevelTitleScene struct {
	levelIndex int
	levelText  *text.Text
	titleText  *text.Text
	nextTimer  *time.Timer
	canvas     *pixelgl.Canvas
}

func (s *LevelTitleScene) Render(win *pixelgl.Window) {
	select {
	case <-s.nextTimer.C:
		Current = PlayLevel(s.levelIndex)
		return
	default:
	}
	s.canvas.Clear(colornames.Black)
	bounds := s.levelText.Bounds()
	matrix := pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*5/9)).Sub(bounds.Center()))
	s.levelText.Draw(s.canvas, matrix)
	bounds = s.titleText.Bounds()
	matrix = pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*4/9)).Sub(bounds.Center()))
	s.titleText.Draw(s.canvas, matrix)
	Draw(win, s.canvas)
}

func (s *LevelTitleScene) Destroy() {
}

func Play() Scene {
	return TitleScene(0)
}

func TitleScene(index int) Scene {
	level := assets.Levels[index]
	levelText := text.New(pixel.ZV, assets.FontSubtitle)
	_, _ = levelText.WriteString("Level " + strconv.Itoa(level.Index+1) + ":")
	titleText := text.New(pixel.ZV, assets.FontTitle)
	titleText.Color = colornames.Yellow
	_, _ = titleText.WriteString(level.Name)
	return &LevelTitleScene{
		levelIndex: index,
		levelText:  levelText,
		titleText:  titleText,
		nextTimer:  time.NewTimer(LevelTitleDelay),
		canvas:     pixelgl.NewCanvas(CanvasBounds),
	}
}
