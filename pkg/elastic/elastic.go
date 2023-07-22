package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
)

// Elasticsearch configuration
var es *elasticsearch.Client
var bufferMutex sync.Mutex
var messageBuffer []esapi.IndexRequest

// Configurable parameters with default values
var batchSize = 10

func init() {
	// Get Elasticsearch server address from environment variable
	esAddress := os.Getenv("ES_ADDRESS")
	if esAddress == "" {
		log.Printf("No Elasticsearch server address provided, skipping Elasticsearch client initialization")
		return
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			esAddress,
		},
	}

	var err error
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Initialize the message buffer
	messageBuffer = make([]esapi.IndexRequest, 0, batchSize)
}

func WriteToElasticsearch(event events.Event) error {

	if es == nil {
		log.Printf("Elasticsearch client is not initialized, skipping write")
		return nil
	}

	// Include RunID in the message
	msgData := struct {
		events.Event
		RunID string `json:"run_id"`
	}{
		Event: event,
		RunID: event.RunID,
	}

	jsonData, err := json.Marshal(msgData)
	if err != nil {
		log.Printf("Failed to marshal event: %s\n", err)
		return err
	}

	indexName := fmt.Sprintf("evals_%s", strings.ToLower(event.RunID)) // Create a new index for each run, ensure it's all lowercase

	// Check if the index exists
	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		log.Fatalf("Error checking if index exists: %s", err)
	}

	// If the index doesn't exist, create it
	if res.StatusCode == 404 {
		res, err := es.Indices.Create(indexName)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
		if res.IsError() {
			log.Printf("Error creating index: %s", res.String())
		}
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: uuid.New().String(), // Use a UUID as the document ID
		Body:       bytes.NewReader(jsonData),
	}

	bufferMutex.Lock()
	messageBuffer = append(messageBuffer, req)
	bufferMutex.Unlock()

	// If the buffer is full, flush it
	if len(messageBuffer) >= batchSize {
		flushMessageBuffer()
	}

	return nil
}

func flushMessageBuffer() {
	bufferMutex.Lock()
	if len(messageBuffer) == 0 {
		bufferMutex.Unlock()
		return
	}

	// Copy the message buffer and clear the original buffer
	tmpBuffer := make([]esapi.IndexRequest, len(messageBuffer))
	copy(tmpBuffer, messageBuffer)
	messageBuffer = messageBuffer[:0]
	bufferMutex.Unlock()

	for _, req := range tmpBuffer {
		operation := func() error {
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Printf("Error getting response: %s", err)
				return err
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("Error indexing document ID=%s", req.DocumentID)
				return fmt.Errorf("error indexing document ID=%s", req.DocumentID)
			}
			return nil
		}

		// Use exponential backoff strategy
		backoffConfig := backoff.NewExponentialBackOff()
		backoffConfig.MaxElapsedTime = 2 * time.Minute

		err := backoff.Retry(operation, backoffConfig)
		if err != nil {
			log.Printf("Failed to write messages to Elasticsearch after several retries: %s\n", err)
			// Consider whether you want to add the messages back to the buffer in case of failure
		}
	}
}

func Shutdown() {
	if es == nil {
		log.Printf("Elasticsearch client is not initialized, skipping shutdown")
		return
	}

	// Flush any remaining messages in the buffer before shutting down
	flushMessageBuffer()
}
