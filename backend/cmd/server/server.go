package main

import (
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

// StartServer starts the http server.
func (s Server) Routes() {
	s.Router.HandlerFunc(http.MethodGet, "/api/graphs", s.handleGraphs())
	s.Router.HandlerFunc(http.MethodGet, "/api/events", s.handleEventsList())
	s.Router.HandlerFunc(http.MethodPost, "/api/events", s.handleEventsCreate())
	s.Router.HandlerFunc(http.MethodGet, "/api/measurements", s.handleMeasurements())
}

func (s Server) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	s.LogError.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
