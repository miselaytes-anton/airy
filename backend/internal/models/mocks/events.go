package mocks

import "github.com/miselaytes-anton/tatadata/backend/internal/models"

type InsertEventMock = func(models.Event, *[]models.Event) (bool, error)

type GetEventsMock = func(models.EventsQuery, *[]models.Event) ([]models.Event, error)

type EventModelMock struct {
	Events []models.Event
	InsertEventMock
	GetEventsMock
}

func (m *EventModelMock) InsertEvent(event models.Event) (bool, error) {
	return m.InsertEventMock(event, &m.Events)
}

func (m *EventModelMock) GetEvents(mq models.EventsQuery) ([]models.Event, error) {
	return m.GetEventsMock(mq, &m.Events)
}
