package api

import (
	"fmt"
	"net/http"
	models "server/models"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type valueGetter func(measurement models.Measurement) float64
type lineItemsPerSensor map[string][]opts.LineData
type measurementsPerSensor map[string][]models.Measurement

var sensorIDs = []string{"bedroom", "livingroom"}

func generateLineItemsFromMeasurements(measurementsPerSensor measurementsPerSensor, getValue valueGetter) lineItemsPerSensor {
	items := make(lineItemsPerSensor)

	for sensorID, measurements := range measurementsPerSensor {
		for _, measurement := range measurements {
			items[sensorID] = append(items[sensorID], opts.LineData{Value: []interface{}{time.Unix(measurement.Timestamp, 0), getValue(measurement)}})
		}
	}

	return items
}

func makeChart(items lineItemsPerSensor, title string, startEpoch int64, endEpoch int64) *charts.Line {
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

	// Create line graphs for each sensor with keys ordered alphabetically
	for _, sensorID := range sensorIDs {
		line.AddSeries(sensorID, items[sensorID]).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	}

	return line
}

func getStartOfDay(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, location)
}

func getEndOfDay(t time.Time, location *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, location)
}

func parseParams(r *http.Request) (models.MeasurementsQuery, error) {
	amsterdam, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		return models.MeasurementsQuery{}, err
	}
	dateStr := r.URL.Query().Get("date")

	var date time.Time
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return models.MeasurementsQuery{}, err
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
			return models.MeasurementsQuery{}, err
		}
		resolution = int(res)
	}

	return models.MeasurementsQuery{
		StartEpoch: startEpoch,
		EndEpoch:   endEpoch,
		Resolution: resolution,
		SensorIDs:  sensorIDs,
	}, nil

}

func graphsHandler(env *ServerEnv) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := parseParams(r)
		if err != nil {
			fmt.Fprintf(w, "Error parsing params: %s", err)
			return
		}

		fmt.Printf("Getting measurements for: %+v\n", params)

		measurements, err := env.Measurements.GetMeasurements(params)
		if err != nil {
			fmt.Fprintf(w, "Error getting measurements: %s", err)
			return
		}

		measurementsPerSensor := make(measurementsPerSensor)

		for _, measurement := range measurements {
			measurementsPerSensor[measurement.SensorID] = append(measurementsPerSensor[measurement.SensorID], measurement)
		}

		iaqLineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.IAQ })
		iaqChart := makeChart(iaqLineItems, "IAQ", params.StartEpoch, params.EndEpoch)
		iaqChart.Render(w)

		humidityLineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.Humidity })
		humidityChart := makeChart(humidityLineItems, "Humidity", params.StartEpoch, params.EndEpoch)
		humidityChart.Render(w)

		temperatureLineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.Temperature })
		temperatureChart := makeChart(temperatureLineItems, "Temperature", params.StartEpoch, params.EndEpoch)
		temperatureChart.Render(w)
	}
}

// ServerEnv represents the environment containing server dependencies.
type ServerEnv struct {
	Measurements interface {
		GetMeasurements(mq models.MeasurementsQuery) ([]models.Measurement, error)
	}
}

// StartServer starts the http server.
func StartServer(env *ServerEnv) {
	http.HandleFunc("/graphs", graphsHandler(env))
	http.ListenAndServe(":8081", nil)
}
