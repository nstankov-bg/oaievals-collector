package influxdb

import (
	"os"
	"testing"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) WriteAPI(org, bucket string) WriteAPI {
	args := m.Called(org, bucket)
	return args.Get(0).(WriteAPI)
}

type MockWriteAPI struct {
	mock.Mock
}

func (m *MockWriteAPI) WritePoint(p *write.Point) {
	m.Called(p)
}

func (m *MockWriteAPI) Flush() {
	m.Called()
}

func TestWriteToInfluxDB(t *testing.T) {
	os.Setenv("INFLUXDB_HOST", "http://localhost:8086")
	os.Setenv("INFLUXDB_ORG", "your-org")
	os.Setenv("INFLUXDB_BUCKET", "your-bucket")
	os.Setenv("INFLUXDB_TOKEN", "your-token")

	mockClient := new(MockClient)
	client = mockClient

	event := events.Event{
		RunID:    "run1",
		Type:     "match",
		EventID:  1,
		Data:     map[string]interface{}{"correct": true},
		SampleID: "sample1",
	}

	mockWriteAPI := new(MockWriteAPI)
	mockClient.On("WriteAPI", "your-org", "your-bucket").Return(mockWriteAPI)

	mockWriteAPI.On("WritePoint", mock.AnythingOfType("*write.Point")).Return()
	mockWriteAPI.On("Flush").Return()

	WriteToInfluxDB(event)

	mockClient.AssertCalled(t, "WriteAPI", "your-org", "your-bucket")
	mockWriteAPI.AssertCalled(t, "WritePoint", mock.AnythingOfType("*write.Point"))
	mockWriteAPI.AssertCalled(t, "Flush")

	os.Unsetenv("INFLUXDB_HOST")
	os.Unsetenv("INFLUXDB_ORG")
	os.Unsetenv("INFLUXDB_BUCKET")
	os.Unsetenv("INFLUXDB_TOKEN")
}
