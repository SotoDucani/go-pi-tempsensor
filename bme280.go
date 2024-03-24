package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

type Bmxx80Device struct {
	DeviceHandle    bmxx80.Dev
	I2CBus          i2c.BusCloser
	temperatureData physic.Temperature
	pressureData    physic.Pressure
	humidityData    physic.RelativeHumidity
}

func (dev *Bmxx80Device) Init() {
	// Load the required drivers
	_, err := host.Init()
	if err != nil {
		panic(err)
	}

	// Open a handle to the first available I2C bus
	bus, err := i2creg.Open("")
	if err != nil {
		panic(err)
	}

	bmeDev, err := bmxx80.NewI2C(bus, 0x77, &bmxx80.DefaultOpts)
	if err != nil {
		panic(err)
	}

	dev.DeviceHandle = *bmeDev
	dev.I2CBus = bus
}

func (dev *Bmxx80Device) Close() {
	dev.DeviceHandle.Halt()
	dev.I2CBus.Close()
}

func (dev *Bmxx80Device) Run(interval time.Duration) {
	for {
		var envData physic.Env
		err := dev.DeviceHandle.Sense(&envData)
		if err != nil {
			panic(err)
		}
		dev.temperatureData = envData.Temperature
		dev.pressureData = envData.Pressure
		dev.humidityData = envData.Humidity
		fmt.Print(dev.PrintAll())
		time.Sleep(interval)
	}
}

func (dev *Bmxx80Device) PrintTemperature() string {
	return fmt.Sprintf("%8s", dev.temperatureData)
}

func (dev *Bmxx80Device) PrintPressure() string {
	return fmt.Sprintf("%10s", dev.pressureData)
}

func (dev *Bmxx80Device) PrintHumidity() string {
	return fmt.Sprintf("%9s", dev.humidityData)
}

func (dev *Bmxx80Device) PrintAll() string {
	return fmt.Sprintf("%8s %10s %9s\n", dev.temperatureData, dev.pressureData, dev.humidityData)
}
