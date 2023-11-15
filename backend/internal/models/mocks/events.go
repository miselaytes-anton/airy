package mocks

import "github.com/miselaytes-anton/tatadata/backend/internal/models"

type InsertEventMock = func(m models.Event, measurements *[]models.Event) (bool, error)

type GetEventsMock = func(mq models.EventsQuery, measurements *[]models.Event) ([]models.Event, error)

type EventModelMock struct {
	Events []models.Event
	InsertEventMock
	GetEventsMock
}

func (m *EventModelMock) InsertEvent(measurement models.Event) (bool, error) {
	return m.InsertEventMock(measurement, &m.Events)
}

func (m *EventModelMock) GetEvents(mq models.EventsQuery) ([]models.Event, error) {
	return m.GetEventsMock(mq, &m.Events)
}
