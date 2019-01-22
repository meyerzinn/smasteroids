package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"gitlab.com/meyerzinn/smasteroids/assets"
	"gitlab.com/meyerzinn/smasteroids/smasteroids"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"sync/atomic"
	"time"
)

const MaxPlayers = 2

const footerMessageText = "Press [ENTER] to start."

type MainscreenScene struct {
	titleMessage      *text.Text
	versionMessage    *text.Text
	controlsMessage   *text.Text
	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
	canvas            *pixelgl.Canvas
}

var keyboardControlsText = [][]string{{
	"Player 1 Controls:",
	"Thrust     - [W]",
	"Turn Left  - [A]",
	"Turn Right - [D]",
	"Fire       - [SPACE]",
	"Boost      - [E]",
}, {
	"Player 2 Controls:",
	"Thrust     - [UP]",
	"Turn Left  - [LEFT]",
	"Turn Right - [RIGHT]",
	"Fire       - [RIGHT CTRL]",
	"Boost      - [RIGHT SHIFT]",
}}

var joinText = []string{
	"Press [UP] to join.",
}

func (s *MainscreenScene) Render(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyEnter) {
		TransitionTo(Play())
	}

	// Handle player join.
	if win.JustPressed(pixelgl.KeyUp) {
		if len(Players) == 2 {
			// Player leaves.
			Players = Players[:1]
		} else {
			Players = append(Players, ControllerInputFn(func(win *pixelgl.Window) Controls {
				// hardcoded for now, should probably change this
				return Controls{
					Thrust: win.Pressed(pixelgl.KeyUp),
					Left:   win.Pressed(pixelgl.KeyLeft),
					Right:  win.Pressed(pixelgl.KeyRight),
					Boost:  win.Pressed(pixelgl.KeyRightShift),
					Shoot:  win.Pressed(pixelgl.KeyRightControl),
				}
			}))
		}
	}

	// Blink the footer.
	select {
	case <-s.footerBlinkTicker.C:
		s.footerActive.Store(false)
		time.AfterFunc(time.Second/10, func() {
			s.footerActive.Store(true)
		})
	default:
	}

	// Clear the window.
	s.canvas.Clear(colornames.Black)
	// Show the game title.
	bounds := s.titleMessage.Bounds()
	s.titleMessage.Draw(s.canvas, pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*4/5)).Sub(bounds.Center())))

	// Show controls message for all currently joined players.
	for i := range Players {
		s.controlsMessage.Clear()
		for _, l := range keyboardControlsText[i] {
			_, _ = fmt.Fprintln(s.controlsMessage, l)
		}
		s.controlsMessage.Draw(s.canvas, pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()*float64(i+1)/3, s.canvas.Bounds().H()/2)).Sub(s.controlsMessage.Bounds().Center())))
	}
	// Show join message for all possible players not joined.
	for i := len(Players); i < MaxPlayers; i++ {
		s.controlsMessage.Clear()
		_, _ = s.controlsMessage.WriteString(joinText[i-1])
		s.controlsMessage.Draw(s.canvas, pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()*float64(i+1)/3, s.canvas.Bounds().H()/2)).Sub(s.controlsMessage.Bounds().Center())))
	}

	// show the footer message
	if s.footerActive.Load().(bool) {
		bounds = s.footerMessage.Bounds()
		s.footerMessage.Draw(s.canvas, pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/2, s.canvas.Bounds().H()*1/5)).Sub(bounds.Center())))
	}

	s.versionMessage.Draw(s.canvas, pixel.IM.Moved(CanvasBounds.Min).Moved(pixel.V(4, 4)))

	Draw(win, s.canvas)
}

func (s *MainscreenScene) Destroy() {
	// stop the footer blinking ticker
	s.footerBlinkTicker.Stop()
}

func Start() Scene {
	titleMessage := text.New(pixel.V(0, 0), assets.FontTitle)
	_, _ = titleMessage.WriteString("SMasteroids")
	controlsMessage := text.New(pixel.ZV, assets.FontInterface)
	footerMessage := text.New(pixel.ZV, assets.FontInterface)
	_, _ = footerMessage.WriteString(footerMessageText)
	footerBlinkTicker := time.NewTicker(time.Second)
	var footerActive atomic.Value
	footerActive.Store(true)
	versionMessage := text.New(pixel.ZV, text.NewAtlas(basicfont.Face7x13, text.ASCII))
	_, _ = versionMessage.WriteString("Version " + smasteroids.Version() + ". Developed by Meyer Zinn.")
	// Add the first player.
	Players = []ControllerInput{
		ControllerInputFn(func(win *pixelgl.Window) Controls {
			return Controls{
				Thrust: win.Pressed(pixelgl.KeyW),
				Left:   win.Pressed(pixelgl.KeyA),
				Right:  win.Pressed(pixelgl.KeyD),
				Shoot:  win.Pressed(pixelgl.KeySpace),
				Boost:  win.Pressed(pixelgl.KeyE),
			}
		}),
	}

	return &MainscreenScene{
		titleMessage:      titleMessage,
		versionMessage:    versionMessage,
		controlsMessage:   controlsMessage,
		footerMessage:     footerMessage,
		footerBlinkTicker: footerBlinkTicker,
		footerActive:      footerActive,
		canvas:            pixelgl.NewCanvas(CanvasBounds),
	}
}
