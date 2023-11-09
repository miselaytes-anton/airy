package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	// imports postgres timezones data
	_ "time/tzdata"
	// postgres driver
	_ "github.com/lib/pq"

	"github.com/miselaytes-anton/tatadata/backend/models"
	"github.com/miselaytes-anton/tatadata/backend/server"
)

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
	events := models.EventModel{DB: db}

	router := http.NewServeMux()
	server := &server.Server{
		Router:       router,
		Measurements: measurements,
		Events:       events,
	}
	server.Routes()
	log.Print("listening on http://localhost:8081")
	http.ListenAndServe(":8081", router)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	db.Close()
	fmt.Println("shutdown complete")
}
