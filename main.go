package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/devices/v3/ssd1306"
	"periph.io/x/host/v3"
)

type EnvData struct {
	temp     string
	pressure string
	humidity string
}

func forever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Shutdown signal received...")
}

func deviceInit() (ssd1306.Dev, bmxx80.Dev) {
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

	// Open a handle to a ssd1306 connected on the I²C bus
	oledDev, err := ssd1306.NewI2C(bus, &ssd1306.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer oledDev.Halt()

	// Open a handle to a bmxx80 connected on the I²C bus
	bmeDev, err := bmxx80.NewI2C(bus, 0x77, &bmxx80.DefaultOpts)
	if err != nil {
		panic(err)
	}
	defer bmeDev.Halt()

	return *oledDev, *bmeDev
}

func main() {
	oledDev, bmeDev := deviceInit()
	var envChan chan EnvData

	go oled(oledDev, envChan)
	go bme_Loop(&bmeDev, envChan)

	forever()
}
