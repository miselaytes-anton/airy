package mocks

import "github.com/oddnoddles/airy-backend/internal/models"

type InsertMeasurementMock = func(models.Measurement, *[]models.Measurement) (string, error)

type GetMeasurementsMock = func(models.MeasurementsQuery, *[]models.Measurement) ([]models.Measurement, error)

type MeasurementModelMock struct {
	Measurements []models.Measurement
	InsertMeasurementMock
	GetMeasurementsMock
}

func (m *MeasurementModelMock) InsertMeasurement(measurement models.Measurement) (string, error) {
	return m.InsertMeasurementMock(measurement, &m.Measurements)
}

func (m *MeasurementModelMock) GetMeasurements(mq models.MeasurementsQuery) ([]models.Measurement, error) {
	return m.GetMeasurementsMock(mq, &m.Measurements)
}

func GetMeasurementsOkMock(mq models.MeasurementsQuery, measurements *[]models.Measurement) ([]models.Measurement, error) {
	return *measurements, nil
}
