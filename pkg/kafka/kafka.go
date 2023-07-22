package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cenkalti/backoff/v4"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/segmentio/kafka-go"
)

// KafkaWriter is an interface that allows for writing messages to Kafka
type KafkaWriter interface {
	WriteMessages(context.Context, ...kafka.Message) error
	Close() error
}

// Global writer variable of the KafkaWriter interface type
var writer KafkaWriter
var bootstrapServers = strings.Split(os.Getenv("KAFKA_BOOTSTRAP_SERVERS"), ",")
var messageBuffer []kafka.Message // Buffer to hold messages
var bufferMutex sync.Mutex        // Mutex to protect the message buffer
var stopChan chan struct{}        // Channel to signal flushing goroutine to stop
var wg sync.WaitGroup             // WaitGroup to wait for flushing goroutine to stop

// Configurable parameters with default values
var batchSize = 10
var batchTimeout = 10 * time.Second
var flushInterval = 1 * time.Second

func init() {
	bootstrapServersEnv := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if bootstrapServersEnv == "" {
		log.Printf("No Kafka bootstrap servers provided, skipping Kafka client initialization")
		return
	}
	bootstrapServers = strings.Split(bootstrapServersEnv, ",")

	// Initialize the writer as a *kafka.Writer, which satisfies the KafkaWriter interface
	writer = &kafka.Writer{
		Addr:         kafka.TCP(bootstrapServers...),
		Async:        true, // Enable async to allow batching
		BatchSize:    batchSize,
		BatchTimeout: batchTimeout,
		Transport:    &kafka.Transport{TLS: nil}, // Update this to your needs
		Logger:       log.New(os.Stdout, "kafka writer: ", 0),
		ErrorLogger:  log.New(os.Stderr, "kafka writer: ", 0),
	}

	// Initialize the message buffer
	messageBuffer = make([]kafka.Message, 0, batchSize)

	// Start the flushing goroutine
	stopChan = make(chan struct{})
	go flushMessageBufferPeriodically()
}

func ensureTopicExists(topic string) error {
	config := sarama.NewConfig()
	config.Net.TLS.Enable = false    // Update this to your needs
	config.Version = sarama.V2_3_0_0 // Update this to your Kafka broker version

	admin, err := sarama.NewClusterAdmin(bootstrapServers, config)
	if err != nil {
		return err
	}
	defer admin.Close()

	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}

	if _, ok := topics[topic]; ok {
		return nil // Topic already exists, no need to create it
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = admin.CreateTopic(topic, topicDetail, false)
	return err
}

func flushMessageBuffer() {
	bufferMutex.Lock()
	if len(messageBuffer) == 0 {
		bufferMutex.Unlock()
		return
	}

	// Copy the message buffer and clear the original buffer
	tmpBuffer := make([]kafka.Message, len(messageBuffer))
	copy(tmpBuffer, messageBuffer)
	messageBuffer = messageBuffer[:0]
	bufferMutex.Unlock()

	operation := func() error {
		return writer.WriteMessages(context.Background(), tmpBuffer...)
	}

	// Use exponential backoff strategy
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.MaxElapsedTime = 2 * time.Minute // example max elapsed time

	err := backoff.Retry(operation, backoffConfig)
	if err != nil {
		log.Printf("Failed to write messages to kafka after several retries: %s\n", err)
		// Consider whether you want to add the messages back to the buffer in case of failure
	}
}

func flushMessageBufferPeriodically() {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			flushMessageBuffer()
		case <-stopChan:
			return
		}
	}
}

// WriteToKafka writes the event to Kafka
func WriteToKafka(event events.Event) error {

	if writer == nil {
		log.Printf("Kafka client is not initialized, skipping write")
		return nil
	}

	// Ensure the topic "evals" exists
	err := ensureTopicExists("evals")
	if err != nil {
		log.Printf("Failed to ensure topic exists: %s\n", err)
		return err
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

	// Add the message to the buffer instead of writing immediately
	msg := kafka.Message{
		Key:   []byte(event.RunID), // If you have a key to use, place it here
		Value: jsonData,
		Topic: "evals", // Using "evals" as topic
	}

	bufferMutex.Lock()
	messageBuffer = append(messageBuffer, msg)
	bufferMutex.Unlock()

	// If the buffer is full, flush it
	if len(messageBuffer) >= batchSize {
		flushMessageBuffer()
	}

	return nil
}

func Shutdown() {
	if writer == nil {
		log.Printf("Kafka client is not initialized, skipping shutdown")
		return
	}
	// Stop the flushing goroutine and wait for it to finish
	close(stopChan)
	wg.Wait()

	// Flush any remaining messages in the buffer before shutting down
	flushMessageBuffer()

	if err := writer.Close(); err != nil {
		log.Printf("Failed to close writer: %s", err)
	}
}
