package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/miselaytes-anton/tatadata/backend/internal/models"
)

func makeEventsQuery(r *http.Request) (models.EventsQuery, error) {
	q := models.EventsQuery{}

	fromStr := r.URL.Query().Get("from")
	fromEpoch, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		return q, fmt.Errorf("invalid from: %s, must be a unix timestamp in ms", fromStr)
	}

	toStr := r.URL.Query().Get("to")
	toEpoch, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		return q, fmt.Errorf("invalid to: %s, must be a unix timestamp in ms", toStr)
	}

	q.StartEpoch = fromEpoch
	q.EndEpoch = toEpoch

	return q, nil
}

// Handles a POST request to /events by inserting event into the database.
// Also handles a GET request to /events by returning events between fromEpoch and toEpoch.
func (s *Server) handleEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var event models.Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			_, err = s.Events.InsertEvent(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
		case "GET":
			q, err := makeEventsQuery(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			events, err := s.Events.GetEvents(q)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(events)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	}
}
