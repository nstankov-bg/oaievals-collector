package elastic

import (
	"testing"

	"github.com/nstankov-bg/oaievals-collector/pkg/events"
)

func TestWriteToElasticsearch(t *testing.T) {
	// Initialize the Elasticsearch client
	init()

	// Create a test event
	event := events.Event{
		RunID: "test-run-id",
		// Fill in other necessary fields
	}

	// Call the function with the test event
	err := WriteToElasticsearch(event)

	// Check that there was no error
	if err != nil {
		t.Errorf("WriteToElasticsearch returned an error: %v", err)
	}

	// Check that the message buffer contains the message
	// (In a real test, you might not want to check internal state like this,
	// but it's one way to verify that the function is doing what it's supposed to.)
	if len(messageBuffer) != 1 {
		t.Errorf("messageBuffer has length %v, want 1", len(messageBuffer))
	}

	// Call flushMessageBuffer
	flushMessageBuffer()

	// Check that the message buffer is now empty
	if len(messageBuffer) != 0 {
		t.Errorf("messageBuffer has length %v, want 0", len(messageBuffer))
	}
}
