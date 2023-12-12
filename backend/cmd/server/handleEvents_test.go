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
		StartTimestamp: 1,
		LocationID:     "bedroom",
		EventType:      "window:open",
	}}

	eventsMock := mocks.EventModelMock{
		Events:          events,
		GetMock:         mocks.GetOkMock,
		GetAllMock:      mocks.GetAllOkMock,
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
		getEventsMock mocks.GetAllMock
	}{
		{
			"invalid to",
			"/api/events?from=1&to=hello",
			http.StatusBadRequest,
			ResponseError{
				Status: "Bad Request",
				Error:  "invalid to: hello, must be a unix timestamp in ms",
			},
			mocks.GetAllOkMock,
		},
		{
			"invalid from",
			"/api/events?from=hello&to=2",
			http.StatusBadRequest,
			ResponseError{
				Status: "Bad Request",
				Error:  "invalid from: hello, must be a unix timestamp in ms",
			},
			mocks.GetAllOkMock,
		},
		{
			"database error",
			"/api/events?from=1&to=2",
			http.StatusInternalServerError,
			ResponseError{
				Status: "Internal Server Error",
				Error:  "internal server error occured",
			},
			mocks.GetAllErrorMock,
		},
	}

	for _, d := range invalidRequests {
		t.Run(
			d.name,
			func(t *testing.T) {
				eventsMock.GetAllMock = d.getEventsMock

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
	type Response = models.Event
	type Request struct {
		StartTimestamp int64  `json:"startTimestamp"`
		LocationID     string `json:"locationId"`
		EventType      string `json:"eventType"`
	}

	event := models.Event{
		ID:             "uuid",
		StartTimestamp: 1,
		LocationID:     "bedroom",
		EventType:      "window:open",
	}

	request := Request{
		StartTimestamp: event.StartTimestamp,
		LocationID:     event.LocationID,
		EventType:      event.EventType,
	}

	eventsMock := mocks.EventModelMock{
		Events:          make([]models.Event, 0),
		GetAllMock:      mocks.GetAllOkMock,
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
			http.StatusOK,
		},
	}

	for _, d := range validRequests {
		t.Run(
			d.name,
			func(t *testing.T) {
				b, err := json.Marshal(request)
				if err != nil {
					log.Fatal(err)
				}
				statusCode, _, body := ts.PostJson(t, d.urlPath, b)

				if diff := cmp.Diff(d.expectedCode, statusCode); diff != "" {
					t.Error(diff)
				}

				if diff := cmp.Diff(event, eventsMock.Events[0]); diff != "" {
					t.Error(diff)
				}

				response := new(Response)

				err = json.Unmarshal(body, &response)

				if err != nil {
					log.Fatal(err)
				}

				if diff := cmp.Diff(*response, event); diff != "" {
					t.Error(diff)
				}
			},
		)
	}

	requestBytes, err := json.Marshal(request)
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
				Error:  "body must not be empty",
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
			requestBytes,
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
