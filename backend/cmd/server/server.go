package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	models "github.com/miselaytes-anton/tatadata/backend/internal/models"
)

// ServerEnv represents the environment containing server dependencies.
type Server struct {
	Router interface {
		HandlerFunc(string, string, http.HandlerFunc)
	}
	Measurements models.MeasurementModelInterface
	Events       models.EventModelInterface
	LogError     *log.Logger
	LogInfo      *log.Logger
}

type ResponseError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// StartServer starts the http server.
func (s Server) routes() {
	s.Router.HandlerFunc(http.MethodGet, "/api/graphs", s.handleGraphs())
	s.Router.HandlerFunc(http.MethodGet, "/api/events", s.handleEventsList())
	s.Router.HandlerFunc(http.MethodPost, "/api/events", s.handleEventsCreate())
	s.Router.HandlerFunc(http.MethodPatch, "/api/events/:id", s.handleEventsUpdate())
	s.Router.HandlerFunc(http.MethodGet, "/api/measurements", s.handleMeasurements())
}

func (s Server) jsonError(w http.ResponseWriter, err error, code int) {
	var responseError ResponseError

	if code == http.StatusInternalServerError {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		s.LogError.Output(2, trace)

		responseError = ResponseError{
			Status: http.StatusText(code),
			Error:  "internal server error occured",
		}
	} else {
		responseError = ResponseError{
			Status: http.StatusText(code),
			Error:  err.Error(),
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(responseError)
}
