package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/julienschmidt/httprouter"

	"github.com/oddnoddles/airy-backend/internal/models"
	"github.com/oddnoddles/airy-backend/internal/models/mocks"
	"github.com/oddnoddles/airy-backend/internal/testserver"
)

func Test_handleMeasurements(t *testing.T) {
	measurements := []models.Measurement{{
		Timestamp:   1,
		SensorID:    "bedroom",
		IAQ:         150,
		CO2:         900,
		VOC:         6,
		Pressure:    760,
		Temperature: 20,
		Humidity:    50,
	}}

	measurementsMock := mocks.MeasurementModelMock{
		Measurements:        measurements,
		GetMeasurementsMock: mocks.GetMeasurementsOkMock,
	}

	router := httprouter.New()
	server := Server{
		Router:       router,
		Measurements: &measurementsMock,
		LogError:     log.New(io.Discard, "", 0),
		LogInfo:      log.New(io.Discard, "", 0),
	}

	server.routes()

	ts := testserver.TestServer{Server: httptest.NewServer(router)}
	defer ts.Server.Close()

	requests := []struct {
		name         string
		urlPath      string
		expectedCode int
	}{
		{
			"valid query",
			"/api/measurements?from=1&to=2&resolution=600",
			http.StatusOK,
		},
		{
			"invalid from",
			"/api/measurements?from=hello&to=2&resolution=600",
			http.StatusBadRequest,
		},
		{
			"invalid to",
			"/api/measurements?from=1&to=hello&resolution=600",
			http.StatusBadRequest,
		},
		{
			"invalid resolution",
			"/api/measurements?from=1&to=-2&resolution=hello",
			http.StatusBadRequest,
		},
	}

	for _, d := range requests {
		t.Run(
			d.name,
			func(t *testing.T) {
				statusCode, _, _ := ts.Get(t, d.urlPath)
				if diff := cmp.Diff(d.expectedCode, statusCode); diff != "" {
					t.Error(diff)
				}
			},
		)
	}
}
