package mocks

import (
	"errors"

	"github.com/miselaytes-anton/tatadata/backend/internal/models"
)

type InsertEventMock = func(models.Event, *[]models.Event) (string, error)
type GetEventsMock = func(models.EventsQuery, *[]models.Event) ([]models.Event, error)
type UpdateEventMock = func(string, int64, *[]models.Event) (models.Event, error)

type EventModelMock struct {
	Events []models.Event
	InsertEventMock
	GetEventsMock
	UpdateEventMock
}

func (m *EventModelMock) InsertEvent(event models.Event) (string, error) {
	return m.InsertEventMock(event, &m.Events)
}

func (m *EventModelMock) UpdateEvent(id string, endTimestamp int64) (models.Event, error) {
	return m.UpdateEventMock(id, endTimestamp, &m.Events)
}

func (m *EventModelMock) GetEvents(mq models.EventsQuery) ([]models.Event, error) {
	return m.GetEventsMock(mq, &m.Events)
}

func GetEventsOkMock(mq models.EventsQuery, events *[]models.Event) ([]models.Event, error) {
	return *events, nil
}

func GetEventsErrorMock(mq models.EventsQuery, events *[]models.Event) ([]models.Event, error) {
	return nil, errors.New("database error")
}

func InsertEventOkMock(m models.Event, events *[]models.Event) (string, error) {
	*events = append(*events, m)
	return "uuid", nil
}

func InsertEventErrorMock(m models.Event, events *[]models.Event) (string, error) {
	return "", errors.New("database error")
}
