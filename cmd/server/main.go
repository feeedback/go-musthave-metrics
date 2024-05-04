package main

import (
	"net/http"

	"github.com/feeedback/go-musthave-metrics/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", handlers.UpdateMetricHandler)
	mux.HandleFunc("GET /{metricType}/{metricName}", handlers.GetMetricHandler)

	http.ListenAndServe(":8080", mux)
}
