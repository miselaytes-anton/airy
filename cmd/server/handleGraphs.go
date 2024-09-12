package main

import (
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/go-playground/validator/v10"

	"github.com/miselaytes-anton/airy/internal/dateutil"
	"github.com/miselaytes-anton/airy/internal/models"
	"github.com/miselaytes-anton/airy/internal/urlquery"
)

type valueGetter func(measurement models.Measurement) float64
type lineItemsPerSensor map[string][]opts.LineData
type markLinesPerSensor map[string][]opts.MarkLineNameXAxisItem
type measurementsPerSensor map[string][]models.Measurement
type eventsPerSensor map[string][]models.Event
type viewConfig struct {
	resolution       int
	startEpochOffset int64
	endEpochOffset   int64
}

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

func generateMarkLinesFromEvents(eventsPerSensor eventsPerSensor) markLinesPerSensor {
	items := make(markLinesPerSensor)

	for sensorID, events := range eventsPerSensor {
		for _, event := range events {
			items[sensorID] = append(items[sensorID], opts.MarkLineNameXAxisItem{Name: event.EventType, XAxis: time.Unix(event.StartTimestamp, 0)})
		}
	}

	return items
}

func makeChart(items lineItemsPerSensor, markLines markLinesPerSensor, title string, startEpoch int64, endEpoch int64) *charts.Line {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros, PageTitle: "Graphs", Width: "100%"}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		// disable animation
		charts.WithAnimation(),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time",
			Type: "time",
			Min:  time.Unix(startEpoch, 0),
			Max:  time.Unix(endEpoch, 0),
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis", TriggerOn: "click"}),
	)

	// Create line graphs for each sensor with keys ordered alphabetically
	for _, sensorID := range SENSOR_IDS {
		seriesOptions := []charts.SeriesOpts{
			charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
			charts.WithMarkLineStyleOpts(opts.MarkLineStyle{Symbol: []string{"none"}, Label: &opts.Label{Show: true, Formatter: "{b}"}}),
		}
		for _, markLine := range markLines[sensorID] {
			seriesOptions = append(seriesOptions, charts.WithMarkLineNameXAxisItemOpts(markLine))
		}

		line.AddSeries(sensorID, items[sensorID]).
			SetSeriesOptions(
				seriesOptions...,
			)
	}

	return line
}

type graphsQuery struct {
	View       *string `validate:"omitempty,oneof=day week"`
	Date       *time.Time
	Resolution *int `validate:"omitempty,gt=0,lte=86400"`
}

// parseGraphsQuery parses the query parameters for the graphs endpoint.
// if a parameter is not present nil is returned.
func parseGraphsQuery(r *http.Request) (*graphsQuery, error) {
	values := r.URL.Query()

	view := urlquery.ReadStringFromQuery(values, "view")

	resolution, err := urlquery.ReadIntFromQuery(values, "resolution")
	if err != nil {
		return nil, err
	}

	date, err := urlquery.ReadDateFromQuery(values, "date", "2006-01-02")

	if err != nil {
		return nil, err
	}

	return &graphsQuery{
		View:       view,
		Date:       date,
		Resolution: resolution,
	}, nil
}

// getEpochs returns the start and end epoch for the given date and view.
func getEpochs(date time.Time, view string, now time.Time, location time.Location) (int64, int64) {
	var startEpoch, endEpoch int64

	switch view {
	case "day":
		// Show 23 hours before and 1 hour after now.
		if dateutil.IsDateEqual(date, now) {
			startEpoch = now.Unix() - defaultsPerView["day"].startEpochOffset
			endEpoch = now.Unix() + defaultsPerView["day"].endEpochOffset
		} else {
			// Show calendar day.
			startEpoch = dateutil.GetStartOfDay(date, &location).Unix()
			endEpoch = dateutil.GetEndOfDay(date, &location).Unix()
		}
	case "week":
		// Show last 7 calendar days before the date and 1 hour after.
		startEpoch = dateutil.GetStartOfDay(date, &location).Unix() - defaultsPerView["week"].startEpochOffset
		endEpoch = dateutil.GetEndOfDay(date, &location).Unix()
	}

	return startEpoch, endEpoch
}

// makeModelsQueries returns the models.MeasurementsQuery and models.EventsQuery for the given graphsQuery.
func makeModelsQueries(q graphsQuery, now time.Time, location time.Location) (models.MeasurementsQuery, models.EventsQuery) {
	var startEpoch, endEpoch int64
	var view string
	var date time.Time
	var resolution int

	if q.View == nil {
		view = "day"
	} else {
		view = *q.View
	}

	if q.Date == nil {
		date = now
	} else {
		date = *q.Date
	}

	if q.Resolution == nil {
		resolution = defaultsPerView[view].resolution
	} else {
		resolution = *q.Resolution
	}

	startEpoch, endEpoch = getEpochs(date, view, now, location)

	return models.MeasurementsQuery{
			StartEpoch: startEpoch,
			EndEpoch:   endEpoch,
			Resolution: resolution,
			SensorIDs:  SENSOR_IDS,
		}, models.EventsQuery{
			StartEpoch: startEpoch,
			EndEpoch:   endEpoch,
		}
}

func renderGraphs(w http.ResponseWriter, measurements []models.Measurement, events []models.Event, startEpoch int64, endEpoch int64) {
	measurementsPerSensor := make(measurementsPerSensor)

	for _, measurement := range measurements {
		measurementsPerSensor[measurement.SensorID] = append(measurementsPerSensor[measurement.SensorID], measurement)
	}

	eventsPerSensor := make(eventsPerSensor)
	for _, event := range events {
		eventsPerSensor[event.LocationID] = append(eventsPerSensor[event.LocationID], event)
	}
	markLinesPerSensor := generateMarkLinesFromEvents(eventsPerSensor)

	co2LineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.CO2 })
	co2Chart := makeChart(co2LineItems, markLinesPerSensor, "CO2", startEpoch, endEpoch)
	co2Chart.Render(w)

	vocLineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.VOC })
	vocChart := makeChart(vocLineItems, markLinesPerSensor, "VOC", startEpoch, endEpoch)
	vocChart.Render(w)

	iaqLineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.IAQ })
	iaqChart := makeChart(iaqLineItems, markLinesPerSensor, "IAQ", startEpoch, endEpoch)
	iaqChart.Render(w)

	humidityLineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.Humidity })
	humidityChart := makeChart(humidityLineItems, markLinesPerSensor, "Humidity", startEpoch, endEpoch)
	humidityChart.Render(w)

	temperatureLineItems := generateLineItemsFromMeasurements(measurementsPerSensor, func(m models.Measurement) float64 { return m.Temperature })
	temperatureChart := makeChart(temperatureLineItems, markLinesPerSensor, "Temperature", startEpoch, endEpoch)
	temperatureChart.Render(w)
}

func (s *Server) handleGraphs() http.HandlerFunc {

	validate := validator.New(validator.WithRequiredStructEnabled())

	// todo: preload
	amsterdam, loadLocationErr := time.LoadLocation("Europe/Amsterdam")

	return func(w http.ResponseWriter, r *http.Request) {
		if loadLocationErr != nil {
			s.jsonError(w, loadLocationErr, http.StatusInternalServerError)
			return
		}

		var now = time.Now().In(amsterdam)

		graphsQuery, err := parseGraphsQuery(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = validate.Struct(graphsQuery)

		if err != nil {
			s.jsonValidationError(w, err)
			return
		}

		measurementsQuery, eventsQuery := makeModelsQueries(*graphsQuery, now, *amsterdam)

		measurements, err := s.Measurements.GetMeasurements(measurementsQuery)
		if err != nil {
			s.jsonError(w, err, http.StatusInternalServerError)
			return
		}

		events, err := s.Events.GetAll(eventsQuery)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderGraphs(w, measurements, events, measurementsQuery.StartEpoch, measurementsQuery.EndEpoch)
	}
}
