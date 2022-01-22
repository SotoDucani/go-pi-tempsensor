package main

import (
	"fmt"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

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
