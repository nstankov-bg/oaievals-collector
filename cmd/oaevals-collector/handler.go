package oaievals_collector

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/nstankov-bg/oaievals-collector/pkg/events"
    "github.com/nstankov-bg/oaievals-collector/pkg/influxdb"
    "github.com/nstankov-bg/oaievals-collector/pkg/monitoring"
)

var mon = monitoring.NewMonitoring()

func eventHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received a request on /events")
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var event events.Event
    err := json.NewDecoder(r.Body).Decode(&event)
    if err != nil {
        http.Error(w, "Error unmarshalling JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    log.Println("Processing event of type: ", event.Type)
    if event.Type == "match" {
        // Accessing the value "correct" from the map.
        correctVal, exists := event.Data["correct"]
        if exists {
            correct, ok := correctVal.(bool)
            if ok {
                if correct {
                    mon.EventCounter.WithLabelValues(event.RunID, "true").Inc()
                    log.Println("Incremented event counter for: ", event.RunID, "with correctness: true")
                } else {
                    mon.EventCounter.WithLabelValues(event.RunID, "false").Inc()
                    log.Println("Incremented event counter for: ", event.RunID, "with correctness: false")
                }
            } else {
                log.Println("The 'correct' value is not a boolean")
            }
        } else {
            log.Println("The 'correct' key does not exist in the data")
        }

        // Write to InfluxDB.
        influxdb.WriteToInfluxDB(event)
    }

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte("Event received"))
    log.Println("Event processed successfully")
}
