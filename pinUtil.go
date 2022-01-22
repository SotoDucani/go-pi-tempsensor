package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
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
