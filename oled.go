package main

import (
	"image"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3"
)

type OledDevice struct {
	DeviceHandle ssd1306.Dev
	I2CBus       i2c.BusCloser
	ImageData    []*image.Gray
	GifData      *gif.GIF
}

func (dev *OledDevice) InitDefault() {
	dev.Init(128, 32, false, false, false)
}

func (dev *OledDevice) Init(w int, h int, rotated bool, sequential bool, swapTopBottom bool) {
	// Load all the drivers
	_, err := host.Init()
	if err != nil {
		panic(err)
	}

	// Open a handle to the first available I²C bus
	bus, err := i2creg.Open("")
	if err != nil {
		panic(err)
	}

	oledOpts := ssd1306.Opts{
		W:             128,
		H:             32,
		Rotated:       false,
		Sequential:    false,
		SwapTopBottom: false,
	}

	// Open a handle to a ssd1306 connected on the I²C bus
	oledDev, err := ssd1306.NewI2C(bus, &oledOpts)
	if err != nil {
		log.Fatal(err)
	}

	dev.DeviceHandle = *oledDev
	dev.I2CBus = bus
}

func (dev *OledDevice) Close() {
	dev.DeviceHandle.Halt()
	dev.I2CBus.Close()
}

func (dev *OledDevice) setupGif(gifPath string) {
	// Open the gif file
	f, err := os.Open(gifPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read the gif file
	g, err := gif.DecodeAll(f)
	if err != nil {
		panic(err)
	}

	imgs := make([]*image.Gray, len(g.Image))
	for i := range g.Image {
		imgs[i] = convertAndResizeAndCenter(dev.DeviceHandle.Bounds().Dx(), dev.DeviceHandle.Bounds().Dy(), g.Image[i])
	}
	dev.ImageData = imgs
	dev.GifData = g
}

func (dev *OledDevice) DisplayGif(gifPath string) {
	dev.setupGif(gifPath)

	for i := 0; ; i++ {
		index := i % len(dev.ImageData)
		c := time.After(time.Duration(10*dev.GifData.Delay[index]) * time.Millisecond)
		img := dev.ImageData[index]
		dev.DeviceHandle.Draw(img.Bounds(), img, image.Point{})
		<-c
	}
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
