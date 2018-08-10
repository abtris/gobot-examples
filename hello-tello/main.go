/*
How to run

        go run hello-tello/main.go
*/
package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8888")
	var flightData *tello.FlightData
	var battery int8
	work := func() {
		drone.TakeOff()

		drone.On(tello.FlightDataEvent, func(data interface{}) {
			flightData = data.(*tello.FlightData)
			battery = flightData.BatteryPercentage
			fmt.Println("Height:", flightData.Height)
		})

		gobot.After(5*time.Second, func() {
			fmt.Println("Battery:", battery)
			drone.Land()
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
