package scenes

import (
	"fmt"
	"github.com/20zinnm/smasteroids/assets"
	"github.com/20zinnm/smasteroids/smasteroids"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"sync/atomic"
	"time"
)

const MaxPlayers = 2

const footerMessageText = "Player 1: Press [Boost] to start."

type MainscreenScene struct {
	titleMessage      *text.Text
	versionMessage    *text.Text
	controlsMessage   *text.Text
	footerMessage     *text.Text
	footerBlinkTicker *time.Ticker
	footerActive      atomic.Value
	canvas            *pixelgl.Canvas
}

var controlsLabels = []string{"Thrust", "Turn Left", "Turn Right", "Shoot", "Boost"}

const joinText = "Connect another joystick to join."

func (s *MainscreenScene) Render(win *pixelgl.Window) {
	if len(Players) > 0 && Players[0].Boost.GetInput(win) {
		TransitionTo(Play())
		return
	}

	// Remove non-primary players pressing boost.
	for i := len(Players) - 1; i > 0; i-- {
		if Players[i].Boost.GetInput(win) {
			copy(Players[i:], Players[i+1:])
			Players[len(Players)-1] = ControlScheme{}
			Players = Players[:len(Players)-1]
			if joystick, ok := playerJoysticks[i]; ok {
				delete(playerJoysticks, i)
				delete(joystickPlayers, joystick)
			}
		}
	}

	// Add new players with joysticks.
	for joystick := pixelgl.Joystick1; joystick < pixelgl.Joystick16 && len(Players) < MaxPlayers; joystick++ {
		if win.JoystickPresent(joystick) {
			if _, ok := joystickPlayers[joystick]; !ok {
				if scheme, ok := joystickControlSchemes[win.JoystickName(joystick)]; ok {
					// we have a known joystick, add the player
					Players = append(Players, scheme(joystick))
					playerJoysticks[len(Players)-1] = joystick
					joystickPlayers[joystick] = len(Players) - 1
				}
			}
		}
	}

	// Default controls is the keyboard scheme.
	if len(Players) == 0 {
		Players = append(Players, defaultKeyboardControls)
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
	// > Show labels.
	s.controlsMessage.Clear()
	fmt.Fprintln(s.controlsMessage, "Controls:")
	for _, l := range controlsLabels {
		fmt.Fprintln(s.controlsMessage, l)
	}
	s.controlsMessage.Draw(s.canvas, pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()/5, s.canvas.Bounds().H()/2)).Sub(s.controlsMessage.Bounds().Center())))
	for i, scheme := range Players {
		s.controlsMessage.Clear()
		fmt.Fprintf(s.controlsMessage, "Player %d\n", i+1)
		fmt.Fprintln(s.controlsMessage, scheme.Thrust.String())
		fmt.Fprintln(s.controlsMessage, scheme.Left.String())
		fmt.Fprintln(s.controlsMessage, scheme.Right.String())
		fmt.Fprintln(s.controlsMessage, scheme.Shoot.String())
		fmt.Fprintln(s.controlsMessage, scheme.Boost.String())
		s.controlsMessage.Draw(s.canvas, pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()*(float64(i)+2)/5, s.canvas.Bounds().H()/2)).Sub(s.controlsMessage.Bounds().Center())))
	}
	// Show join message for all possible players not joined.
	if len(Players) != MaxPlayers {
		s.controlsMessage.Clear()
		fmt.Fprintf(s.controlsMessage, joinText)
		s.controlsMessage.Draw(s.canvas, pixel.IM.Moved(s.canvas.Bounds().Min.Add(pixel.V(s.canvas.Bounds().W()*0.5, s.canvas.Bounds().H()*2/5)).Sub(s.controlsMessage.Bounds().Center())))
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
	_, _ = titleMessage.WriteString("Crewsteroids")
	controlsMessage := text.New(pixel.ZV, assets.FontInterface)
	footerMessage := text.New(pixel.ZV, assets.FontInterface)
	_, _ = footerMessage.WriteString(footerMessageText)
	footerBlinkTicker := time.NewTicker(time.Second)
	var footerActive atomic.Value
	footerActive.Store(true)
	versionMessage := text.New(pixel.ZV, text.NewAtlas(basicfont.Face7x13, text.ASCII))
	_, _ = versionMessage.WriteString("Version " + smasteroids.Version() + ". Developed by Meyer Zinn.")

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
