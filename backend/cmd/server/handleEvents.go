package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
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
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(events)
		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleEventsCreate() http.HandlerFunc {
	type request struct {
		StartTimestamp int64  `json:"startTimestamp" validate:"required,gte=0,lte=2147483647"`
		EndTimestamp   int64  `json:"endTimestamp,omitempty" validate:"omitempty,gtfield=StartTimestamp,lte=2147483647"`
		LocationID     string `json:"locationId" validate:"required,oneof=bedroom livingroom"`
		EventType      string `json:"eventType" validate:"required"`
	}

	type response = models.Event

	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		var request request
		err := s.readJson(w, r, &request)

		if err != nil {
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		err = validate.Struct(request)

		if err != nil {
			s.jsonValidationError(w, err)
			return
		}

		event := models.Event{
			StartTimestamp: request.StartTimestamp,
			EndTimestamp:   request.EndTimestamp,
			LocationID:     request.LocationID,
			EventType:      request.EventType,
		}

		event, err = s.Events.InsertEvent(event)

		if err != nil {
			if errors.Is(err, models.ErrDuplicateEvent) {
				s.jsonError(w, err, http.StatusConflict)
				return
			}

			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		var response response = event

		err = json.NewEncoder(w).Encode(response)

		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleEventsUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event models.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil || event.EndTimestamp <= 0 {
			err := errors.New("invalid event format, expected endTimestamp in ms")
			s.jsonError(w, err, http.StatusBadRequest)
			return
		}

		params := httprouter.ParamsFromContext(r.Context())

		id := params.ByName("id")

		event, err = s.Events.UpdateEvent(id, event.EndTimestamp)

		if err == sql.ErrNoRows {
			err = errors.New("event not found")
			s.jsonError(w, err, http.StatusNotFound)
			return
		}

		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(event)

		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}
	}
}
