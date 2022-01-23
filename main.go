package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type EnvData struct {
	temp     string
	pressure string
	humidity string
}

func forever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Shutdown signal received...")
}

func main() {
	var envChan chan EnvData

	go oled(envChan)
	go bme_Loop(envChan)

	forever()
}
