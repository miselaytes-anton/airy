package models

import (
	"database/sql"

	"github.com/lib/pq"
)

// Measurement represents a single measurement.
type Measurement struct {
	Timestamp   int64
	SensorID    string
	IAQ         float64
	CO2         float64
	VOC         float64
	Pressure    float64
	Temperature float64
	Humidity    float64
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
func (m MeasurementModel) InsertMeasurement(measurement Measurement) (bool, error) {
	query := `insert into "measurements"("timestamp", "sensor_id", "iaq",  "co2", "voc", "pressure", "temperature", "humidity") values($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := m.DB.Exec(query, measurement.Timestamp, measurement.SensorID, measurement.IAQ, measurement.CO2, measurement.VOC, measurement.Pressure, measurement.Temperature, measurement.Humidity)

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetMeasurements returns measurements aggregated by resolution (ms) between fromEpoch and toEpoch.
func (m MeasurementModel) GetMeasurements(mq MeasurementsQuery) ([]Measurement, error) {
	query := `
	select (floor("timestamp"/$1)*$1)::numeric::integer as timestamp, 
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
