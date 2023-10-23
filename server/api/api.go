package api

import (
	"net/http"

	models "github.com/miselaytes-anton/tatadata/server/models"
)

// ServerEnv represents the environment containing server dependencies.
type ServerEnv struct {
	Measurements interface {
		GetMeasurements(mq models.MeasurementsQuery) ([]models.Measurement, error)
	}
	Events interface {
		GetEvents(q models.EventsQuery) ([]models.Event, error)
		InsertEvent(e models.Event) (bool, error)
	}
}

// StartServer starts the http server.
func StartServer(env *ServerEnv) {
	http.HandleFunc("/api/graphs", graphsHandler(env))
	http.HandleFunc("/api/events", eventsHandler(env))
	http.HandleFunc("/api/measurements", measurementsHandler(env))
	http.ListenAndServe(":8081", nil)
}
