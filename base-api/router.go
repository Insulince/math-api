package base_api

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus"
)

type Router struct {
	*mux.Router
}

func CreateRouter(api *Api, operation *Operation) (router *Router) {
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
	var operateCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "math_api_operate",
			Help: "Math API Operate",
		},
	)
	var operateLatency = promauto.NewSummary(
		prometheus.SummaryOpts{
			Name: "math_api_operate_latency",
			Help: "Math API Operate Latency",
		},
	)

	router = &Router{
		mux.NewRouter().StrictSlash(true),
	}

	router.HandleFunc("/home", WrapHandler(api.HomeHandler, NewPrometheusMetrics(homeLatency, homeCounter))).Methods(http.MethodGet)

	router.HandleFunc("/health", WrapHandler(api.HealthCheckHandler, NewPrometheusMetrics(healthCheckLatency, healthCheckCounter))).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/", WrapHandler(api.OperateHandler(operation), NewPrometheusMetrics(operateLatency, operateCounter))).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(WrapHandler(api.NotFoundHandler, NewPrometheusMetrics(notFoundLatency, notFoundCounter)))

	return router
}
