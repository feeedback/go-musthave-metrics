package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	storageModule "github.com/feeedback/go-musthave-metrics/internal/storage"
)

var storage *storageModule.MemStorage

func init() {
	storage = &storageModule.MemStorage{
		MetricsCounter: make(map[string]storageModule.MetricCounter),
		MetricsGauge:   make(map[string]storageModule.MetricGauge),
	}
}

func GetMetricHandler(w http.ResponseWriter, req *http.Request) {
	metricType := req.PathValue("metricType")
	metricName := req.PathValue("metricName")

	metricValue, metricExists := storage.GetMetric(metricName, storageModule.MetricType(metricType))
	if !metricExists {
		http.Error(w, "Metric not exists", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "%v", metricValue)
}

func UpdateMetricHandler(w http.ResponseWriter, req *http.Request) {
	metricTypeRaw := req.PathValue("metricType")
	metricName := req.PathValue("metricName")
	metricValueStr := req.PathValue("metricValue")

	if !storageModule.IsMetricType(metricTypeRaw) {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}
	metricType := storageModule.MetricType(metricTypeRaw)

	var metricValue interface{}
	var err error

	switch metricType {
	case storageModule.Gauge:
		metricValue, err = strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}
	case storageModule.Counter:
		metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}
	}

	if metricName == "" {
		http.Error(w, "Empty metric name", http.StatusNotFound)
		return
	}

	metric := storageModule.Metric{
		Type:  metricType,
		Name:  metricName,
		Value: metricValue,
	}

	storage.UpdateMetric(metric)

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "OK")
}
