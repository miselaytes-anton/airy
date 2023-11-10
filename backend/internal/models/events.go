package models

import (
	"database/sql"
)

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
	Timestamp  int64  `json:"timestamp"`
	LocationID string `json:"locationId"`
	EventType  string `json:"eventType"`
}

// GetEvents returns events between fromEpoch and toEpoch.
func (m EventModel) GetEvents(q EventsQuery) ([]Event, error) {
	query := `
	select timestamp, location_id, type from "events"
	where "timestamp" >= $1 and "timestamp" <= $2
	order by timestamp asc
	`

	rows, err := m.DB.Query(query, q.StartEpoch, q.EndEpoch)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := make([]Event, 0)

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Timestamp, &event.LocationID, &event.EventType)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// InsertEvent inserts a new event into the database.
func (m EventModel) InsertEvent(e Event) (bool, error) {
	query := `insert into "events"("timestamp", "location_id", "type") values($1, $2, $3)`
	_, err := m.DB.Exec(query, e.Timestamp, e.LocationID, e.EventType)

	if err != nil {
		return false, err
	}

	return true, nil
}
