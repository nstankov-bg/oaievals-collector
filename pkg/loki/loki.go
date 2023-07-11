package loki

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/cortexproject/cortex/pkg/util/flagext"
	kitlog "github.com/go-kit/kit/log"
	"github.com/grafana/loki/pkg/promtail/client"
	"github.com/nstankov-bg/oaievals-collector/pkg/events"
	"github.com/prometheus/common/model"
)

type Client interface {
	Handle(labels model.LabelSet, time time.Time, line string) error
}

type RealClient struct {
	client client.Client
}

func (r *RealClient) Handle(labels model.LabelSet, t time.Time, line string) error {
	return r.client.Handle(labels, t, line)
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

	realClient, err := client.New(config, logger)
	if err != nil {
		log.Fatalf("Could not create Loki client: %s", err)
	}

	lokiClient = &RealClient{client: realClient}
}

func WriteToLoki(event events.Event) {
	labels := model.LabelSet{
		"job":       model.LabelValue(fmt.Sprint("MYJOB")),
		"run_id":    model.LabelValue(fmt.Sprint(event.RunID)),
		"event_id":  model.LabelValue(fmt.Sprint(event.EventID)),
		"sample_id": model.LabelValue(fmt.Sprint(event.SampleID)),
	}

	logLine := ""

	// Add additional fields based on the data in the event
	for key, value := range event.Data {
		logLine += fmt.Sprintf("%s: %v ", key, value)
	}

	// Log the message
	err := lokiClient.Handle(labels, time.Now(), logLine)
	if err != nil {
		log.Fatalf("Could not send log to Loki: %s", err)
	}
}
