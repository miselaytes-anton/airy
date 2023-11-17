package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/julienschmidt/httprouter"

	"github.com/miselaytes-anton/tatadata/backend/internal/models"
	"github.com/miselaytes-anton/tatadata/backend/internal/models/mocks"
	"github.com/miselaytes-anton/tatadata/backend/internal/testserver"
)

func Test_handleEventsList(t *testing.T) {
	events := []models.Event{{
		Timestamp:  1,
		LocationID: "bedroom",
		EventType:  "window:open",
	}}

	eventsMock := mocks.EventModelMock{
		Events:          events,
		GetEventsMock:   mocks.GetEventsOkMock,
		InsertEventMock: mocks.InsertEventOkMock,
	}

	router := httprouter.New()
	server := Server{
		Router:   router,
		Events:   &eventsMock,
		LogError: log.New(io.Discard, "", 0),
		LogInfo:  log.New(io.Discard, "", 0),
	}

	server.routes()

	ts := testserver.TestServer{Server: httptest.NewServer(router)}
	defer ts.Server.Close()

	validRequests := []struct {
		name         string
		urlPath      string
		expectedCode int
	}{
		{
			"valid query",
			"/api/events?from=1&to=2",
			http.StatusOK,
		},
	}

	for _, d := range validRequests {
		t.Run(
			d.name,
			func(t *testing.T) {
				statusCode, _, body := ts.Get(t, d.urlPath)

				receivedEvents := new([]models.Event)

				err := json.Unmarshal(body, &receivedEvents)

				if err != nil {
					log.Fatal(err)
				}

				if diff := cmp.Diff(d.expectedCode, statusCode); diff != "" {
					t.Error(diff)
				}

				if diff := cmp.Diff(events, *receivedEvents); diff != "" {
					t.Error(diff)
				}

			},
		)
	}

	invalidRequests := []struct {
		name          string
		urlPath       string
		expectedCode  int
		expectedError ResponseError
		getEventsMock mocks.GetEventsMock
	}{
		{
			"invalid to",
			"/api/events?from=1&to=hello",
			http.StatusBadRequest,
			ResponseError{
				Status: "Bad Request",
				Error:  "invalid to: hello, must be a unix timestamp in ms",
			},
			mocks.GetEventsOkMock,
		},
		{
			"invalid from",
			"/api/events?from=hello&to=2",
			http.StatusBadRequest,
			ResponseError{
				Status: "Bad Request",
				Error:  "invalid from: hello, must be a unix timestamp in ms",
			},
			mocks.GetEventsOkMock,
		},
		{
			"database error",
			"/api/events?from=1&to=2",
			http.StatusInternalServerError,
			ResponseError{
				Status: "Internal Server Error",
				Error:  "internal server error occured",
			},
			mocks.GetEventsErrorMock,
		},
	}

	for _, d := range invalidRequests {
		t.Run(
			d.name,
			func(t *testing.T) {
				eventsMock.GetEventsMock = d.getEventsMock

				statusCode, _, body := ts.Get(t, d.urlPath)

				if diff := cmp.Diff(d.expectedCode, statusCode); diff != "" {
					t.Error(diff)
				}

				responseError := new(ResponseError)

				err := json.Unmarshal(body, &responseError)

				if err != nil {
					log.Fatal(err)
				}

				if diff := cmp.Diff(d.expectedError, *responseError); diff != "" {
					t.Error(diff)
				}
			},
		)
	}
}

func Test_handleEventsCreate(t *testing.T) {
	event := models.Event{
		Timestamp:  1,
		LocationID: "bedroom",
		EventType:  "window:open",
	}

	eventsMock := mocks.EventModelMock{
		Events:          make([]models.Event, 0),
		GetEventsMock:   mocks.GetEventsOkMock,
		InsertEventMock: mocks.InsertEventOkMock,
	}

	router := httprouter.New()
	server := Server{
		Router:   router,
		Events:   &eventsMock,
		LogError: log.New(io.Discard, "", 0),
		LogInfo:  log.New(io.Discard, "", 0),
	}

	server.routes()

	ts := testserver.TestServer{Server: httptest.NewServer(router)}
	defer ts.Server.Close()

	validRequests := []struct {
		name         string
		urlPath      string
		expectedCode int
	}{
		{
			"valid request",
			"/api/events",
			http.StatusCreated,
		},
	}

	for _, d := range validRequests {
		t.Run(
			d.name,
			func(t *testing.T) {
				b, err := json.Marshal(event)
				if err != nil {
					log.Fatal(err)
				}
				statusCode, _, _ := ts.PostJson(t, d.urlPath, b)

				if diff := cmp.Diff(d.expectedCode, statusCode); diff != "" {
					t.Error(diff)
				}

				if diff := cmp.Diff(event, eventsMock.Events[0]); diff != "" {
					t.Error(diff)
				}
			},
		)
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}
	invalidRequests := []struct {
		name            string
		urlPath         string
		expectedCode    int
		expectedError   ResponseError
		insertEventMock mocks.InsertEventMock
		requestBody     []byte
	}{
		{
			"invalid request",
			"/api/events",
			http.StatusBadRequest,
			ResponseError{
				Status: "Bad Request",
				Error:  "invalid event format, expected timestamp in ms, locationId and eventType",
			},
			mocks.InsertEventOkMock,
			make([]byte, 0),
		},
		{
			"database error",
			"/api/events",
			http.StatusInternalServerError,
			ResponseError{
				Status: "Internal Server Error",
				Error:  "internal server error occured",
			},
			mocks.InsertEventErrorMock,
			eventBytes,
		},
	}

	for _, d := range invalidRequests {
		t.Run(
			d.name,
			func(t *testing.T) {
				eventsMock.InsertEventMock = d.insertEventMock
				statusCode, _, body := ts.PostJson(t, d.urlPath, d.requestBody)

				if diff := cmp.Diff(d.expectedCode, statusCode); diff != "" {
					t.Error(diff)
				}

				responseError := new(ResponseError)

				err = json.Unmarshal(body, &responseError)

				if err != nil {
					log.Fatal(err)
				}

				if diff := cmp.Diff(d.expectedError, *responseError); diff != "" {
					t.Error(diff)
				}
			},
		)
	}
}
