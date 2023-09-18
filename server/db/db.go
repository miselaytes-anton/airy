package db

import (
	"database/sql"
	"time"
)

// Measurement represents a single measurement.
type Measurement struct {
	Timestamp   time.Time
	SensorID    string
	Temperature uint64
	Humidity    uint64
	CO2         uint64
	VOC         uint64
}

// InsertMeasurement inserts a new measurement into the database.
func InsertMeasurement(db *sql.DB, measurement Measurement) (bool, error) {
	query := `insert into "measurements"("timestamp", "sensor_id", "temperature", "humidity", "co2", "voc") values($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, "now()", measurement.SensorID, measurement.Temperature, measurement.Humidity, measurement.CO2, measurement.VOC)

	if err != nil {
		return false, err
	}

	return true, nil
}
