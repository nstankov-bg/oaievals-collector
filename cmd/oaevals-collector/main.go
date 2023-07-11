package oaievals_collector

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run() {
	log.Println("Starting up the event collector...")

	go ProcessFilesInDirectory() // Start processing files in a new goroutine.

	mux := http.NewServeMux()

	mux.HandleFunc("/events", eventHandler)

	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
