package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Joystick Testing",
		Bounds:    pixel.R(0, 0, 1080, 720),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	defer win.Destroy()
	if err != nil {
		panic(err)
	}
	for joystick := pixelgl.Joystick1; joystick < pixelgl.Joystick16; joystick++ {
		present := win.JoystickPresent(joystick)
		fmt.Printf("Test whether joystick (index=%d) is present: %v\n", joystick, present)
		if present {
			fmt.Printf("> Name: %s\n", win.JoystickName(joystick))
			buttons := win.JoystickButtonCount(joystick)
			fmt.Printf("> Buttons (%d):\n", buttons)
			for button := 0; button < buttons; button++ {
				fmt.Printf("> > Test whether button (index=%d) is pressed: %v\n", button, win.JoystickPressed(joystick, button))
			}
			axes := win.JoystickAxisCount(joystick)
			fmt.Printf("> Axes (%d):\n", axes)
			for axis := 0; axis < axes; axis++ {
				fmt.Printf("> > Test axis (index=%d): %v\n", axis, win.JoystickAxis(joystick, axis))
			}
		}
	}
}
