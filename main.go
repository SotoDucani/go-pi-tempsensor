package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3"
)

func getPinState(i int) rpio.State {
	pin := rpio.Pin(i)
	curState := pin.Read()
	fmt.Printf("Pin %v curState: %v\n", i, curState)
	return curState
}

func setPinState(i int, targetState int) {
	pin := rpio.Pin(i)
	pin.Output()
}

func getAllPinState() {
	count := 27
	for i := 0; i < count; i++ {
		getPinState(i)
	}
}

func flipPinLoop() {
	pinList := []int{17, 18, 27, 22, 23, 24, 25, 4, 5, 6, 13, 19, 26, 12, 16, 20, 21}
	for _, pin := range pinList {
		fmt.Printf("PIN:%v\n", pin)
		pinObj := rpio.Pin(pin)
		pinObj.Output()
		pinObj.Toggle()
		time.Sleep(time.Second)
		pinObj.Toggle()
	}
}

func bme280() {
	if _, err := host.Init(); err != nil {
		panic(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		panic(err)
	}
	defer bus.Close()

	dev, err := bmxx80.NewI2C(bus, 0x77, &bmxx80.DefaultOpts)
	if err != nil {
		panic(err)
	}
	defer dev.Halt()

	var env physic.Env
	if err = dev.Sense(&env); err != nil {
		panic(err)
	}
	fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
}

// convertAndResizeAndCenter takes an image, resizes and centers it on a
// image.Gray of size w*h.
func convertAndResizeAndCenter(w, h int, src image.Image) *image.Gray {
	src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
	img := image.NewGray(image.Rect(0, 0, w, h))
	r := src.Bounds()
	r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img
}

func oled() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	// Open a handle to a ssd1306 connected on the I²C bus:
	dev, err := ssd1306.NewI2C(bus, &ssd1306.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}

	// Decodes an animated GIF as specified on the command line:
	if len(os.Args) != 2 {
		log.Fatal("please provide the path to an animated GIF")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Converts every frame to image.Gray and resize them:
	imgs := make([]*image.Gray, len(g.Image))
	for i := range g.Image {
		imgs[i] = convertAndResizeAndCenter(dev.Bounds().Dx(), dev.Bounds().Dy(), g.Image[i])
	}

	// Display the frames in a loop:
	for i := 0; ; i++ {
		index := i % len(imgs)
		c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
		img := imgs[index]
		dev.Draw(img.Bounds(), img, image.Point{})
		<-c
	}
}

func forever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Shutdown signal received...")
}

func stats_Loop() {
	//var interval time.Duration = 20
	for {
		//getAllPinState()
		//bme280()
		//flipPinLoop()
		//time.Sleep(interval * time.Second)
		oled()
	}
}

func main() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer rpio.Close()

	go stats_Loop()

	forever()
}
