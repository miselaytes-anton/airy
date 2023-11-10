package main

import (
	"net/http"

	models "github.com/miselaytes-anton/tatadata/backend/internal/models"
)

// ServerEnv represents the environment containing server dependencies.
type Server struct {
	Router interface {
		HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	}
	Measurements interface {
		GetMeasurements(mq models.MeasurementsQuery) ([]models.Measurement, error)
	}
	Events interface {
		GetEvents(q models.EventsQuery) ([]models.Event, error)
		InsertEvent(e models.Event) (bool, error)
	}
}

// StartServer starts the http server.
func (s Server) Routes() {
	s.Router.HandleFunc("/api/graphs", s.handleGraphs())
	s.Router.HandleFunc("/api/events", s.handleEvents())
	s.Router.HandleFunc("/api/measurements", s.handleMeasurements())
}
