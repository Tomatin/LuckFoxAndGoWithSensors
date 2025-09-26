package leds

import (
	"LuckFoxGo/io/gpio"
)

const (
	LED_A = "GPIO1_C7"
)

var (
	led_a *gpio.GPIO
)

func Atach() {
	var err error

	if led_a, err = gpio.Atach(LED_A); err != nil {
		panic("error in LED_A")
	}

	led_a.GPIOAsOutput()
}

func Detach() {
	led_a.GPIOSetValue(false)
	led_a.Detach()
}

func LED_A_On() {
	led_a.GPIOSetValue(true)
}

func LED_A_Off() {
	led_a.GPIOSetValue(false)
}
