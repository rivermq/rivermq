package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTPResponseCtr allows the counting of Http Responses and their status codes
	HTTPResponseCtr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_resp_total",
			Help: "Number of http responses",
		},
		[]string{"code"},
	)

	// MessageCtr allows the counting of Messages and their status
	MessageCtr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "message_total",
			Help: "Number of messages received",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(HTTPResponseCtr)
	prometheus.MustRegister(MessageCtr)
}
