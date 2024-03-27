package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func forever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Shutdown signal received...")
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
	Oled.DisplayText("Hello World!")

	// Setup the Prometheus metrics
	Prom.Init(desiredTempUnits)
	go ServePromServer(&Prom)

	// Run the app forever
	forever()
}
