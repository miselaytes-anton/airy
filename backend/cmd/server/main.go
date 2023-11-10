package main

import (
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
	"github.com/miselaytes-anton/tatadata/backend/internal/log"
	"github.com/miselaytes-anton/tatadata/backend/internal/models"
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
	events := models.EventModel{DB: db}

	router := http.NewServeMux()
	server := &Server{
		Router:       router,
		Measurements: measurements,
		Events:       events,
		LogError:     log.Error,
		LogInfo:      log.Info,
	}
	server.Routes()

	log.Info.Println("server is listening on :8081")
	log.Info.Println("visit http://localhost:8081/api/graphs")

	srv := &http.Server{
		Addr: ":8081", ErrorLog: log.Error, Handler: router,
	}
	srv.ListenAndServe()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	log.Info.Println("signal caught - exiting")
	db.Close()
	log.Info.Println("shutdown complete")
}
