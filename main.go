package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
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

	dev, err := bmxx80.NewI2C(bus, 0x76, &bmxx80.DefaultOpts)
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

func forever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Shutdown signal received...")
}

func stats_Loop() {
	var interval time.Duration = 20
	for {
		//getAllPinState()
		bme280()
		//flipPinLoop()
		time.Sleep(interval * time.Second)
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
