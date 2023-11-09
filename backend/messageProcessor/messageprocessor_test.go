package messageprocessor_test

import (
	"errors"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	messageprocessor "github.com/miselaytes-anton/tatadata/server/messageProcessor"
	"github.com/miselaytes-anton/tatadata/server/models"
)

type insertMeasurement = func(m models.Measurement, measurements *[]models.Measurement) (bool, error)

type MeasurementsStub struct {
	measurements []models.Measurement
	insertMeasurement
}

func (m *MeasurementsStub) InsertMeasurement(measurement models.Measurement) (bool, error) {
	return m.insertMeasurement(measurement, &m.measurements)
}

func insertMeasurementOk(m models.Measurement, measurements *[]models.Measurement) (bool, error) {
	*measurements = append(*measurements, m)
	return true, nil
}

func insertMeasurementError(m models.Measurement, measurements *[]models.Measurement) (bool, error) {
	return false, errors.New("database error")
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
		insertMeasurement
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
			insertMeasurementOk,
		},
		{
			"empty message",
			"",
			make([]models.Measurement, 0),
			insertMeasurementOk,
		},
		{
			"invalid message",
			"bedroom something",
			make([]models.Measurement, 0),
			insertMeasurementOk,
		},
		{
			"valid message, database error",
			"bedroom 51.86 607.44 0.52 100853 27.25 60.22",
			make([]models.Measurement, 0),
			insertMeasurementError,
		},
	}

	for _, d := range data {
		t.Run(
			d.name,
			func(t *testing.T) {
				measurementsStub := MeasurementsStub{
					measurements:      make([]models.Measurement, 0),
					insertMeasurement: d.insertMeasurement,
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
