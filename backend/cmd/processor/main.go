package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	// imports postgres timezones data
	_ "time/tzdata"
	// postgres driver
	_ "github.com/lib/pq"

	"github.com/miselaytes-anton/tatadata/backend/models"
	"github.com/miselaytes-anton/tatadata/backend/processor"
)

const (
	mqttClientID     = "tatadata"
	measurementTopic = "measurement"
	measurementQOS   = 1
)

func getBrokerAdress() string {
	value, ok := os.LookupEnv("BROKER_ADDRESS")
	if !ok {
		panic("BROKER_ADDRESS environment variable not set")
	}
	return value
}

func getPostgresAddress() string {
	value, ok := os.LookupEnv("POSTGRES_ADDRESS")
	if !ok {
		panic("POSTGRES_ADDRESS environment variable not set")
	}
	return value
}

func main() {
	db, err := sql.Open("postgres", getPostgresAddress())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	measurements := models.MeasurementModel{DB: db}

	handler := processor.MeasurementHandler{
		Measurements: measurements,
	}

	options := processor.MqttClientOpts{
		BrokerAddress: getBrokerAdress(),
		ClientID:      mqttClientID,
		MessageHandlers: processor.MessageHandlers{
			measurementTopic: {
				Handler: handler.OnMessageHandler,
				QOS:     measurementQOS,
			},
		},
	}

	mqttClient := processor.MakeMqttClient(options)

	p := processor.MessageProcessor{
		Client: mqttClient,
	}
	p.EnableMqttLogging()
	p.StartProcessing()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	db.Close()
	p.StopProcessing()
	fmt.Println("shutdown complete")
}
