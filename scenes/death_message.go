package scenes

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/meyerzinn/smastroids/game"
	"sync/atomic"
	"time"
)

type deathMessageScene struct {
	teacher           game.TeacherData
	
	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
}

func (s *deathMessageScene) Render(win *pixelgl.Window, canvas *pixelgl.Canvas) {

}
