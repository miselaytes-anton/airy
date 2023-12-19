package mocks

import (
	"errors"

	"github.com/miselaytes-anton/tatadata/internal/models"
)

type InsertEventMock = func(models.Event, *[]models.Event) (models.Event, error)
type GetAllMock = func(models.EventsQuery, *[]models.Event) ([]models.Event, error)
type GetMock = func(string, *[]models.Event) (models.Event, error)
type UpdateEventMock = func(models.Event, *[]models.Event) (models.Event, error)

type EventModelMock struct {
	Events []models.Event
	InsertEventMock
	GetAllMock
	GetMock
	UpdateEventMock
}

func (m *EventModelMock) InsertEvent(event models.Event) (models.Event, error) {
	return m.InsertEventMock(event, &m.Events)
}

func (m *EventModelMock) UpdateEvent(e models.Event) (models.Event, error) {
	return m.UpdateEventMock(e, &m.Events)
}

func (m *EventModelMock) Get(id string) (models.Event, error) {
	return m.GetMock(id, &m.Events)
}

func (m *EventModelMock) GetAll(mq models.EventsQuery) ([]models.Event, error) {
	return m.GetAllMock(mq, &m.Events)
}

func GetEventOkMock(id string, events *[]models.Event) (models.Event, error) {
	for _, event := range *events {
		if event.ID == id {
			return event, nil
		}
	}
	return models.Event{}, nil
}

func GetAllEventsOkMock(mq models.EventsQuery, events *[]models.Event) ([]models.Event, error) {
	return *events, nil
}

func GetAllEventsErrorMock(mq models.EventsQuery, events *[]models.Event) ([]models.Event, error) {
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

func UpdateEventOkMock(e models.Event, events *[]models.Event) (models.Event, error) {
	return e, nil
}
