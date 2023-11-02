package messageprocessor_test

import (
	"fmt"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	messageprocessor "github.com/miselaytes-anton/tatadata/server/messageProcessor"
	"github.com/miselaytes-anton/tatadata/server/models"
)

type MeasurementsStub struct {
	measurements []models.Measurement
}

func (m *MeasurementsStub) InsertMeasurement(measurement models.Measurement) (bool, error) {
	fmt.Println("INSERTING")
	m.measurements = append(m.measurements, measurement)
	return true, nil
}

type MqttClientStub struct {
	mqtt.Client
}

type MessageStub struct {
	mqtt.Message
	payload func() []byte
}

func (m MessageStub) Payload() []byte {
	return m.payload()
}

func TestOnMessageHandler(t *testing.T) {
	data := []struct {
		name     string
		message  string
		expected []models.Measurement
	}{
		{
			"valid message",
			"bedroom 51.86 607.44 0.52 100853 27.25 60.22",
			[]models.Measurement{{
				SensorID:    "bedroom",
				IAQ:         51.86,
				CO2:         607.44,
				VOC:         0.52,
				Pressure:    100853,
				Temperature: 27.25,
				Humidity:    60.22,
			}},
		},
		{
			"empty message",
			"",
			make([]models.Measurement, 0),
		},
		{
			"invalid message",
			"bedroom something",
			make([]models.Measurement, 0),
		},
	}

	for _, d := range data {
		t.Run(
			d.name,
			func(t *testing.T) {
				measurementsStub := MeasurementsStub{
					measurements: make([]models.Measurement, 0),
				}
				handler := messageprocessor.MeasurementHandler{
					Measurements: &measurementsStub,
				}
				messageStub := MessageStub{payload: func() []byte {
					return []byte(d.message)
				}}

				handler.OnMessageHandler(MqttClientStub{}, messageStub)

				if diff := cmp.Diff(d.expected, measurementsStub.measurements, cmpopts.IgnoreFields(models.Measurement{}, "Timestamp")); diff != "" {
					t.Error(diff)
				}

			},
		)
	}
}
