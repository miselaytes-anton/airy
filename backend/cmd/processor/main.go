package main

import (
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	// imports postgres timezones data
	_ "time/tzdata"
	// postgres driver
	_ "github.com/lib/pq"

	"github.com/miselaytes-anton/tatadata/backend/internal/config"
	"github.com/miselaytes-anton/tatadata/backend/internal/log"
	"github.com/miselaytes-anton/tatadata/backend/internal/models"
)

const (
	mqttClientID     = "tatadata"
	measurementTopic = "measurement"
	measurementQOS   = 1
)

func main() {
	db, err := sql.Open("postgres", config.GetPostgresAddress())
	if err != nil {
		log.Error.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Error.Fatal(err)
	}

	measurements := models.MeasurementModel{DB: db}

	handler := MeasurementHandler{
		Measurements: measurements,
	}

	options := MqttClientOpts{
		BrokerAddress: config.GetBrokerAdress(),
		ClientID:      mqttClientID,
		MessageHandlers: MessageHandlers{
			measurementTopic: {
				Handler: handler.OnMessageHandler,
				QOS:     measurementQOS,
			},
		},
	}

	mqttClient := NewMqttClient(options)

	p := Processor{
		Client:   mqttClient,
		LogError: log.Error,
		LogInfo:  log.Info,
	}
	p.EnableMqttLogging()
	p.StartProcessing()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	log.Info.Println("signal caught - exiting")
	db.Close()
	p.StopProcessing()
	log.Info.Println("shutdown complete")
}
