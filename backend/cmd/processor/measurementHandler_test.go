package main

import (
	"errors"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/miselaytes-anton/tatadata/backend/internal/log"
	"github.com/miselaytes-anton/tatadata/backend/internal/models"
)

type insertMeasurement = func(m models.Measurement, measurements *[]models.Measurement) (bool, error)

type measurementsStub struct {
	measurements []models.Measurement
	insertMeasurement
}

func (m *measurementsStub) InsertMeasurement(measurement models.Measurement) (bool, error) {
	return m.insertMeasurement(measurement, &m.measurements)
}

func (m *measurementsStub) GetMeasurements(mq models.MeasurementsQuery) ([]models.Measurement, error) {
	return nil, nil
}

func insertMeasurementOk(m models.Measurement, measurements *[]models.Measurement) (bool, error) {
	*measurements = append(*measurements, m)
	return true, nil
}

func insertMeasurementError(m models.Measurement, measurements *[]models.Measurement) (bool, error) {
	return false, errors.New("database error")
}

type mqttClientStub struct {
	mqtt.Client
}

type messageStub struct {
	mqtt.Message
	payload func() []byte
}

func (m messageStub) Payload() []byte {
	return m.payload()
}

func Test_parseMeasurementMessage(t *testing.T) {
	data := []struct {
		name     string
		message  string
		expected models.Measurement
		errMsg   string
	}{
		{
			"valid message",
			"bedroom 51.86 607.44 0.52 100853 27.25 60.22",
			models.Measurement{
				SensorID:    "bedroom",
				IAQ:         51.86,
				CO2:         607.44,
				VOC:         0.52,
				Pressure:    100853,
				Temperature: 27.25,
				Humidity:    60.22,
			},
			"",
		},
		{
			"empty message",
			"",
			models.Measurement{},
			"EOF",
		},
		{
			"invalid message",
			"bedroom something",
			models.Measurement{
				SensorID:    "bedroom",
				IAQ:         0,
				CO2:         0,
				VOC:         0,
				Pressure:    0,
				Temperature: 0,
				Humidity:    0,
			},
			"strconv.ParseFloat: parsing \"\": invalid syntax",
		},
	}

	for _, d := range data {
		t.Run(
			d.name,
			func(t *testing.T) {
				m, err := parseMeasurementMessage(d.message)
				if diff := cmp.Diff(d.expected, m); diff != "" {
					t.Error(diff)
				}

				var errMsg string
				if err != nil {
					errMsg = err.Error()
				}

				if errMsg != d.errMsg {
					t.Errorf("Expected error message `%s`, got `%s`", d.errMsg, errMsg)
				}
			},
		)
	}

}

func Test_handle(t *testing.T) {
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
				measurementsStub := measurementsStub{
					measurements:      make([]models.Measurement, 0),
					insertMeasurement: d.insertMeasurement,
				}
				handler := measurementHandler{
					Measurements: &measurementsStub,
					LogError:     log.Error,
					LogInfo:      log.Info,
				}
				messageStub := messageStub{payload: func() []byte {
					return []byte(d.message)
				}}

				handler.handle(mqttClientStub{}, messageStub)

				if diff := cmp.Diff(d.expected, measurementsStub.measurements, cmpopts.IgnoreFields(models.Measurement{}, "Timestamp")); diff != "" {
					t.Error(diff)
				}

			},
		)
	}
}
