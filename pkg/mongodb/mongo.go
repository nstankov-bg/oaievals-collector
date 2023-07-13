package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	mtx    sync.Mutex
)

func init() {
	if os.Getenv("MONGODB_URI") == "" {
		log.Printf("No MongoDB URI provided, skipping MongoDB client initialization")
		return
	}
	Initialize()
}

func Initialize() {
	mtx.Lock()
	defer mtx.Unlock()

	if client != nil {
		return
	}

	connectToMongoDB()
}

func connectToMongoDB() {
	var err error
	uri := os.Getenv("MONGODB_URI")

	username := os.Getenv("MONGODB_USER")
	password := os.Getenv("MONGODB_PASSWORD")

	if username != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s", username, password, uri)
	} else {
		uri = fmt.Sprintf("mongodb://%s", uri)
	}

	client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Unable to create MongoDB client: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to MongoDB: %v\n", err)
	}

	log.Println("Connected to MongoDB successfully.")
}

func WriteToMongoDB(event events.Event) error {
	if client == nil {
		log.Printf("MongoDB client is not initialized, skipping write")
		return nil
	}

	database := os.Getenv("MONGODB_DATABASE")
	if database == "" {
		database = "test" // default database name
	}

	collectionName := os.Getenv("MONGODB_COLLECTION")
	if collectionName == "" {
		collectionName = "events" // default collection name
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(database).Collection(collectionName)
	res, err := collection.InsertOne(ctx, bson.M{
		"time":      time.Now(),
		"run_id":    event.RunID,
		"event_id":  event.EventID,
		"sample_id": event.SampleID,
		"data":      event.Data,
	})
	if err != nil {
		log.Printf("Unable to insert event into MongoDB: %v\n", err)
		return err
	}

	log.Printf("Successfully inserted event with ID %v into MongoDB\n", res.InsertedID)
	return nil
}
