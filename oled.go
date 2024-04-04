package main

import (
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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
	dev.Init(128, 32, false, true, false)
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
		W:             w,
		H:             h,
		Rotated:       rotated,
		Sequential:    sequential,
		SwapTopBottom: swapTopBottom,
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

func (dev *OledDevice) DisplayText(str string) {
	img := generateTextImage(dev.DeviceHandle.Bounds().Dx(), dev.DeviceHandle.Bounds().Dy(), str)
	output, err := os.Create("txt_display.jpg")
	if err != nil {
		panic(err)
	}
	defer output.Close()
	err = jpeg.Encode(output, img, nil)
	if err != nil {
		panic(err)
	}
	err = dev.DeviceHandle.Draw(image.Rect(0, 0, dev.DeviceHandle.Bounds().Dx(), dev.DeviceHandle.Bounds().Dy()), img, image.Point{})
	if err != nil {
		panic(err)
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

func generateTextImage(w int, h int, str string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	face := basicfont.Face7x13
	lines := strings.Split(str, "\n")
	totalTxtHeight := face.Height * len(lines)
	startY := (h-totalTxtHeight)/2 + face.Height

	for _, line := range lines {
		txtWidth := face.Width * len(line)
		startX := (w - txtWidth) / 2

		drawer := font.Drawer{
			Dst:  img,
			Src:  image.White,
			Face: face,
			Dot:  fixed.P(startX, startY),
		}
		drawer.DrawString(line)
		startY += face.Height
	}
	return img
}
