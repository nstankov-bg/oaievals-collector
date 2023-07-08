package monitoring

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestEventMetricsRegistration(t *testing.T) {
	mon := NewMonitoring()
	mon.EventCounter.WithLabelValues("test_run_id", "true").Inc()

	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	assert.NoError(t, err, "Failed to gather metrics")
	assert.NotEmpty(t, metricFamilies, "No metrics gathered")

	found := false
	for _, metricFamily := range metricFamilies {
		if metricFamily.GetName() == "events_total" {
			found = true
			break
		}
	}

	assert.True(t, found, "Metric 'events_total' not found")
}
