package loki

import (
	"testing"
	"time"

	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/mock"
)

// Define a mock struct to be used in your unit tests of myFunc.
type MockClient struct {
	mock.Mock
}

// Mock the Handle method
func (m *MockClient) Handle(labels model.LabelSet, t time.Time, line string) error {
	args := m.Called(labels, t, line)
	return args.Error(0)
}

func TestWriteToLoki(t *testing.T) {
	// create an instance of our test object
	mockClient := new(MockClient)

	// setup expectations
	mockEvent := events.Event{
		RunID:    "testRun",
		EventID:  456,
		SampleID: "testSample",
		Data: map[string]interface{}{
			"testKey": "testValue",
		},
	}

	expectedLabels := model.LabelSet{
		"job":       model.LabelValue("MYJOB"),
		"run_id":    model.LabelValue("testRun"),
		"event_id":  model.LabelValue("456"), // Adjusted to match the mockEvent EventID
		"sample_id": model.LabelValue("testSample"),
	}

	expectedLogLine := "testKey: testValue "

	mockClient.On("Handle", expectedLabels, mock.AnythingOfType("time.Time"), expectedLogLine).Return(nil)

	// set the global client variable to our mock client
	lokiClient = mockClient

	// call the code we are testing
	WriteToLoki(mockEvent)

	// assert that the expectations were met
	mockClient.AssertExpectations(t)
}
