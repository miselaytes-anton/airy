package main

import (
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	// imports postgres timezones data
	_ "time/tzdata"
	// postgres driver
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"

	"github.com/miselaytes-anton/airy/internal/config"
	"github.com/miselaytes-anton/airy/internal/log"
	"github.com/miselaytes-anton/airy/internal/models"
)

func enableMqttLogging() {
	mqtt.ERROR = log.Error
	mqtt.CRITICAL = log.Critical
	mqtt.WARN = log.Warning
	// mqtt.DEBUG = log.Debug
}

func connectToDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.GetPostgresAddress())
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, err
}

func main() {
	const (
		mqttClientID                = "tatadata"
		measurementTopic            = "measurement"
		measurementQOS              = 1
		waithBeforeMqttDisconnectMs = 1000
	)
	enableMqttLogging()

	db, err := connectToDb()

	if err != nil {
		log.Error.Fatal(err)
	}

	measurements := models.MeasurementModel{DB: db}

	handler := measurementHandler{
		Measurements: measurements,
		LogError:     log.Error,
		LogInfo:      log.Info,
	}

	options := mqttClientOpts{
		BrokerAddress: config.GetBrokerAdress(),
		ClientID:      mqttClientID,
		MessageHandlers: messageHandlers{
			measurementTopic: {
				Handler: handler.handle,
				QOS:     measurementQOS,
			},
		},
		LogError: log.Error,
		LogInfo:  log.Info,
	}

	mqttClient := NewMqttClient(options)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Error.Fatal(token.Error())
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	log.Info.Println("signal caught - exiting")
	db.Close()
	mqttClient.Disconnect(waithBeforeMqttDisconnectMs)
	log.Info.Println("shutdown complete")
}
