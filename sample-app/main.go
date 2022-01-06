package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sample_app_requests_count",
			Help: "Request counter",
		},
		[]string{"success"},
	)
)

func init() {
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
	prometheus.MustRegister(requestCounter)
}

type response struct {
	Node       string `json:"node"`
	Namespace  string `json:"namespace"`
	Pod        string `json:"pod"`
	Deployment string `json:"deployment"`
}

func handleSuccess(response response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCounter.
			WithLabelValues("true").
			Inc()
		_ = json.NewEncoder(w).Encode(response)
	}
}

func handleError(status int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCounter.
			WithLabelValues("false").
			Inc()
		w.WriteHeader(status)
	}
}

func main() {
	metricsPort, prometheusEnabled := os.LookupEnv("METRICS_PORT")
	httpPort, _ := os.LookupEnv("HTTP_PORT")
	node, _ := os.LookupEnv("KUBERNETES_NODE")
	namespace, _ := os.LookupEnv("KUBERNETES_NAMESPACE")
	pod, _ := os.LookupEnv("KUBERNETES_POD")
	deployment, _ := os.LookupEnv("KUBERNETES_DEPLOYMENT")
	resp := response{
		Node:       node,
		Namespace:  namespace,
		Pod:        pod,
		Deployment: deployment,
	}

	if prometheusEnabled {
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/monitoring/metrics", promhttp.Handler())
		go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), metricsMux)) }()
	}

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/success", handleSuccess(resp))
	httpMux.HandleFunc("/error", handleError(http.StatusInternalServerError))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), httpMux))
}
