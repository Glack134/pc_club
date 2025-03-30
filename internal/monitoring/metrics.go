package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	SessionsActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_sessions_count",
		Help: "Current number of active gaming sessions",
	})

	SessionsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sessions_created_total",
		Help: "Total number of created gaming sessions",
	})
)

func RecordSessionStart() {
	SessionsActive.Inc()
	SessionsCreated.Inc()
}

func RecordSessionEnd() {
	SessionsActive.Dec()
}
