package messageprocessor

import (
	"testing"

	models "github.com/miselaytes-anton/tatadata/server/models"
)

func Test_parseMeasurementMessage(t *testing.T) {
	m, err := parseMeasurementMessage("bedroom 51.86 607.44 0.52 100853 27.25 60.22")

	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	expected := models.Measurement{
		SensorID:    "bedroom",
		IAQ:         51.86,
		CO2:         607.44,
		VOC:         0.52,
		Pressure:    100853,
		Temperature: 27.25,
		Humidity:    60.22,
	}

	const float64EqualityThreshold = 1e-9

	if m.CO2-expected.CO2 > float64EqualityThreshold {
		t.Errorf("CO2: expected %f, got %f", expected.CO2, m.CO2)
	}

	if m.IAQ-expected.IAQ > float64EqualityThreshold {
		t.Errorf("IAQ: expected %f, got %f", expected.IAQ, m.IAQ)
	}

	if m.VOC-expected.VOC > float64EqualityThreshold {
		t.Errorf("VOC: expected %f, got %f", expected.VOC, m.VOC)
	}

	if m.Pressure != expected.Pressure {
		t.Errorf("Pressure: expected %f, got %f", expected.Pressure, m.Pressure)
	}

	if m.Temperature-expected.Temperature > float64EqualityThreshold {
		t.Errorf("Temperature: expected %f, got %f", expected.Temperature, m.Temperature)
	}

	if m.Humidity-expected.Humidity > float64EqualityThreshold {
		t.Errorf("Humidity: expected %f, got %f", expected.Humidity, m.Humidity)
	}

	_, err = parseMeasurementMessage("")

	if err == nil {
		t.Errorf("err: expected error to no be nil")
	}
}
