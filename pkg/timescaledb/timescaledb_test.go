package timescaledb

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/stretchr/testify/require"
)

func createMockEvent(eventTime time.Time) events.Event {
	return events.Event{
		CreatedAt: eventTime.Format(time.RFC3339),
		RunID:     "run1",
		EventID:   123,
		SampleID:  "sample1",
		Data:      map[string]interface{}{"key": "value"},
	}
}

func fetchEvent(runID string, eventID int, sampleID string) (time.Time, string, int, string, map[string]interface{}, error) {
	row := pool.QueryRow(context.Background(), "SELECT time, run_id, event_id, sample_id, data FROM events WHERE run_id=$1 AND event_id=$2 AND sample_id=$3", runID, eventID, sampleID)

	var readTime time.Time
	var rID string
	var eID int
	var sID string
	var dataJson []byte
	var data map[string]interface{}

	err := row.Scan(&readTime, &rID, &eID, &sID, &dataJson)
	if err != nil {
		return readTime, rID, eID, sID, data, err
	}

	err = json.Unmarshal(dataJson, &data)
	return readTime, rID, eID, sID, data, err
}

func TestWriteToTimescaleDB(t *testing.T) {
	Initialize()

	event := createMockEvent(time.Now())

	err := WriteToTimescaleDB(event)
	require.NoError(t, err)

	readTime, runID, eventID, sampleID, data, err := fetchEvent(event.RunID, event.EventID, event.SampleID)

	require.NoError(t, err)
	require.Equal(t, event.RunID, runID)
	require.Equal(t, event.EventID, eventID)
	require.Equal(t, event.SampleID, sampleID)
	require.Equal(t, event.Data, data)

	// We only check if the readTime is not zero value.
	// It is not feasible to compare the exact time due to possible delay between insertion and selection.
	require.False(t, readTime.IsZero())
}
