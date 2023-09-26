package api

import (
	"database/sql"
	"fmt"
	"net/http"
	database "server/db"
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

func graphsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {

		//todo: get start and end epoch from query params
		startEpoch := time.Now().Unix()/86400*86400 - 3600*2
		endEpoch := startEpoch + 86400
		fmt.Println(startEpoch, endEpoch, time.Now().Unix())
		var resolution int64 = 3600
		sensorID := "livingroom"
		measurements, err := database.GetMeasurements(db, resolution, startEpoch, endEpoch, sensorID)
		if err != nil {
			fmt.Fprintf(w, "Error getting measurements: %s", err)
			return
		}

		iaqLineItems := generateLineItemsFromMeasurements(measurements, func(m database.Measurement) float64 { return m.IAQ })
		iaqChart := makeChart(sensorID, iaqLineItems, "IAQ", startEpoch, endEpoch)
		iaqChart.Render(w)

		humidityLineItems := generateLineItemsFromMeasurements(measurements, func(m database.Measurement) float64 { return m.Humidity })
		humidityChart := makeChart(sensorID, humidityLineItems, "Humidity", startEpoch, endEpoch)
		humidityChart.Render(w)

		temperatureLineItems := generateLineItemsFromMeasurements(measurements, func(m database.Measurement) float64 { return m.Temperature })
		temperatureChart := makeChart(sensorID, temperatureLineItems, "Temperature", startEpoch, endEpoch)
		temperatureChart.Render(w)
	}
}

// StartServer starts the http server.
func StartServer(db *sql.DB) {
	http.HandleFunc("/graphs", graphsHandler(db))
	http.ListenAndServe(":8081", nil)
}
