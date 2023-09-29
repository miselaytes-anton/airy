package api

import (
	"database/sql"
	"fmt"
	"net/http"
	database "server/db"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type valueGetter func(measurement database.Measurement) float64

func generateLineItemsFromMeasurements(measurements []database.Measurement, getValue valueGetter) []opts.LineData {
	items := make([]opts.LineData, 0)

	for _, measurement := range measurements {
		items = append(items, opts.LineData{
			Value: []interface{}{
				time.Unix(measurement.Timestamp, 0),
				getValue(measurement),
			},
		})
	}

	return items
}

func makeChart(sensorID string, lineItems []opts.LineData, title string, startEpoch int64, endEpoch int64) *charts.Line {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros, PageTitle: "Graphs"}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time",
			Type: "time",
			Min:  time.Unix(startEpoch, 0),
			Max:  time.Unix(endEpoch, 0),
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis", TriggerOn: "click"}),
	)

	// Put data into instance
	line.AddSeries(sensorID, lineItems).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	return line
}

func getStartOfDay(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)
}

func getEndOfDay(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, location)
}

func parseParams(r *http.Request) (database.MeasurementsQuery, error) {
	amsterdam, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		return database.MeasurementsQuery{}, err
	}
	dateStr := r.URL.Query().Get("date")

	var date time.Time
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return database.MeasurementsQuery{}, err
		}
	} else {
		date = time.Now().In(amsterdam)
	}

	startEpoch := getStartOfDay(date, amsterdam).Unix()
	endEpoch := getEndOfDay(date, amsterdam).Unix()

	resolutionStr := r.URL.Query().Get("resolution")
	var resolution int
	if resolutionStr == "" {
		resolution = 3600
	} else {
		res, err := strconv.ParseInt(resolutionStr, 10, 32)
		if err != nil {
			return database.MeasurementsQuery{}, err
		}
		resolution = int(res)
	}

	sensorID := r.URL.Query().Get("sensor_id")
	if sensorID == "" {
		sensorID = "livingroom"
	}

	return database.MeasurementsQuery{
		StartEpoch: startEpoch,
		EndEpoch:   endEpoch,
		Resolution: resolution,
		SensorID:   sensorID,
	}, nil

}

func graphsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := parseParams(r)
		if err != nil {
			fmt.Fprintf(w, "Error parsing params: %s", err)
			return
		}

		fmt.Printf("Getting measurements for: %+v\n", params)

		measurements, err := database.GetMeasurements(db, params)
		if err != nil {
			fmt.Fprintf(w, "Error getting measurements: %s", err)
			return
		}

		iaqLineItems := generateLineItemsFromMeasurements(measurements, func(m database.Measurement) float64 { return m.IAQ })
		iaqChart := makeChart(params.SensorID, iaqLineItems, "IAQ", params.StartEpoch, params.EndEpoch)
		iaqChart.Render(w)

		humidityLineItems := generateLineItemsFromMeasurements(measurements, func(m database.Measurement) float64 { return m.Humidity })
		humidityChart := makeChart(params.SensorID, humidityLineItems, "Humidity", params.StartEpoch, params.EndEpoch)
		humidityChart.Render(w)

		temperatureLineItems := generateLineItemsFromMeasurements(measurements, func(m database.Measurement) float64 { return m.Temperature })
		temperatureChart := makeChart(params.SensorID, temperatureLineItems, "Temperature", params.StartEpoch, params.EndEpoch)
		temperatureChart.Render(w)
	}
}

// StartServer starts the http server.
func StartServer(db *sql.DB) {
	http.HandleFunc("/graphs", graphsHandler(db))
	http.ListenAndServe(":8081", nil)
}
