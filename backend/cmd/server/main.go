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

	"github.com/miselaytes-anton/tatadata/backend/internal/config"
	"github.com/miselaytes-anton/tatadata/backend/internal/models"
)

func main() {
	db, err := sql.Open("postgres", config.GetPostgresAddress())
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
	server := &Server{
		Router:       router,
		Measurements: measurements,
		Events:       events,
	}
	server.Routes()
	log.Print("server is listening, view http://localhost:8081/api/graphs")
	http.ListenAndServe(":8081", router)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")
	db.Close()
	fmt.Println("shutdown complete")
}
