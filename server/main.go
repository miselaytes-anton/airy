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

	api "server/api"
	messageprocessor "server/messageProcessor"

	_ "github.com/lib/pq"
)

const (
	defaultBrokerAddress    = "tcp://mosquitto:1883"
	defaultPostgressAddress = "postgres://tatadata:tatadata@postgres:5432/tatadata?sslmode=disable"
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

	p := messageprocessor.NewMessageProcessor(getBrokerAdress(), db)
	p.EnableMqttLogging()
	p.StartProcessing()
	api.StartServer(db)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	db.Close()
	p.StopProcessing()
	fmt.Println("shutdown complete")
}
