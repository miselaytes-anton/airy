package mocks

import (
	"errors"

	"github.com/miselaytes-anton/tatadata/backend/internal/models"
)

type InsertEventMock = func(models.Event, *[]models.Event) (models.Event, error)
type GetEventsMock = func(models.EventsQuery, *[]models.Event) ([]models.Event, error)
type UpdateEventMock = func(string, int64, *[]models.Event) (models.Event, error)

type EventModelMock struct {
	Events []models.Event
	InsertEventMock
	GetEventsMock
	UpdateEventMock
}

func (m *EventModelMock) InsertEvent(event models.Event) (models.Event, error) {
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

func InsertEventOkMock(e models.Event, events *[]models.Event) (models.Event, error) {
	e.ID = "uuid"
	*events = append(*events, e)
	return e, nil
}

func InsertEventErrorMock(e models.Event, events *[]models.Event) (models.Event, error) {
	return models.Event{}, errors.New("database error")
}
