package aht10

import (
	"LuckFoxGo/io/i2c"
	"fmt"
	"time"
)

type Sensor struct {
	Temperature float32
	Humidity    float32
}

const (
	I2C_ADDR = 0x38
	BUS      = "/dev/i2c-3"
)

var (
	i2c_sensor *i2c.I2C
)

func Atach() {
	var err error

	// Open i2c conduit
	if i2c_sensor, err = i2c.Open(BUS, I2C_ADDR); err != nil {
		panic(err)
	}

	// Init ATH10
	initCmd := []byte{0xE1, 0x08, 0x00}

	if err := i2c_sensor.Write(initCmd); err != nil {
		panic(fmt.Errorf("error init the ATH10: %w", err))
	}

	time.Sleep(50 * time.Millisecond)
}

func Detach() {
	i2c_sensor.Close()
}

func GetSensorMeasurement() (*Sensor, error) {
	// Fire a measurement
	cmd := []byte{0xAC, 0x33, 0x00}
	if err := i2c_sensor.Write(cmd); err != nil {
		return nil, err
	}

	time.Sleep(100 * time.Millisecond)

	// Read 6 bytes
	buf, err := i2c_sensor.Read(6)
	if err != nil {
		return nil, err
	}

	sensor := new(Sensor)

	// Humidity (20 bits)
	rawHum := (uint32(buf[1]) << 12) | (uint32(buf[2]) << 4) | (uint32(buf[3]) >> 4)
	sensor.Humidity = float32(rawHum) * 100.0 / (1 << 20)

	// Temperature (20 bits)
	rawTemp := ((uint32(buf[3]) & 0x0F) << 16) | (uint32(buf[4]) << 8) | uint32(buf[5])
	sensor.Temperature = (float32(rawTemp) * 200.0 / (1 << 20)) - 50.0

	return sensor, nil
}
