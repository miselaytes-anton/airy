package mocks

import "github.com/miselaytes-anton/tatadata/backend/internal/models"

type InsertMeasurementMock = func(models.Measurement, *[]models.Measurement) (bool, error)

type GetMeasurementsMock = func(models.MeasurementsQuery, *[]models.Measurement) ([]models.Measurement, error)

type MeasurementModelMock struct {
	Measurements []models.Measurement
	InsertMeasurementMock
	GetMeasurementsMock
}

func (m *MeasurementModelMock) InsertMeasurement(measurement models.Measurement) (bool, error) {
	return m.InsertMeasurementMock(measurement, &m.Measurements)
}

func (m *MeasurementModelMock) GetMeasurements(mq models.MeasurementsQuery) ([]models.Measurement, error) {
	return m.GetMeasurementsMock(mq, &m.Measurements)
}
