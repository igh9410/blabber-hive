package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var Reg = prometheus.NewRegistry()

func InitPrometheusMetrics() {
	// Register default collectors
	Reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	Reg.MustRegister(collectors.NewGoCollector())

	// Register custom metrics here if you have any
	// e.g., reg.MustRegister(yourCustomMetric)
}
