package events

import (
	"testing"
)

func TestEvent(t *testing.T) {
	event := Event{
		RunID:    "run1",
		Type:     "match",
		EventID:  1,
		Data:     Data{"correct": true},
		SampleID: "sample1",
	}

	if event.RunID != "run1" {
		t.Errorf("Expected run1, got %s", event.RunID)
	}

	if event.Type != "match" {
		t.Errorf("Expected match, got %s", event.Type)
	}

	if event.EventID != 1 {
		t.Errorf("Expected 1, got %d", event.EventID)
	}

	correct, exists := event.Data["correct"]
	if !exists || correct != true {
		t.Errorf("Expected true, got %v", correct)
	}

	if event.SampleID != "sample1" {
		t.Errorf("Expected sample1, got %s", event.SampleID)
	}
}
