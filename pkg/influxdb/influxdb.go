package influxdb

import (
	"os"
	"time"
	"log"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
)

type Client interface {
	WriteAPI(org, bucket string) WriteAPI
}

type WriteAPI interface {
	WritePoint(p *write.Point)
	Flush()
}

type RealClient struct {
	client influxdb2.Client
}

func (r *RealClient) WriteAPI(org, bucket string) WriteAPI {
	return &RealWriteAPI{api: r.client.WriteAPI(org, bucket)}
}

type RealWriteAPI struct {
	api api.WriteAPI
}

func (r *RealWriteAPI) WritePoint(p *write.Point) {
	r.api.WritePoint(p)
}

func (r *RealWriteAPI) Flush() {
	r.api.Flush()
}

var client Client

func init() {
	token := os.Getenv("INFLUXDB_TOKEN")
	if token == "" {
		log.Fatalf("No InfluxDB token provided")
	}
	client = &RealClient{client: influxdb2.NewClient(os.Getenv("INFLUXDB_HOST"), token)}
}

func WriteToInfluxDB(event events.Event) {
    // get non-blocking write client
    writeAPI := client.WriteAPI(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))

    // create point
    p := influxdb2.NewPointWithMeasurement("events").
        AddTag("run_id", event.RunID).
        AddField("event_id", event.EventID).
        AddField("sample_id", event.SampleID).
        SetTime(time.Now())
    
    // Add additional fields based on the data in the event
    for key, value := range event.Data {
        switch v := value.(type) {
        case bool:
            p.AddField(key, v)
        case string:
            p.AddField(key, v)
        case float64:
            p.AddField(key, v)
        case int:
            p.AddField(key, v)
        default:
            log.Printf("Unsupported type for field %s", key)
        }
    }

    // write point immediately
    writeAPI.WritePoint(p)
    // Ensures background processes finishes
    writeAPI.Flush()
}

