package main

import (
	"errors"
	"io"
	"log"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/miselaytes-anton/tatadata/backend/internal/models"
	"github.com/miselaytes-anton/tatadata/backend/internal/models/mocks"
)

func getMeasurementsOkMock(mq models.MeasurementsQuery, measurements *[]models.Measurement) ([]models.Measurement, error) {
	return *measurements, nil
}

func insertMeasurementOkMock(m models.Measurement, measurements *[]models.Measurement) (string, error) {
	*measurements = append(*measurements, m)
	return "uuid", nil
}

func insertMeasurementErrorMock(m models.Measurement, measurements *[]models.Measurement) (string, error) {
	return "", errors.New("database error")
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

				if diff := cmp.Diff(d.errMsg, errMsg); diff != "" {
					t.Error(diff)
				}
			},
		)
	}

}

func Test_handle(t *testing.T) {
	data := []struct {
		name                  string
		message               string
		expected              []models.Measurement
		insertMeasurementMock mocks.InsertMeasurementMock
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
			insertMeasurementOkMock,
		},
		{
			"empty message",
			"",
			make([]models.Measurement, 0),
			insertMeasurementOkMock,
		},
		{
			"invalid message",
			"bedroom something",
			make([]models.Measurement, 0),
			insertMeasurementOkMock,
		},
		{
			"valid message, database error",
			"bedroom 51.86 607.44 0.52 100853 27.25 60.22",
			make([]models.Measurement, 0),
			insertMeasurementErrorMock,
		},
	}

	for _, d := range data {
		t.Run(
			d.name,
			func(t *testing.T) {
				measurementsMock := mocks.MeasurementModelMock{
					Measurements:          make([]models.Measurement, 0),
					InsertMeasurementMock: d.insertMeasurementMock,
					GetMeasurementsMock:   getMeasurementsOkMock,
				}
				handler := measurementHandler{
					Measurements: &measurementsMock,
					LogError:     log.New(io.Discard, "", 0),
					LogInfo:      log.New(io.Discard, "", 0),
				}
				messageStub := messageStub{payload: func() []byte {
					return []byte(d.message)
				}}

				handler.handle(mqttClientStub{}, messageStub)

				if diff := cmp.Diff(d.expected, measurementsMock.Measurements, cmpopts.IgnoreFields(models.Measurement{}, "Timestamp")); diff != "" {
					t.Error(diff)
				}

			},
		)
	}
}
