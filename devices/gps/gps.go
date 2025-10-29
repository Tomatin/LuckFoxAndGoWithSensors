package gps

import (
	"LuckFoxGo/io/uart"
	"LuckFoxGo/timeout"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RMC struct {
	Time      time.Time
	Latitude  float64
	Longitude float64
	Knots     float64
	Course    float64
	Valid     bool
}

const (
	GPS_UART        = "/dev/ttyS3"
	GPS_BAUD_RATE   = 9600
	NMEA_RMC_HEADER = "$GPRMC"
)

var (
	uart_gps           *uart.UART
	ENABLE_RMC_FRAME   = []byte{0xB5, 0x62, 0x06, 0x01, 0x03, 0x00, 0xF0, 0x04, 0x01, 0xF5, 0x0E}
	DISABLE_ALL_FRAMES = [][]byte{
		{0xB5, 0x62, 0x06, 0x01, 0x03, 0x00, 0xF0, 0x00, 0x00, 0xFA, 0x0F}, // GGA
		{0xB5, 0x62, 0x06, 0x01, 0x03, 0x00, 0xF0, 0x01, 0x00, 0xFB, 0x11}, // GLL
		{0xB5, 0x62, 0x06, 0x01, 0x03, 0x00, 0xF0, 0x02, 0x00, 0xFC, 0x13}, // GSA
		{0xB5, 0x62, 0x06, 0x01, 0x03, 0x00, 0xF0, 0x03, 0x00, 0xFD, 0x15}, // GSV
		{0xB5, 0x62, 0x06, 0x01, 0x03, 0x00, 0xF0, 0x05, 0x00, 0xFF, 0x19}, // VTG
		{0xB5, 0x62, 0x06, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // TXT
	}
)

func Atach() {
	// Wait two seconds to init GPS
	time.Sleep(2 * time.Second)

	var err error

	// Init uart
	if uart_gps, err = uart.Open(GPS_UART, GPS_BAUD_RATE); err != nil {
		panic(err)
	}

	// Disable all NMEA frames except RMC
	for _, msg := range DISABLE_ALL_FRAMES {
		uart_gps.WriteBytes(msg)
		time.Sleep(100 * time.Millisecond)
	}

	// Just in case enable the RMC frame
	uart_gps.WriteBytes(ENABLE_RMC_FRAME)

	uart_gps.Flush()
}

func Detach() {
	uart_gps.Close()
}

func GpsGetRMCFrame() (*RMC, error) {
	// Get NMEA RMC frame
	buf := make([]byte, 0, 256)

	// Init timeout
	timeout := timeout.TimerEventSet(30 * time.Second)

	for {
		if timeout.TimerEventHasExpired() {
			return nil, fmt.Errorf("Timeout")
		}

		// Read until new line char
		tmp, n, err := read_one_byte()

		if err != nil {
			return nil, err
		}

		if n == 0 {
			continue
		}

		b := tmp

		if b == '\r' {
			continue
		}

		if b == '\n' {
			break
		}

		buf = append(buf, b)
	}

	// Build fields from the entire frame
	sbuf := string(buf)
	frame := strings.Split(sbuf, ",")

	rmc := new(RMC)

	// Check if frame is RMC
	if frame[0] != NMEA_RMC_HEADER {
		return rmc, nil
	}

	// frame[ 0] => RMC header
	// frame[ 1] => UTC time
	// frame[ 2] => Status
	// frame[ 3] => Latitude
	// frame[ 4] => North/South indicator
	// frame[ 5] => Longitude
	// frame[ 6] => East/West indicator
	// frame[ 7] => Speed over ground
	// frame[ 8] => Course over ground
	// frame[ 9] => Date
	// frame[10] =>
	// frame[11] =>
	// frame[12] => Checksum

	// Verify checksum
	if frame[12][2:] == fmt.Sprintf("%02X", nmea_checksum(sbuf)) {
		// Verify header and valid frame
		if frame[2] == "A" {
			// Get knots
			rmc.Knots = parse_float(frame[7])

			// Get course
			rmc.Course = parse_float(frame[8])

			// Get latitude and longitude
			rmc.Latitude = nmea_to_maps_coordinates(frame[3], frame[4])
			rmc.Longitude = nmea_to_maps_coordinates(frame[5], frame[6])

			// Get date and time
			var err error
			rmc.Time, err = time.Parse("020106 150405", frame[9]+" "+frame[1])
			local_time := time.FixedZone("UTC-3", -3*60*60)
			rmc.Time = rmc.Time.In(local_time)

			rmc.Valid = true

			if err != nil {
				return rmc, err
			}
		}
	}

	return rmc, nil
}

func nmea_checksum(in string) (check_sum int) {
	nmea_data := []byte(in)

	for i := 1; i < len(in)-3; i++ {
		check_sum ^= (int)(nmea_data[i])
	}

	return
}

func nmea_to_maps_coordinates(coord string, direction string) float64 {
	val, _ := strconv.ParseFloat(coord, 64)

	degrees := float64(int(val / 100))
	minutes := val - (degrees * 100)
	decimal := degrees + (minutes / 60.0)

	if direction == "S" || direction == "W" {
		decimal = -decimal
	}

	return decimal
}

func parse_float(in string) (val float64) {
	val, err := strconv.ParseFloat(in, 64)

	if err != nil {
		return 0
	}

	return
}

func read_one_byte() (the_char byte, n int, err error) {
	maxRetries := 5
	tmp := make([]byte, 1)

	for range maxRetries {
		n, err := uart_gps.ReadBytes(tmp)

		if err == nil {
			return tmp[0], n, nil
		}

		time.Sleep(time.Second * 1)
	}

	return 0, 0, err
}
