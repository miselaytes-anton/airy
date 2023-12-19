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

func Test_handleGraphs(t *testing.T) {
	events := []models.Event{{
		StartTimestamp: 1,
		LocationID:     "bedroom",
		EventType:      "window:open",
	}}

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

	eventsMock := mocks.EventModelMock{
		Events:     events,
		GetAllMock: mocks.GetAllEventsOkMock,
	}

	measurementsMock := mocks.MeasurementModelMock{
		Measurements:        measurements,
		GetMeasurementsMock: mocks.GetMeasurementsOkMock,
	}

	router := httprouter.New()
	server := Server{
		Router:       router,
		Events:       &eventsMock,
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
			"no query",
			"/api/graphs",
			http.StatusOK,
		},
		{
			"day view",
			"/api/graphs?view=day",
			http.StatusOK,
		},
		{
			"day view, date",
			"/api/graphs?view=day&date=2020-01-01",
			http.StatusOK,
		},
		{
			"day view, date, resolution",
			"/api/graphs?view=day&date=2020-01-01&resolution=600",
			http.StatusOK,
		},
		{
			"week view",
			"/api/graphs?view=week",
			http.StatusOK,
		},
		{
			"week view, date",
			"/api/graphs?view=week&date=2020-01-01",
			http.StatusOK,
		},
		{
			"week view, date, resolution",
			"/api/graphs?view=week&date=2020-01-01&resolution=86400",
			http.StatusOK,
		},
		{
			"invalid view",
			"/api/graphs?view=month",
			http.StatusBadRequest,
		},
		{
			"invalid date",
			"/api/graphs?date=123",
			http.StatusBadRequest,
		},
		{
			"invalid resolution:string",
			"/api/graphs?resolution=hello",
			http.StatusBadRequest,
		},
		{
			"invalid resolution:0",
			"/api/graphs?resolution=0",
			http.StatusBadRequest,
		},
		{
			"invalid resolution:too big",
			"/api/graphs?resolution=99999999",
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
