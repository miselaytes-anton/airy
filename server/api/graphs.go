package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/miselaytes-anton/tatadata/server/models"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type valueGetter func(measurement models.Measurement) float64
type lineItemsPerSensor map[string][]opts.LineData
type measurementsPerSensor map[string][]models.Measurement
type viewConfig struct {
	resolution       int
	startEpochOffset int64
	endEpochOffset   int64
}

var sensorIDs = []string{"livingroom", "bedroom"}

var defaultsPerView = map[string]viewConfig{
	"day": {
		resolution:       600,
		startEpochOffset: 23 * 3600,
		endEpochOffset:   3600 - 1,
	},
	"week": {
		resolution:       3600,
		startEpochOffset: 24 * 3600 * 7,
		endEpochOffset:   3600 - 1,
	},
}

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

func dateEqual(date1, date2 time.Time) bool {

	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

func makeMeasurementsQueryFromGetGraphsRequest(r *http.Request) (models.MeasurementsQuery, error) {
	amsterdam, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		return models.MeasurementsQuery{}, err
	}

	view := r.URL.Query().Get("view")
	if view == "" {
		view = "day"
	} else if view != "day" && view != "week" {
		return models.MeasurementsQuery{}, fmt.Errorf("unknown view: %s, can be 'day' or 'week'", view)
	}

	config := defaultsPerView[view]

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

	var startEpoch, endEpoch int64
	var now = time.Now().In(amsterdam)

	// if today
	// then show last 24 hours
	if view == "day" && dateEqual(date, now) {
		startEpoch = now.Unix() - config.startEpochOffset
		endEpoch = now.Unix() + config.endEpochOffset
	} else if view == "day" {
		// if not today and day view
		// then show calendar day
		startEpoch = getStartOfDay(date, amsterdam).Unix()
		endEpoch = getEndOfDay(date, amsterdam).Unix()
	} else if view == "week" {
		// if not today and week view
		// then show last 7 days
		startEpoch = getStartOfDay(date, amsterdam).Unix() - config.startEpochOffset
		endEpoch = getEndOfDay(date, amsterdam).Unix()
	}

	resolutionStr := r.URL.Query().Get("resolution")
	var resolution int
	if resolutionStr == "" {
		resolution = config.resolution
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
		params, err := makeMeasurementsQueryFromGetGraphsRequest(r)
		if err != nil {
			fmt.Fprintf(w, "error parsing params: %s", err)
			return
		}

		fmt.Printf("Getting measurements for: %+v\n", params)

		measurements, err := env.Measurements.GetMeasurements(params)
		if err != nil {
			fmt.Fprintf(w, "error getting measurements: %s", err)
			return
		}

		measurementsPerSensor := make(measurementsPerSensor)

		for _, measurement := range measurements {
			measurementsPerSensor[measurement.SensorID] = append(measurementsPerSensor[measurement.SensorID], measurement)
		}

		co2LineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.CO2 })
		co2Chart := makeChart(co2LineItems, "CO2", params.StartEpoch, params.EndEpoch)
		co2Chart.Render(w)

		vocLineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.VOC })
		vocChart := makeChart(vocLineItems, "VOC", params.StartEpoch, params.EndEpoch)
		vocChart.Render(w)

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
