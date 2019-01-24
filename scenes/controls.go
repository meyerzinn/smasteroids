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

	joystickControlSchemes = map[string]JoystickControlSchemeFactory{
		"8Bitdo SFC30 GamePad": make8BitdoSFC30GamePadControlScheme,
	}
)

type JoystickControlSchemeFactory func(pixelgl.Joystick) ControlScheme

func make8BitdoSFC30GamePadControlScheme(joystick pixelgl.Joystick) ControlScheme {
	return ControlScheme{
		Left:   JoystickAxisInputMethod{Joystick: joystick, Axis: 0, Inverse: true, Threshold: .1, Alias: "LEFT"},
		Right:  JoystickAxisInputMethod{Joystick: joystick, Axis: 0, Inverse: false, Threshold: .1, Alias: "RIGHT"},
		Shoot:  JoystickButtonInputMethod{Joystick: joystick, Button: 0, Alias: "A"},
		Boost:  JoystickButtonInputMethod{Joystick: joystick, Button: 7, Alias: "R"},
		Thrust: JoystickButtonInputMethod{Joystick: joystick, Button: 6, Alias: "L"},
	}
}

type AnyInputMethod struct {
	Methods []InputMethod
}

func (im AnyInputMethod) GetInput(win *pixelgl.Window) bool {
	for i := 0; i < len(im.Methods); i++ {
		if im.Methods[i].GetInput(win) {
			return true
		}
	}
	return false
}

func (im AnyInputMethod) String() string {
	switch len(im.Methods) {
	case 0:
		return ""
	case 1:
		return im.Methods[0].String()
	default:
		const deliminator = " | "
		var val string
		for _, method := range im.Methods {
			val += method.String() + " | "
		}
		return val[:len(val)-3]
	}
}

func AnyInput(methods ...InputMethod) InputMethod {
	return AnyInputMethod{
		Methods: methods,
	}
}

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
		if im.Inverse && extent < 0 {
			return true
		}
		if !im.Inverse && extent > 0 {
			return true
		}
	}
	return false
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
