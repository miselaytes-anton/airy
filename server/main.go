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

	api "server/api"
	messageprocessor "server/messageProcessor"
	models "server/models"

	_ "github.com/lib/pq"
)

const (
	defaultBrokerAddress    = "tcp://mosquitto:1883"
	defaultPostgressAddress = "postgres://tatadata:tatadata@postgres:5432/tatadata?sslmode=disable"
	mqttClientID            = "tatadata"
	measurementTopic        = "measurement"
	measurementQOS          = 1
)

func getBrokerAdress() string {
	value, ok := os.LookupEnv("BROKER_ADDRESS")
	if ok {
		return value
	}
	return defaultBrokerAddress
}

func getPostgresAddress() string {
	value, ok := os.LookupEnv("POSTGRES_ADDRESS")
	if ok {
		return value
	}
	return defaultPostgressAddress
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
