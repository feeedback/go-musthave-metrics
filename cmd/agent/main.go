package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/feeedback/go-musthave-metrics/internal/storage"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://localhost:8080"
)

var (
	pollCount int64
	memStats  runtime.MemStats
)

var metricsNames = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

func generateRandomValue() float64 {
	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)

	randomFloat := randGen.Float64()

	return randomFloat
}

func sendMetric[ValueType storage.CounterValue | storage.GaugeValue](metricType storage.MetricType, metricName string, metricValue ValueType) {
	url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, metricType, metricName, metricValue)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "text/plain")
	http.DefaultClient.Do(req)
}

func collectMetrics() {
	for {
		runtime.GC()

		runtime.ReadMemStats(&memStats)
		pollCount++

		time.Sleep(pollInterval)
	}
}

func reportMetrics() {
	for {
		for i := 0; i < len(metricsNames); i++ {
			name := metricsNames[i]
			value := reflect.ValueOf(memStats).FieldByName(name).Interface()

			switch value := value.(type) {
			case uint64:
				sendMetric(storage.Gauge, name, storage.GaugeValue(value))
			case uint32:
				sendMetric(storage.Gauge, name, storage.GaugeValue(value))
			case float64:
				sendMetric(storage.Gauge, name, storage.GaugeValue(value))
			case float32:
				sendMetric(storage.Gauge, name, storage.GaugeValue(value))
			default:
				fmt.Printf("Unknown type %T for metric %s\n", value, name)
			}
		}

		sendMetric(storage.Counter, "PollCount", storage.CounterValue(pollCount))
		sendMetric(storage.Gauge, "RandomValue", storage.GaugeValue(generateRandomValue()))

		time.Sleep(reportInterval)
	}
}

func main() {
	go collectMetrics()
	go reportMetrics()

	// Keep the main function running
	select {}
}
