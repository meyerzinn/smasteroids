package scenes

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

var (
	defaultKeyboardControls = ControlScheme{
		Thrust: KeyboardInputMethod{Button: pixelgl.KeyW},
		Left:   KeyboardInputMethod{Button: pixelgl.KeyA},
		Right:  KeyboardInputMethod{Button: pixelgl.KeyD},
		Boost:  KeyboardInputMethod{Button: pixelgl.KeyE},
		Shoot:  KeyboardInputMethod{Button: pixelgl.KeySpace},
	}

	activeJoystickers = make(map[pixelgl.Joystick]struct{})

	joystickControlSchemes = map[string]ControlScheme{
		"8Bitdo SFC30 GamePad": {
			Left:   JoystickButtonInputMethod{Button: 6, Alias: "L"},
			Right:  JoystickButtonInputMethod{Button: 7, Alias: "R"},
			Shoot:  JoystickButtonInputMethod{Button: 0, Alias: "A"},
			Boost:  JoystickButtonInputMethod{Button: 3, Alias: "X"},
			Thrust: JoystickAxisInputMethod{Axis: 1, Inverse: false, Threshold: .1, Alias: "UP"},
		},
	}
)

//type AnyInputMethod struct {
//	Methods []InputMethod
//}
//
//func (im AnyInputMethod) GetInput(win *pixelgl.Window) bool {
//	panic("implement me")
//}
//
//func (im AnyInputMethod) String() string {
//	panic("implement me")
//}
//
//func AnyInput(methods ...InputMethod) InputMethod

type ControlScheme struct {
	Thrust, Left, Right, Shoot, Boost InputMethod
}

func (cs ControlScheme) Controls(win *pixelgl.Window) Controls {
	return Controls{
		Thrust: cs.Thrust.GetInput(win),
		Left:   cs.Left.GetInput(win),
		Right:  cs.Right.GetInput(win),
		Shoot:  cs.Shoot.GetInput(win),
		Boost:  cs.Boost.GetInput(win),
	}
}

type KeyboardInputMethod struct {
	Button pixelgl.Button
}

func (im KeyboardInputMethod) GetInput(win *pixelgl.Window) bool {
	return win.Pressed(im.Button)
}

func (im KeyboardInputMethod) String() string {
	return im.Button.String()
}

type JoystickButtonInputMethod struct {
	Joystick pixelgl.Joystick
	Button   int
	Alias    string
}

func (im JoystickButtonInputMethod) GetInput(win *pixelgl.Window) bool {
	return win.JoystickPressed(im.Joystick, im.Button)
}

func (im JoystickButtonInputMethod) String() string {
	if len(im.Alias) > 0 {
		return im.Alias
	}
	return fmt.Sprintf("Button%d", im.Button)
}

type JoystickAxisInputMethod struct {
	Joystick pixelgl.Joystick
	Axis     int
	// Inverse makes GetInput return true when |axis|>threshold && axis < 0.
	Inverse   bool
	Threshold float64
	Alias     string
}

func (im JoystickAxisInputMethod) GetInput(win *pixelgl.Window) bool {
	extent := win.JoystickAxis(im.Joystick, im.Axis)
	if math.Abs(extent) > im.Threshold {
		return !im.Inverse
	}
	return im.Inverse
}

func (im JoystickAxisInputMethod) String() string {
	if len(im.Alias) > 0 {
		return im.Alias
	}
	return fmt.Sprintf("Axis%d", im.Axis)
}

type InputMethod interface {
	GetInput(win *pixelgl.Window) bool
	String() string
}
