package scenes

import (
	"github.com/20zinnm/smasteroids/assets"
	"github.com/20zinnm/smasteroids/smasteroids"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

const LevelTitleDelay = 3 * time.Second

type levelTitleScene struct {
	levelIndex int
	levelText  *text.Text
	titleText  *text.Text
	nextTimer  *time.Timer
}

func (s *levelTitleScene) Render(win *pixelgl.Window) {
	select {
	case <-s.nextTimer.C:
		TransitionTo(NewLevelScene(s.levelIndex))
		return
	default:
	}
	win.Clear(colornames.Black)
	win.SetMatrix(pixel.IM)
	bounds := s.levelText.Bounds()
	matrix := pixel.IM.Moved(win.Bounds().Min.Add(pixel.V(win.Bounds().W()/2, win.Bounds().H()*5/9)).Sub(bounds.Center()))
	s.levelText.Draw(win, matrix)
	bounds = s.titleText.Bounds()
	matrix = pixel.IM.Moved(win.Bounds().Min.Add(pixel.V(win.Bounds().W()/2, win.Bounds().H()*4/9)).Sub(bounds.Center()))
	s.titleText.Draw(win, matrix)
}

func NewTitleScene(index int) Scene {
	level := smasteroids.Levels[index]
	levelText := text.New(pixel.ZV, assets.FontSubtitle)
	_, _ = levelText.WriteString("Level " + strconv.Itoa(index+1) + ":")
	titleText := text.New(pixel.ZV, assets.FontTitle)
	titleText.Color = colornames.Yellow
	_, _ = titleText.WriteString(level.Name)
	return &levelTitleScene{
		levelIndex: index,
		levelText:  levelText,
		titleText:  titleText,
		nextTimer:  time.NewTimer(LevelTitleDelay),
	}
}
