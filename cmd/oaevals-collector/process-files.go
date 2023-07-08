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

func ProcessFilesInDirectory() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("Event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file: ", event.Name)
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Println("Error getting file info: ", err)
						continue
					}
					processFile(fileInfo)
				}
			case err := <-watcher.Errors:
				log.Println("Error: ", err)
			}
		}
	}()

	err = watcher.Add("/data")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func processFile(fileInfo os.FileInfo) {
	filePath := filepath.Join("/data", fileInfo.Name())
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event events.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			log.Println("Error unmarshalling event from file: ", err)
			continue
		}
		processEvent(event)
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading file: ", err)
	}

	newPath := filepath.Join("/data/processed", fileInfo.Name())
	if err := os.Rename(filePath, newPath); err != nil {
		log.Println("Error moving file: ", err)
	}
}

func processEvent(event events.Event) {
	log.Println("Processing event of type: ", event.Type)
	if event.Type == "match" {
		correctVal, exists := event.Data["correct"]
		if exists {
			correct, ok := correctVal.(bool)
			if ok {
				if correct {
					mon.EventCounter.WithLabelValues(event.RunID, "true").Inc()
					log.Println("Incremented event counter for: ", event.RunID, "with correctness: true")
				} else {
					mon.EventCounter.WithLabelValues(event.RunID, "false").Inc()
					log.Println("Incremented event counter for: ", event.RunID, "with correctness: false")
				}
			} else {
				log.Println("The 'correct' value is not a boolean")
			}
		} else {
			log.Println("The 'correct' key does not exist in the data")
		}

		influxdb.WriteToInfluxDB(event)
	}
}
