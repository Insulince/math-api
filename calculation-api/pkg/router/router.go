package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"math-api/calculation-api/pkg/api"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math-api/base-api"
)

type Router struct {
	*mux.Router
}

func CreateRouter(api *api.Api) (router *Router) {
	var homeCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "math_api_home",
			Help: "Math API Home",
		},
	)
	var homeLatency = promauto.NewSummary(
		prometheus.SummaryOpts{
			Name: "math_api_home_latency",
			Help: "Math API Home Latency",
		},
	)
	var healthCheckCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "math_api_health_check",
			Help: "Math API Health Check",
		},
	)
	var healthCheckLatency = promauto.NewSummary(
		prometheus.SummaryOpts{
			Name: "math_api_health_check_latency",
			Help: "Math API Health Check Latency",
		},
	)
	var notFoundCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "math_api_not_found",
			Help: "Math API Not Found",
		},
	)
	var notFoundLatency = promauto.NewSummary(
		prometheus.SummaryOpts{
			Name: "math_api_not_found_latency",
			Help: "Math API Not Found Latency",
		},
	)
	var calculateCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "math_api_calculate",
			Help: "Math API Calculate",
		},
	)
	var calculateLatency = promauto.NewSummary(
		prometheus.SummaryOpts{
			Name: "math_api_calculate_latency",
			Help: "Math API Calculate Latency",
		},
	)

	router = &Router{
		mux.NewRouter().StrictSlash(true),
	}

	router.HandleFunc("/home", base_api.WrapHandler(api.HomeHandler, base_api.NewPrometheusMetrics(homeLatency, homeCounter))).Methods(http.MethodGet)

	router.HandleFunc("/health", base_api.WrapHandler(api.HealthCheckHandler, base_api.NewPrometheusMetrics(healthCheckLatency, healthCheckCounter))).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/", base_api.WrapHandler(api.CalculateHandler, base_api.NewPrometheusMetrics(calculateLatency, calculateCounter))).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(base_api.WrapHandler(api.NotFoundHandler, base_api.NewPrometheusMetrics(notFoundLatency, notFoundCounter)))

	return router
}
