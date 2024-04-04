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

func updateTempDisplay(interval time.Duration, bmxx *Bmxx80Device, oled *OledDevice) {
	for {
		txt := fmt.Sprintf("%s %s\n%s", bmxx.PrintTemperature(), bmxx.humidityData, bmxx.pressureData)
		oled.DisplayText(txt)
		time.Sleep(interval * time.Second)
	}
}

func updatePromMetrics(interval time.Duration, bmxx *Bmxx80Device, em *PrometheusMetrics) {
	for {
		em.temperature.Set(bmxx.temperatureData)
		em.humidity.Set(float64(bmxx.humidityData))
		em.pressure.Set(float64(bmxx.pressureData))
		time.Sleep(interval)
	}
}

func main() {
	var Bmxx Bmxx80Device
	var Oled OledDevice
	var Prom PrometheusMetrics
	desiredTempUnits := "Fahrenheit"
	desiredUpdateIntervals := time.Duration(5)

	// Setup and start running the BME280
	Bmxx.Init(desiredTempUnits)
	defer Bmxx.Close()
	go Bmxx.Run(desiredUpdateIntervals * time.Second)

	// Setup and start running the OLED screen
	Oled.InitDefault()
	defer Oled.Close()
	// Sleep for a second so our first display update has data
	time.Sleep(time.Second)
	go updateTempDisplay(desiredUpdateIntervals, &Bmxx, &Oled)

	// Setup and start running the Prometheus metrics
	Prom.Init(desiredTempUnits)
	go ServePromServer(&Prom)
	go updatePromMetrics(desiredUpdateIntervals, &Bmxx, &Prom)

	// Run the app forever
	forever()
}
