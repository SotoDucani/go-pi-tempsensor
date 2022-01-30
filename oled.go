package main

import (
	"image"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3"
)

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

func oled(envChan chan EnvData) {
	// Load all the drivers
	if _, err := host.Init(); err != nil {
		panic(err)
	}

	// Open a handle to the first available I²C bus
	bus, err := i2creg.Open("")
	if err != nil {
		panic(err)
	}
	defer bus.Close()

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
	defer oledDev.Halt()

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
		imgs[i] = convertAndResizeAndCenter(oledDev.Bounds().Dx(), oledDev.Bounds().Dy(), g.Image[i])
	}

	// Display the frames in a loop:
	for i := 0; ; i++ {
		index := i % len(imgs)
		c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
		img := imgs[index]
		oledDev.Draw(img.Bounds(), img, image.Point{})
		<-c
	}

	/*
		img := image1bit.NewVerticalLSB(oledDev.Bounds())
		f := basicfont.Face7x13
		drawer := font.Drawer{
			Dst:  img,
			Src:  &image.Uniform{image1bit.On},
			Face: f,
			Dot:  fixed.P(0, img.Bounds().Dy()-1-f.Descent),
		}
		drawer.DrawString("Hello World!")
		err = oledDev.Draw(img.Bounds(), img, image.Point{})
		if err != nil {
			panic(err)
		}
	*/
}
