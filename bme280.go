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
	DeviceHandle    *bmxx80.Dev
	I2CBus          i2c.BusCloser
	temperatureData float64
	pressureData    physic.Pressure
	humidityData    physic.RelativeHumidity
	TempUnits       string
}

func (dev *Bmxx80Device) Init(TempUnits string) {
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

	dev.DeviceHandle = bmeDev
	dev.I2CBus = bus
	dev.TempUnits = TempUnits
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
		switch dev.TempUnits {
		case "Fahrenheit":
			dev.temperatureData = envData.Temperature.Fahrenheit()
		case "Celsius":
			dev.temperatureData = envData.Temperature.Celsius()
		default:
			dev.temperatureData = envData.Temperature.Fahrenheit()
		}

		dev.pressureData = envData.Pressure
		dev.humidityData = envData.Humidity
		fmt.Print(dev.PrintAll())
		time.Sleep(interval)
	}
}

func (dev *Bmxx80Device) PrintTemperature() string {
	switch dev.TempUnits {
	case "Fahrenheit":
		return fmt.Sprintf("%.1f°F", dev.temperatureData)
	case "Celsius":
		return fmt.Sprintf("%.1f°C", dev.temperatureData)
	default:
		return fmt.Sprintf("%.1f°F", dev.temperatureData)
	}
}

func (dev *Bmxx80Device) PrintPressure() string {
	return fmt.Sprintf("%10s", dev.pressureData)
}

func (dev *Bmxx80Device) PrintHumidity() string {
	return fmt.Sprintf("%9s", dev.humidityData)
}

func (dev *Bmxx80Device) PrintAll() string {

	return fmt.Sprintf("%s %s %s\n", dev.PrintTemperature(), dev.PrintPressure(), dev.PrintHumidity())
}
