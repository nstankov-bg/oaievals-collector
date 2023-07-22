package kafka

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Initialize Prometheus metrics
var (
	MessagesSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "kafka_messages_sent_total",
		Help: "The total number of messages sent",
	})
	SendErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "kafka_send_errors_total",
		Help: "The total number of send errors",
	})
)
