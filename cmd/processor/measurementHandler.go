package main

import (
	"fmt"
	"log"
	"time"

	"github.com/oddnoddles/airy-backend/internal/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type measurementHandler struct {
	Measurements models.MeasurementModelInterface
	LogError     *log.Logger
	LogInfo      *log.Logger
}

// parseMeasurementMessage parses a measurement message which comes in the form of "bedroom 51.86 607.44 0.52 100853 27.25 60.22"
func parseMeasurementMessage(msg string) (models.Measurement, error) {
	var m models.Measurement
	if _, err := fmt.Sscanf(msg, "%s %g %g %g %g %g %g", &m.SensorID, &m.IAQ, &m.CO2, &m.VOC, &m.Pressure, &m.Temperature, &m.Humidity); err != nil {
		return m, err
	}

	return m, nil
}

func (h measurementHandler) handle(_ mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	h.LogInfo.Printf("received message: %s\n", payload)
	m, err := parseMeasurementMessage(payload)
	if err != nil {
		h.LogError.Printf("message could not be parsed (%s): %s", payload, err)
		return
	}

	m.Timestamp = time.Now().Unix()

	h.LogInfo.Printf("inserting measurement: %+v\n", m)

	_, err = h.Measurements.InsertMeasurement(m)
	if err != nil {
		h.LogError.Printf("measurement could not be inserted into database: %s", err)
	}
}
