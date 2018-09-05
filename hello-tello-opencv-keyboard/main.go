/*
You must have ffmpeg and OpenCV installed in order to run this code. It will connect to the Tello
and then open a window using OpenCV showing the streaming video.

How to run

	go run examples/tello_opencv.go
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os/exec"
	"strconv"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gobot.io/x/gobot/platforms/keyboard"
	"gocv.io/x/gocv"
)

const (
	frameX    = 640
	frameY    = 480
	frameSize = frameX * frameY * 3
)

func resetDronePostion(drone *tello.Driver) {
	drone.Forward(0)
	drone.Backward(0)
	drone.Up(0)
	drone.Down(0)
	drone.Left(0)
	drone.Right(0)
	drone.Clockwise(0)
}

func main() {
	drone := tello.NewDriver("8890")
	keys := keyboard.NewDriver()

	var flightData *tello.FlightData

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
			drone.Forward(20)
		case keyboard.S:
			fmt.Println(key.Char)
			drone.Backward(20)
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
		case keyboard.ArrowUp:
			fmt.Println(key.Char)
			drone.FrontFlip()
		case keyboard.ArrowDown:
			fmt.Println(key.Char)
			drone.BackFlip()
		case keyboard.ArrowLeft:
			fmt.Println(key.Char)
			drone.LeftFlip()
		case keyboard.ArrowRight:
			fmt.Println(key.Char)
			drone.RightFlip()
		case keyboard.Escape:
			resetDronePostion(drone)
		}
	})

	window := gocv.NewWindow("Tello")
	xmlFile := "haarcascade_frontalface_default.xml"
	ffmpeg := exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
		"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1")
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()

	work := func() {
		if err := ffmpeg.Start(); err != nil {
			fmt.Println(err)
			return
		}

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
			drone.SetExposure(0)

			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})

		drone.On(tello.FlightDataEvent, func(data interface{}) {
			flightData = data.(*tello.FlightData)
			fmt.Println("Height:", flightData.Height)
		})

	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{keys, drone},
		work,
	)

	// calling Start(false) lets the Start routine return immediately without an additional blocking goroutine
	robot.Start(false)

	// now handle video frames from ffmpeg stream in main thread, to be macOS/Windows friendly
	for {
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}
		img, _ := gocv.NewMatFromBytes(frameY, frameX, gocv.MatTypeCV8UC3, buf)
		if img.Empty() {
			continue
		}

		// detect faces
		// color for the rect when faces detected
		blue := color.RGBA{0, 0, 255, 0}
		// load classifier to recognize faces
		classifier := gocv.NewCascadeClassifier()
		defer classifier.Close()
		if !classifier.Load(xmlFile) {
			fmt.Printf("Error reading cascade file: %v\n", xmlFile)
			return
		}
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
