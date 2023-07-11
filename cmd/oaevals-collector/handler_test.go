// handler_test.go
package oaievals_collector

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nstankov-bg/oaievals-collector/pkg/timescaledb"
)

func TestEventHandler(t *testing.T) {
	timescaledb.Initialize() // initialize TimescaleDB before running the test

	req, err := http.NewRequest("POST", "/events", strings.NewReader(`{
        "run_id": "123",
        "type": "match",
        "event_id": 1,
        "data": {
            "correct": true
        },
        "sample_id": "456"
    }`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(eventHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	expected := `Event received`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
