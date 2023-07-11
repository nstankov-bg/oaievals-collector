package timescaledb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
)

var (
	pool *pgxpool.Pool
	mtx  sync.Mutex
)

func init() {
	Initialize()
}

func Initialize() {
	mtx.Lock()
	defer mtx.Unlock()

	if pool != nil {
		return
	}

	connectToTimescaleDB()
	createEventsTable()
}

func connectToTimescaleDB() {
	var err error
	pool, err = pgxpool.Connect(context.Background(), os.Getenv("TIMESCALEDB_HOST"))
	if err != nil {
		log.Fatalf("Unable to connect to TimescaleDB: %v\n", err)
	} else {
		log.Println("Connected to TimescaleDB successfully.")
	}
}

func createEventsTable() {
	_, err := pool.Exec(context.Background(), `
    CREATE TABLE IF NOT EXISTS events (
        id SERIAL,
        time TIMESTAMPTZ NOT NULL,
        run_id TEXT NOT NULL,
        event_id INTEGER NOT NULL,
        sample_id TEXT NOT NULL,
        data JSONB
    );
    SELECT create_hypertable('events', 'time', if_not_exists => TRUE);
`)
	if err != nil {
		log.Fatalf("Unable to create events table in TimescaleDB: %v\n", err)
	} else {
		log.Println("Events table created successfully in TimescaleDB.")
	}
}

func WriteToTimescaleDB(event events.Event) error {
	sql, args := prepareInsertStatement(event)

	log.Printf("Executing SQL: %s with args: %v\n", sql, args)

	if err := executeStatement(sql, args); err != nil {
		log.Printf("Error executing statement: %v\n", err)
		return fmt.Errorf("error inserting into TimescaleDB: %w", err)
	}

	return nil
}

func executeStatement(sql string, args []interface{}) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, sql, args...)
	if err != nil {
		log.Printf("Error in executeStatement: %v\n", err)
	}
	return err
}

func prepareInsertStatement(event events.Event) (string, []interface{}) {
	dataJson, err := json.Marshal(event.Data)
	if err != nil {
		log.Fatalf("Error preparing insert statement: %v\n", err)
	}

	sql := "INSERT INTO events (time, run_id, event_id, sample_id, data) VALUES ($1, $2, $3, $4, $5)"
	return sql, []interface{}{time.Now(), event.RunID, event.EventID, event.SampleID, dataJson}
}
