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

	Bmxx.Init()
	defer Bmxx.Close()
	go Bmxx.Run(5 * time.Second)

	Oled.InitDefault()
	defer Oled.Close()
	go Oled.DisplayGif("./ballerine.gif")

	forever()
}
