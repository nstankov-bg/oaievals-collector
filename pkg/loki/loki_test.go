package loki

import (
	"testing"
	"time"

	"github.com/grafana/loki/clients/pkg/promtail/api"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

// MockClient is a mock implementation of the Client interface for testing
type MockClient struct {
	received chan api.Entry
}

func (m *MockClient) Chan() chan<- api.Entry {
	return m.received
}

func (m *MockClient) Stop() {}

func (m *MockClient) Name() string {
	return "mock"
}

func TestWriteToLoki(t *testing.T) {
	// Set up a mock Loki client
	mockClient := &MockClient{
		received: make(chan api.Entry, 1),
	}
	lokiClient = mockClient

	// Create a test event
	event := events.Event{
		RunID:    "run1",
		EventID:  123,
		SampleID: "sample1",
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Call the function under test
	WriteToLoki(event)

	// Assert that a log entry was received by the mock client
	select {
	case entry := <-mockClient.received:
		expectedLabels := model.LabelSet{
			"job":       "MYJOB",
			"run_id":    "run1",
			"event_id":  "123",
			"sample_id": "sample1",
		}
		assert.Equal(t, expectedLabels, entry.Labels)
		assert.Contains(t, entry.Entry.Line, "key1: value1")
		assert.Contains(t, entry.Entry.Line, "key2: value2")
	case <-time.After(time.Second):
		t.Fatal("No log entry received by mock client")
	}
}
