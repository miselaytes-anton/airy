package processor

import (
	"fmt"
	"time"

	"github.com/miselaytes-anton/tatadata/backend/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MeasurementHandler struct defines the handlers for measurment message.
type MeasurementHandler struct {
	Measurements interface {
		InsertMeasurement(mq models.Measurement) (bool, error)
	}
}

// parseMeasurementMessage parses a measurement message which comes in the form of "bedroom 51.86 607.44 0.52 100853 27.25 60.22"
func parseMeasurementMessage(msg string) (models.Measurement, error) {
	var m models.Measurement
	if _, err := fmt.Sscanf(msg, "%s %g %g %g %g %g %g", &m.SensorID, &m.IAQ, &m.CO2, &m.VOC, &m.Pressure, &m.Temperature, &m.Humidity); err != nil {
		return m, err
	}

	return m, nil
}

// OnMessageHandler is called when a message is received
func (h MeasurementHandler) OnMessageHandler(_ mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	fmt.Printf("Received message: %s\n", payload)
	m, err := parseMeasurementMessage(payload)
	if err != nil {
		fmt.Printf("Message could not be parsed (%s): %s", payload, err)
		return
	}

	m.Timestamp = time.Now().Unix()

	fmt.Printf("Inserting measurement: %+v\n", m)

	_, err = h.Measurements.InsertMeasurement(m)
	if err != nil {
		fmt.Printf("Measurement could not be inserted into database: %s", err)
	}
}
