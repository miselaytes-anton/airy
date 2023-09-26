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

func generateLineItemsFromMeasurements(measurements []database.Measurement) []opts.LineData {
	items := make([]opts.LineData, 0)

	for _, measurement := range measurements {
		items = append(items, opts.LineData{
			Value: []interface{}{
				time.Unix(measurement.Timestamp, 0),
				measurement.IAQ,
			},
		})
	}

	return items
}

func httpserver(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		startEpoch := time.Now().Unix() / 86400 * 86400
		endEpoch := startEpoch + 86400
		var resolution int64 = 3600
		sensorID := "livingroom"
		measurements, err := database.GetMeasurements(db, resolution, startEpoch, endEpoch, sensorID)
		if err != nil {
			fmt.Fprintf(w, "Error getting measurements: %s", err)
			return
		}

		// create a new line instance
		line := charts.NewLine()
		// set some global options like Title/Legend/ToolTip or anything else
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
			charts.WithTitleOpts(opts.Title{
				Title:    "IAQ",
				Subtitle: "Hourly IQQ values",
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: "Time",
				Type: "time",
			}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis", TriggerOn: "click"}),
		)

		// Put data into instance
		line.AddSeries(sensorID, generateLineItemsFromMeasurements(measurements)).
			// AddSeries("Category B", generateLineItems()).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
		line.Render(w)
	}
}

// StartServer starts the http server.
func StartServer(db *sql.DB) {
	http.HandleFunc("/", httpserver(db))
	http.ListenAndServe(":8081", nil)
}
