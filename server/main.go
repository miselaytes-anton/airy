package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"

	messageprocessor "server/messageProcessor"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
)

const (
	defaultBrokerAddress    = "tcp://mosquitto:1883"
	defaultPostgressAddress = "postgresql://tatadata:tatadata@localhost:5432/tatadata?sslmode=disable"
)

func getBrokerAdress() string {
	value, ok := os.LookupEnv("BROKER_ADDRESS")
	if ok {
		return value
	}
	return defaultBrokerAddress
}

// host, port, user, password, dbname := "localhost", 5432, "tatadata", "tatadata", "tatadata"
// psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

	// close database
	defer db.Close()

	// Enable mqtt logging.
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	handlers := messageprocessor.Handlers{Db: db}
	mqttClient := messageprocessor.MakeMqttClient(getBrokerAdress(), handlers)
	messageprocessor.StartProcessing(mqttClient)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	mqttClient.Disconnect(1000)
	fmt.Println("shutdown complete")
}
