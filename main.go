package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
)

func forever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Shutdown signal received...")
}

func stats_Loop() {
	var interval time.Duration = 5
	for {
		bme280()
		time.Sleep(interval * time.Second)
	}
}

func main() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer rpio.Close()
	go oled()
	go stats_Loop()

	forever()
}
