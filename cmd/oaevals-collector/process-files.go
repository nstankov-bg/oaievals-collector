package oaievals_collector

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/nstankov-bg/oaievals-collector/pkg/influxdb"
)

var (
	dataDir      string
	processedDir string
)

func init() {
	baseDir := os.Getenv("WATCHED_FOLDER")

	if baseDir == "" {
		baseDir, _ = filepath.Abs(".")
	}

	dataDir = filepath.Join(baseDir, "data")
	processedDir = filepath.Join(baseDir, "data/processed")

	if err := createDirIfNotExist(dataDir); err != nil {
		log.Fatalf("Failed to create directory %s: %v", dataDir, err)
	}

	if err := createDirIfNotExist(processedDir); err != nil {
		log.Fatalf("Failed to create directory %s: %v", processedDir, err)
	}
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o755)
	}
	return nil
}

func ProcessFilesInDirectory() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Printf("Error getting file info for %s: %v", event.Name, err)
						continue
					}
					processFile(fileInfo)
				}
			case err, ok := <-watcher.Errors:
				if ok {
					log.Printf("Error: %v", err)
				}
			}
		}
	}()

	if err = watcher.Add(dataDir); err != nil {
		log.Fatalf("Failed to add directory to watcher: %v", err)
	}
	<-done
}

func processFile(fileInfo os.FileInfo) {
	filePath := filepath.Join(dataDir, fileInfo.Name())
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event events.Event
		if err = json.Unmarshal(scanner.Bytes(), &event); err != nil {
			log.Printf("Error unmarshalling event from file %s: %v", filePath, err)
			continue
		}
		processEvent(event)
	}

	if err = scanner.Err(); err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
	}

	newPath := filepath.Join(processedDir, fileInfo.Name())
	if err = os.Rename(filePath, newPath); err != nil {
		log.Printf("Error moving file from %s to %s: %v", filePath, newPath, err)
	}
}

func processEvent(event events.Event) {
	if event.Type == "match" {
		if correctVal, exists := event.Data["correct"]; exists {
			if correct, ok := correctVal.(bool); ok {
				if correct {
					mon.EventCounter.WithLabelValues(event.RunID, "true").Inc()
				} else {
					mon.EventCounter.WithLabelValues(event.RunID, "false").Inc()
				}
			} else {
				log.Printf("The 'correct' value in event %v is not a boolean", event)
			}
		} else {
			log.Printf("The 'correct' key does not exist in the data of event %v", event)
		}

		influxdb.WriteToInfluxDB(event)
	}
}
