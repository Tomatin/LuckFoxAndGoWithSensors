package backend

import (
	"LuckFoxGo/devices/aht10"
	"LuckFoxGo/devices/gps"
	"encoding/json"
	"log"
	"net/http"
)

const (
	PORT = ":8080"
)

type JSONResponse struct {
	Temperature float32
	Humidity    float32
	Time        string
	Latitude    float64
	Longitude   float64
	Knots       float64
	Course      float64
	Valid       bool
}

func WebServerInit() {
	// Handler for serving the dashboard.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/mnt/sdcard/static/dashboard.html")
	})

	// Handler for serving JSON data
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get GPS data
		rmc, err := gps.GpsGetRMCFrame()

		if err != nil {
			panic(err)
		}

		// Get sensor data
		sensor, err := aht10.GetSensorMeasurement()

		if err != nil {
			panic(err)
		}

		// Build JSON response
		jsonData := JSONResponse{
			Temperature: sensor.Temperature,
			Humidity:    sensor.Humidity,
			Time:        rmc.Time.Format("2006-01-02 15:04:05"),
			Latitude:    rmc.Latitude,
			Longitude:   rmc.Longitude,
			Knots:       rmc.Knots,
			Course:      rmc.Course,
			Valid:       rmc.Valid,
		}

		json.NewEncoder(w).Encode(jsonData)
	})

	// Launch WebServer
	err := http.ListenAndServe(PORT, nil)

	if err != nil {
		log.Fatal(err)
	}
}
