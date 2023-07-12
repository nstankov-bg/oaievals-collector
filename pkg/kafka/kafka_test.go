package kafka

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

type MockWriter struct {
	messages []kafka.Message
}

func (mw *MockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	mw.messages = append(mw.messages, msgs...)
	return nil
}

func (mw *MockWriter) Close() error {
	return nil
}

func TestWriteToKafka(t *testing.T) {
	mockWriter := &MockWriter{}
	writer = mockWriter

	event := events.Event{
		RunID:    "run1",
		Type:     "match",
		EventID:  1,
		Data:     map[string]interface{}{"correct": true},
		SampleID: "sample1",
	}

	err := WriteToKafka(event)
	require.NoError(t, err)

	require.Len(t, mockWriter.messages, 1)

	msg := mockWriter.messages[0]
	require.Equal(t, event.Type, string(msg.Key))

	var msgEvent events.Event
	err = json.Unmarshal(msg.Value, &msgEvent)
	require.NoError(t, err)

	require.Equal(t, event.Type, msgEvent.Type)
	require.Equal(t, event.RunID, msgEvent.RunID)
}
