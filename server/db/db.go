package db

import (
	"database/sql"
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

// InsertMeasurement inserts a new measurement into the database.
func InsertMeasurement(db *sql.DB, measurement Measurement) (bool, error) {
	query := `insert into "measurements"("timestamp", "sensor_id", "iaq",  "co2", "voc", "pressure", "temperature", "humidity") values($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(query, measurement.Timestamp, measurement.SensorID, measurement.IAQ, measurement.CO2, measurement.VOC, measurement.Pressure, measurement.Temperature, measurement.Humidity)

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetMeasurements returns measurements aggregated by resolution (ms) between fromEpoch and toEpoch.
func GetMeasurements(db *sql.DB, resolution int64, fromEpoch int64, toEpoch int64, sensorID string) ([]Measurement, error) {
	query := `
	select (floor("timestamp"/$1)*$1)::numeric::integer as timestamp, sensor_id, avg(iaq) as iaq
	from "measurements"
	where sensor_id=$4 and "timestamp" >= $2 and "timestamp" <= $3
	group by (floor(timestamp/$1)*$1)::numeric::integer, sensor_id
	`

	rows, err := db.Query(query, resolution, fromEpoch, toEpoch, sensorID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	measurements := make([]Measurement, 0)

	for rows.Next() {
		var measurement Measurement
		err := rows.Scan(&measurement.Timestamp, &measurement.SensorID, &measurement.IAQ)
		if err != nil {
			return nil, err
		}
		measurements = append(measurements, measurement)
	}

	return measurements, nil
}
