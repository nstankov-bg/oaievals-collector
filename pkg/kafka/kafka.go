package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
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

func init() {
	// Initialize the writer as a *kafka.Writer, which satisfies the KafkaWriter interface
	writer = &kafka.Writer{
		Addr:         kafka.TCP(bootstrapServers...),
		Async:        true, // Enable async to allow batching
		BatchSize:    1000, // Increase batch size
		BatchTimeout: 10 * time.Second,
		Transport:    &kafka.Transport{TLS: nil}, // Update this to your needs
		Logger:       log.New(os.Stdout, "kafka writer: ", 0),
		ErrorLogger:  log.New(os.Stderr, "kafka writer: ", 0),
	}

	// Initialize the message buffer
	messageBuffer = make([]kafka.Message, 0, 1000)
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
	// Only attempt to write messages if there are messages in the buffer
	if len(messageBuffer) > 0 {
		err := writer.WriteMessages(context.Background(), messageBuffer...)
		if err != nil {
			log.Printf("Failed to write messages to kafka: %s\n", err)
		}
		// Clear the message buffer
		messageBuffer = messageBuffer[:0]
	}
}

// WriteToKafka writes the event to Kafka
func WriteToKafka(event events.Event) error {
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
	messageBuffer = append(messageBuffer, msg)

	// If the buffer is full, flush it
	if len(messageBuffer) >= 1000 {
		flushMessageBuffer()
	}

	return nil
}

func Shutdown() {
	// Flush any remaining messages in the buffer before shutting down
	flushMessageBuffer()

	if err := writer.Close(); err != nil {
		log.Printf("Failed to close writer: %s", err)
	}
}
