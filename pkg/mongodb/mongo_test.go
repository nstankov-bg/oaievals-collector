package mongodb

import (
	"os"
	"testing"
	"time"

	"github.com/nstankov-bg/oaievals-collector/pkg/events"
)

func TestWriteToMongoDB(t *testing.T) {
	// Set the environment variables
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("MONGODB_USER", "your mongodb username")
	os.Setenv("MONGODB_PASSWORD", "your mongodb password")
	os.Setenv("MONGODB_DATABASE", "mydatabase")
	os.Setenv("MONGODB_COLLECTION", "mycollection")

	Initialize()

	// Create a test event
	event := events.Event{
		RunID:    "testRun",
		EventID:  1,
		SampleID: "testSample",
		Data: map[string]interface{}{
			"field1": "value1",
			"field2": 2,
			"field3": time.Now(),
		},
	}

	err := WriteToMongoDB(event)
	if err != nil {
		t.Fatalf("Error writing to MongoDB: %v", err)
	}
}
