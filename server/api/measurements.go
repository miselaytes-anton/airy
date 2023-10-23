package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/miselaytes-anton/tatadata/server/models"
)

func makeMeasurementsQuery(r *http.Request) (models.MeasurementsQuery, error) {
	q := models.MeasurementsQuery{}

	fromStr := r.URL.Query().Get("from")
	fromEpoch, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		return q, fmt.Errorf("invalid from: %s, must be a unix timestamp in ms", fromStr)
	}

	toStr := r.URL.Query().Get("to")
	toEpoch, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		return q, fmt.Errorf("invalid to: %s, must be a unix timestamp in ms", toStr)
	}

	resolutionStr := r.URL.Query().Get("resolution")
	resolution, err := strconv.ParseInt(resolutionStr, 10, 64)
	if err != nil {
		return q, fmt.Errorf("invalid resolution: %s, must be an integer", resolutionStr)
	}

	q.StartEpoch = fromEpoch
	q.EndEpoch = toEpoch
	q.Resolution = int(resolution)
	q.SensorIDs = sensorIDs

	return q, nil
}

func measurementsHandler(env *ServerEnv) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			q, err := makeMeasurementsQuery(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			measurements, err := env.Measurements.GetMeasurements(q)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(measurements)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	}
}
