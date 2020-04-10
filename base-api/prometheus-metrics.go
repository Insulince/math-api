package base_api

import "github.com/prometheus/client_golang/prometheus"

type PrometheusMetrics struct {
	Summary prometheus.Summary `json:"summary"`
	Counter prometheus.Counter `json:"counter"`
}

func NewPrometheusMetrics(summary prometheus.Summary, counter prometheus.Counter) (prometheusMetrics *PrometheusMetrics) {
	return &PrometheusMetrics{
		Summary: summary,
		Counter: counter,
	}
}