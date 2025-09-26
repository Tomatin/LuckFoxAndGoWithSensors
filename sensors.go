package main

import (
	"LuckFoxGo/backend"
	"LuckFoxGo/devices/aht10"
	"LuckFoxGo/devices/gps"
	"LuckFoxGo/devices/leds"
	"fmt"
	"runtime"
	"time"
)

const (
	PRG_NAME = "LuckFox Pro/Max with sensors"
	VERSION  = "1.0"
)

func blink(on, off int) {
	for {
		leds.LED_A_On()
		time.Sleep(time.Millisecond * time.Duration(on))

		leds.LED_A_Off()
		time.Sleep(time.Millisecond * time.Duration(off))
	}
}

func main() {
	fmt.Println(PRG_NAME, "v", VERSION)
	fmt.Printf("Go version: %s\n\n", runtime.Version())

	// Init devices
	gps.Atach()
	aht10.Atach()
	leds.Atach()

	defer func() {
		fmt.Println("Sensor power down")
		gps.Detach()
		leds.Detach()
		aht10.Detach()
	}()

	// Blink led
	go blink(100, 3000)

	// Init WebServer
	backend.WebServerInit()
}
