package loki

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	kitlog "github.com/go-kit/log"
	flagext "github.com/grafana/dskit/flagext"
	"github.com/grafana/loki/clients/pkg/promtail/api"
	"github.com/grafana/loki/clients/pkg/promtail/client"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

type Client interface {
	Chan() chan<- api.Entry
	Stop()
	Name() string
}

type RealClient struct {
	client client.Client
}

func (r *RealClient) Chan() chan<- api.Entry {
	return r.client.Chan()
}

func (r *RealClient) Stop() {
	r.client.Stop()
}

func (r *RealClient) Name() string {
	return r.client.Name()
}

var lokiClient Client

func init() {
	lokiURL := os.Getenv("LOKI_URL")
	if lokiURL == "" {
		log.Fatalf("No Loki URL provided")
	}

	_, err := url.Parse(lokiURL)
	if err != nil {
		log.Fatalf("Could not parse Loki URL: %s", err)
	}

	urlValue := flagext.URLValue{}
	_ = urlValue.Set(lokiURL)

	config := client.Config{
		URL:       urlValue,
		BatchWait: 1 * time.Second,
		BatchSize: 1024,
		Timeout:   10 * time.Second,
	}

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))

	// Please replace the placeholders with suitable arguments
	reg := prometheus.NewPedanticRegistry()
	metrics := client.NewMetrics(reg)
	realClient, err := client.New(metrics, config, 0, 0, false, logger)
	if err != nil {
		log.Fatalf("Could not create Loki client: %s", err)
	}

	if realClient == nil {
		log.Fatalf("Client is nil")
	}

	lokiClient = &RealClient{client: realClient}
}

func WriteToLoki(event events.Event) {
	labels := model.LabelSet{
		"job":       model.LabelValue("MYJOB"),
		"run_id":    model.LabelValue(fmt.Sprint(event.RunID)),
		"event_id":  model.LabelValue(fmt.Sprint(event.EventID)),
		"sample_id": model.LabelValue(fmt.Sprint(event.SampleID)),
	}

	logLine := ""

	// Add additional fields based on the data in the event
	for key, value := range event.Data {
		logLine += key + ": " + fmt.Sprint(value) + " "
	}

	// Prepare the log entry
	entry := api.Entry{
		Labels: labels,
		Entry: logproto.Entry{
			Timestamp: time.Now(),
			Line:      logLine,
		},
	}

	// Log the message
	lokiClient.Chan() <- entry
}
