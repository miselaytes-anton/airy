package messageprocessor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

	if diff := cmp.Diff(expected, m); diff != "" {
		t.Error(diff)
	}

	_, err = parseMeasurementMessage("")

	if err == nil {
		t.Errorf("err: expected error to not be nil")
	}
}
