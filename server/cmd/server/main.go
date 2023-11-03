package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	_ "time/tzdata"

	// postgres driver
	_ "github.com/lib/pq"

	"github.com/miselaytes-anton/tatadata/server/api"
	messageprocessor "github.com/miselaytes-anton/tatadata/server/messageProcessor"
	"github.com/miselaytes-anton/tatadata/server/models"
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

	handler := messageprocessor.MeasurementHandler{
		Measurements: measurements,
	}

	options := messageprocessor.MqttClientOpts{
		BrokerAddress: getBrokerAdress(),
		ClientID:      mqttClientID,
		MessageHandlers: messageprocessor.MessageHandlers{
			measurementTopic: {
				Handler: handler.OnMessageHandler,
				QOS:     measurementQOS,
			},
		},
	}

	mqttClient := messageprocessor.MakeMqttClient(options)

	p := messageprocessor.MessageProcessor{
		Client: mqttClient,
	}
	p.EnableMqttLogging()
	p.StartProcessing()

	serverEnv := &api.ServerEnv{
		Measurements: measurements,
		Events:       models.EventModel{DB: db},
	}
	api.StartServer(serverEnv)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	db.Close()
	p.StopProcessing()
	fmt.Println("shutdown complete")
}
