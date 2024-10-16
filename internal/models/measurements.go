package models

import (
	"database/sql"

	"github.com/lib/pq"
)

type MeasurementModelInterface interface {
	GetMeasurements(MeasurementsQuery) ([]Measurement, error)
	InsertMeasurement(Measurement) (string, error)
}

// Measurement represents a single measurement.
type Measurement struct {
	ID          string  `json:"id,omitempty"`
	Timestamp   int64   `json:"timestamp"`
	SensorID    string  `json:"sensorId"`
	IAQ         float64 `json:"iaq"`
	CO2         float64 `json:"co2"`
	VOC         float64 `json:"voc"`
	Pressure    float64 `json:"pressure"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

// MeasurementsQuery represents a query for measurements.
type MeasurementsQuery struct {
	StartEpoch, EndEpoch int64
	Resolution           int
	SensorIDs            []string
}

// MeasurementModel represents a measurement model.
type MeasurementModel struct {
	DB *sql.DB
}

// InsertMeasurement inserts a new measurement into the database.
func (m MeasurementModel) InsertMeasurement(measurement Measurement) (string, error) {
	query := `insert into "measurements"("timestamp", "sensor_id", "iaq",  "co2", "voc", "pressure", "temperature", "humidity") values($1, $2, $3, $4, $5, $6, $7, $8)`
	err := m.DB.QueryRow(
		query,
		measurement.Timestamp,
		measurement.SensorID,
		measurement.IAQ,
		measurement.CO2,
		measurement.VOC,
		measurement.Pressure,
		measurement.Temperature,
		measurement.Humidity,
	).Scan(&measurement)

	if err != nil {
		return "", err
	}

	return measurement.ID, nil
}

// GetMeasurements returns measurements aggregated by resolution (ms) between fromEpoch and toEpoch.
func (m MeasurementModel) GetMeasurements(mq MeasurementsQuery) ([]Measurement, error) {
	query := `
	select
	(floor("timestamp"/$1)*$1)::numeric::integer as timestamp, 
	sensor_id, 
	avg(iaq) as iaq, 
	avg(humidity) as humidity,
	avg(temperature) as temperature,
	avg(pressure) as pressure, 
	avg(co2) as co2, 
	avg(voc) as voc
	from "measurements"
	where sensor_id = any($4) and "timestamp" >= $2 and "timestamp" <= $3
	group by (floor("timestamp"/$1)*$1)::numeric::integer, sensor_id
	order by timestamp asc
	`

	rows, err := m.DB.Query(query, mq.Resolution, mq.StartEpoch, mq.EndEpoch, pq.Array(mq.SensorIDs))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	measurements := make([]Measurement, 0)

	for rows.Next() {
		var measurement Measurement
		err := rows.Scan(&measurement.Timestamp, &measurement.SensorID, &measurement.IAQ, &measurement.Humidity, &measurement.Temperature, &measurement.Pressure, &measurement.CO2, &measurement.VOC)
		if err != nil {
			return nil, err
		}
		measurements = append(measurements, measurement)
	}

	return measurements, nil
}
