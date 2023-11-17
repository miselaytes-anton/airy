package main

import (
	"encoding/json"
	"errors"
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

func (s *Server) handleEventsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := makeEventsQuery(r)
		if err != nil {
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		events, err := s.Events.GetEvents(q)
		if err != nil {
			s.jsonServerError(w, err)
			return
		}

		err = json.NewEncoder(w).Encode(events)
		if err != nil {
			s.jsonServerError(w, err)
			return
		}
	}
}

func (s *Server) handleEventsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event models.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			err := errors.New("invalid event format, expected timestamp in ms, locationId and eventType")
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		_, err = s.Events.InsertEvent(event)
		if err != nil {
			s.jsonServerError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
