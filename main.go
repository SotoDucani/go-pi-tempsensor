package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func forever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutdown signal received...")
}

func updateTempDisplay(interval int, bmxx *Bmxx80Device, oled *OledDevice) {
	for {
		txt := fmt.Sprintf("%s %s\n%s", bmxx.PrintTemperature(), bmxx.humidityData, bmxx.pressureData)
		oled.DisplayText(txt)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func main() {
	var Bmxx Bmxx80Device
	var Oled OledDevice
	var Prom PrometheusMetrics
	desiredTempUnits := "Fahrenheit"

	// Setup the BME280
	Bmxx.Init(desiredTempUnits)
	defer Bmxx.Close()
	go Bmxx.Run(5 * time.Second)

	// Setup the OLED screen
	Oled.InitDefault()
	defer Oled.Close()
	//go Oled.DisplayGif("./ballerine.gif")
	//Oled.DisplayText("Hello World!")
	time.Sleep(time.Second)
	go updateTempDisplay(5, &Bmxx, &Oled)

	// Setup the Prometheus metrics
	Prom.Init(desiredTempUnits)
	go ServePromServer(&Prom)

	// Run the app forever
	forever()
}
