package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
)

func bme_Loop(bmeDev bmxx80.Dev, envChan chan EnvData) {
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
