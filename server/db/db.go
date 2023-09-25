package db

import (
	"database/sql"
	"time"
)

// Measurement represents a single measurement.
type Measurement struct {
	Timestamp   time.Time
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
	_, err := db.Exec(query, "now()", measurement.SensorID, measurement.IAQ, measurement.CO2, measurement.VOC, measurement.Pressure, measurement.Temperature, measurement.Humidity)

	if err != nil {
		return false, err
	}

	return true, nil
}
