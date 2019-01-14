package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/meyerzinn/smastroids/assets"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

type levelTitleScene struct {
	levelIndex int
	levelText  *text.Text
	titleText  *text.Text
	nextTimer  *time.Timer
}

func (s *levelTitleScene) Render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	select {
	case <-s.nextTimer.C:
		Current = Level(s.levelIndex)
		return
	default:
	}
	canvas.Clear(colornames.Black)
	bounds := s.levelText.Bounds()
	matrix := pixel.IM.Moved(canvas.Bounds().Min.Add(pixel.V(canvas.Bounds().W()/2, canvas.Bounds().H()*5/9)).Sub(bounds.Center()))
	s.levelText.Draw(canvas, matrix)
	bounds = s.titleText.Bounds()
	matrix = pixel.IM.Moved(canvas.Bounds().Min.Add(pixel.V(canvas.Bounds().W()/2, canvas.Bounds().H()*4/9)).Sub(bounds.Center()))
	s.titleText.Draw(canvas, matrix)
}

func (s *levelTitleScene) Destroy() {
}

func Play() Scene {
	return TitleScene(0)
}

func TitleScene(index int) Scene {
	level := assets.Levels[index]
	levelText := text.New(pixel.ZV, assets.FontSubtitle)
	_, _ = levelText.WriteString("Level " + strconv.Itoa(level.Number) + ":")
	titleText := text.New(pixel.ZV, assets.FontTitle)
	titleText.Color = colornames.Yellow
	_, _ = titleText.WriteString(level.Name)
	return &levelTitleScene{
		levelIndex: index,
		levelText:  levelText,
		titleText:  titleText,
		nextTimer:  time.NewTimer(4 * time.Second),
	}
}
