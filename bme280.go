package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

func bmeInit() bmxx80.Dev {
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

	// Open a handle to a bmxx80 connected on the I²C bus
	bmeDev, err := bmxx80.NewI2C(bus, 0x77, &bmxx80.DefaultOpts)
	if err != nil {
		panic(err)
	}
	defer bmeDev.Halt()

	return *bmeDev
}

func bme_Loop(envChan chan EnvData) {
	bmeDev := bmeInit()
	var interval time.Duration = 5
	for {
		bme280(bmeDev, envChan)
		time.Sleep(interval * time.Second)
	}
}

func bme280(bmeDev bmxx80.Dev, envChan chan EnvData) {
	var env physic.Env
	err := bmeDev.Sense(&env)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
}
