package models

import (
	"database/sql"
)

type EventModelInterface interface {
	GetEvents(q EventsQuery) ([]Event, error)
	InsertEvent(e Event) (string, error)
	UpdateEvent(id string, endTimestamp int64) (Event, error)
}

// EventModel represents an event model.
type EventModel struct {
	DB *sql.DB
}

// EventsQuery represents a query for measurements.
type EventsQuery struct {
	StartEpoch, EndEpoch int64
}

// Event represents a single event.
type Event struct {
	ID             string `json:"id"`
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
	LocationID     string `json:"locationId"`
	EventType      string `json:"eventType"`
}

// GetEvents returns events between fromEpoch and toEpoch.
func (m EventModel) GetEvents(q EventsQuery) ([]Event, error) {
	query := `
	select id, start_timestamp, end_timestamp, location_id, type from "events"
	where "start_timestamp" >= $1 and "start_timestamp" <= $2
	order by start_timestamp asc
	`

	rows, err := m.DB.Query(query, q.StartEpoch, q.EndEpoch)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := make([]Event, 0)

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.StartTimestamp, &event.EndTimestamp, &event.LocationID, &event.EventType)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// InsertEvent inserts a new event into the database.
func (m EventModel) InsertEvent(e Event) (string, error) {
	query := `insert into "events"("start_timestamp", "end_timestamp", "location_id", "type") values($1, $2, $3, $4) RETURNING id`
	err := m.DB.QueryRow(query, e.StartTimestamp, e.EndTimestamp, e.LocationID, e.EventType).Scan(&e.ID)

	if err != nil {
		return "", err
	}

	return e.ID, nil
}

func (m EventModel) UpdateEvent(id string, endTimestamp int64) (Event, error) {
	e := Event{}
	query := `update "events" set "end_timestamp" = $2 where "id" = $1 
			returning id,start_timestamp,end_timestamp,location_id,type`
	err := m.DB.QueryRow(
		query,
		id,
		endTimestamp,
	).Scan(&e.ID, &e.StartTimestamp, &e.EndTimestamp, &e.LocationID, &e.EventType)

	if err != nil {
		return Event{}, err
	}

	return e, nil
}
