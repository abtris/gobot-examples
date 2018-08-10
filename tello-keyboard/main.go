/*
How to run:
Connect to the drone's Wi-Fi network from your computer. It will be named something like "TELLO-XXXXXX".

Once you are connected you can run the Gobot code on your computer to control the drone.

        go run tello-keyboard/main.go
*/

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gobot.io/x/gobot/platforms/keyboard"
)

func main() {
	drone := tello.NewDriver("8888")
	keys := keyboard.NewDriver()

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			keys.On(keyboard.Key, func(data interface{}) {
				key := data.(keyboard.KeyEvent)
				switch key.Key {
				case keyboard.A:
					fmt.Println(key.Char)
					drone.Clockwise(-25)
				case keyboard.D:
					fmt.Println(key.Char)
					drone.Clockwise(25)
				case keyboard.W:
					fmt.Println(key.Char)
					drone.Forward(25)
				case keyboard.S:
					fmt.Println(key.Char)
					drone.Backward(25)
				case keyboard.K:
					fmt.Println(key.Char)
					drone.Down(20)
				case keyboard.J:
					fmt.Println(key.Char)
					drone.Up(20)
				case keyboard.Q:
					fmt.Println(key.Char)
					drone.Land()
				case keyboard.P:
					fmt.Println(key.Char)
					drone.TakeOff()
				case keyboard.Escape:
					drone.Forward(0)
					drone.Backward(0)
					drone.Up(0)
					drone.Down(0)
					drone.Left(0)
					drone.Right(0)
					drone.Clockwise(0)
				}
			})
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{keys, drone},
		work,
	)

	robot.Start()
}
