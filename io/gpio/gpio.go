package gpio

import (
	"fmt"
	"os"
	"path/filepath"
)

// An abstraction layer for handling a
// general purpose io in LuckFox board

const (
	ROOT_GPIO_SYSFS = "/sys/class/gpio"
	EXPORT_GPIO     = ROOT_GPIO_SYSFS + "/export"
	UNEXPORT_GPIO   = ROOT_GPIO_SYSFS + "/unexport"
	FPERM           = 0644
)

var (
	gpios = map[string]int{
		"GPIO0_A0": 0, "GPIO0_A1": 1, "GPIO0_A2": 2, "GPIO0_A3": 3, "GPIO0_A4": 4, "GPIO0_A5": 5, "GPIO0_A6": 6, "GPIO0_A7": 7,
		"GPIO0_B0": 8, "GPIO0_B1": 9, "GPIO0_B2": 10, "GPIO0_B3": 11, "GPIO0_B4": 12, "GPIO0_B5": 13, "GPIO0_B6": 14, "GPIO0_B7": 15,
		"GPIO0_C0": 16, "GPIO0_C1": 17, "GPIO0_C2": 18, "GPIO0_C3": 19, "GPIO0_C4": 20, "GPIO0_C5": 21, "GPIO0_C6": 22, "GPIO0_C7": 23,
		"GPIO0_D0": 24, "GPIO0_D1": 25, "GPIO0_D2": 26, "GPIO0_D3": 27, "GPIO0_D4": 28, "GPIO0_D5": 29, "GPIO0_D6": 30, "GPIO0_D7": 31,

		"GPIO1_A0": 32, "GPIO1_A1": 33, "GPIO1_A2": 34, "GPIO1_A3": 35, "GPIO1_A4": 36, "GPIO1_A5": 37, "GPIO1_A6": 38, "GPIO1_A7": 39,
		"GPIO1_B0": 40, "GPIO1_B1": 41, "GPIO1_B2": 42, "GPIO1_B3": 43, "GPIO1_B4": 44, "GPIO1_B5": 45, "GPIO1_B6": 46, "GPIO1_B7": 47,
		"GPIO1_C0": 48, "GPIO1_C1": 49, "GPIO1_C2": 50, "GPIO1_C3": 51, "GPIO1_C4": 52, "GPIO1_C5": 53, "GPIO1_C6": 54, "GPIO1_C7": 55,
		"GPIO1_D0": 56, "GPIO1_D1": 57, "GPIO1_D2": 58, "GPIO1_D3": 59, "GPIO1_D4": 60, "GPIO1_D5": 61, "GPIO1_D6": 62, "GPIO1_D7": 63,

		"GPIO2_A0": 64, "GPIO2_A1": 65, "GPIO2_A2": 66, "GPIO2_A3": 67, "GPIO2_A4": 68, "GPIO2_A5": 69, "GPIO2_A6": 70, "GPIO2_A7": 71,
		"GPIO2_B0": 72, "GPIO2_B1": 73, "GPIO2_B2": 74, "GPIO2_B3": 75, "GPIO2_B4": 76, "GPIO2_B5": 77, "GPIO2_B6": 78, "GPIO2_B7": 79,
		"GPIO2_C0": 80, "GPIO2_C1": 81, "GPIO2_C2": 82, "GPIO2_C3": 83, "GPIO2_C4": 84, "GPIO2_C5": 85, "GPIO2_C6": 86, "GPIO2_C7": 87,
		"GPIO2_D0": 88, "GPIO2_D1": 89, "GPIO2_D2": 90, "GPIO2_D3": 91, "GPIO2_D4": 92, "GPIO2_D5": 93, "GPIO2_D6": 94, "GPIO2_D7": 95,

		"GPIO3_A0": 96, "GPIO3_A1": 97, "GPIO3_A2": 98, "GPIO3_A3": 99, "GPIO3_A4": 100, "GPIO3_A5": 101, "GPIO3_A6": 102, "GPIO3_A7": 103,
		"GPIO3_B0": 104, "GPIO3_B1": 105, "GPIO3_B2": 106, "GPIO3_B3": 107, "GPIO3_B4": 108, "GPIO3_B5": 109, "GPIO3_B6": 110, "GPIO3_B7": 111,
		"GPIO3_C0": 112, "GPIO3_C1": 113, "GPIO3_C2": 114, "GPIO3_C3": 115, "GPIO3_C4": 116, "GPIO3_C5": 117, "GPIO3_C6": 118, "GPIO3_C7": 119,
		"GPIO3_D0": 120, "GPIO3_D1": 121, "GPIO3_D2": 122, "GPIO3_D3": 123, "GPIO3_D4": 124, "GPIO3_D5": 125, "GPIO3_D6": 126, "GPIO3_D7": 127,

		"GPIO4_A0": 128, "GPIO4_A1": 129, "GPIO4_A2": 130, "GPIO4_A3": 131, "GPIO4_A4": 132, "GPIO4_A5": 133, "GPIO4_A6": 134, "GPIO4_A7": 135,
		"GPIO4_B0": 136, "GPIO4_B1": 137, "GPIO4_B2": 138, "GPIO4_B3": 139, "GPIO4_B4": 140, "GPIO4_B5": 141, "GPIO4_B6": 142, "GPIO4_B7": 143,
		"GPIO4_C0": 144, "GPIO4_C1": 145, "GPIO4_C2": 146, "GPIO4_C3": 147, "GPIO4_C4": 148, "GPIO4_C5": 149, "GPIO4_C6": 150, "GPIO4_C7": 151,
		"GPIO4_D0": 152, "GPIO4_D1": 153, "GPIO4_D2": 154, "GPIO4_D3": 155, "GPIO4_D4": 156, "GPIO4_D5": 157, "GPIO4_D6": 158, "GPIO4_D7": 159,
	}
)

type GPIO struct {
	path string
	pin  string
}

func Atach(pin string) (*GPIO, error) {
	// Export pin list
	f, err := os.OpenFile(EXPORT_GPIO, os.O_WRONLY, FPERM)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	// Check if the pin name is valid
	pin_name, ok := gpios[pin]

	if !ok {
		return nil, fmt.Errorf("unrecognized pin name")
	}

	pin_number := fmt.Sprintf("%d", pin_name)
	f.WriteString(pin_number)

	// Build pin name path
	gpio := new(GPIO)
	gpio.path = filepath.Join(ROOT_GPIO_SYSFS, "gpio")
	gpio.path += pin_number

	return gpio, nil
}

func (gpio *GPIO) Detach() {
	// Unexport pin list
	f, err := os.OpenFile(UNEXPORT_GPIO, os.O_WRONLY, FPERM)

	if err != nil {
		panic(err)
	}

	f.WriteString(gpio.pin)
	f.Close()
}

func (gpio *GPIO) GPIOAsOutput() {
	os.WriteFile(filepath.Join(gpio.path, "direction"), []byte("out"), FPERM)
}

func (gpio *GPIO) GPIOAsInput() {
	os.WriteFile(filepath.Join(gpio.path, "direction"), []byte("in"), FPERM)
}

func (gpio *GPIO) GPIOSetValue(value bool) {
	os.WriteFile(filepath.Join(gpio.path, "value"), []byte(bool_to_text(value)), FPERM)
}

func (gpio *GPIO) GPIOGetValue() bool {
	n, _ := os.ReadFile(filepath.Join(gpio.path, "value"))

	return string(n) == "1"
}

func bool_to_text(value bool) string {
	if value {
		return "1"
	} else {
		return "0"
	}
}
