
package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Monitoring struct {
	EventCounter *prometheus.CounterVec
}

func NewMonitoring() *Monitoring {
	counterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_total",
			Help: "Total number of events, subdivided by run_id and correctness",
		},
		[]string{"run_id", "correct"},
	)
	prometheus.MustRegister(counterVec)

	return &Monitoring{EventCounter: counterVec}
}
