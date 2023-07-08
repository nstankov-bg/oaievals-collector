package oaievals_collector

import (
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	go Run() // Run the server in a goroutine

	// Wait for the server to start up
	time.Sleep(time.Second)

	// Test /metrics endpoint
	resp, err := http.Get("http://localhost:8080/metrics")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Test /events endpoint with GET (should be not allowed)
	resp, err = http.Get("http://localhost:8080/events")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed, got %v", resp.StatusCode)
	}
}
