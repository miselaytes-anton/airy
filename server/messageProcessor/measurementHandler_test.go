package messageprocessor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	models "github.com/miselaytes-anton/tatadata/server/models"
)

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
