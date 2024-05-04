package storage

import (
	"sync"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type (
	CounterValue int64
	GaugeValue   float64
)
type Metric struct {
	Type  MetricType
	Name  string
	Value any
}

type MetricCounter struct {
	Type  MetricType
	Name  string
	Value CounterValue
}
type MetricGauge struct {
	Type  MetricType
	Name  string
	Value GaugeValue
}

type MemStorage struct {
	MetricsCounter map[string]MetricCounter
	MetricsGauge   map[string]MetricGauge
	mutex          sync.Mutex
}

func (s *MemStorage) UpdateMetric(metric Metric) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch metric.Type {
	case Gauge:
		value := GaugeValue(metric.Value.(float64))

		if existingMetric, ok := s.MetricsGauge[metric.Name]; ok {
			existingMetric.Value += value

			s.MetricsGauge[metric.Name] = existingMetric
		} else {
			s.MetricsGauge[metric.Name] = MetricGauge{metric.Type, metric.Name, value}
		}

	case Counter:
		incrementValue := CounterValue(metric.Value.(int64))

		if existingMetric, ok := s.MetricsCounter[metric.Name]; ok {
			existingMetric.Value = incrementValue

			s.MetricsCounter[metric.Name] = existingMetric
		} else {
			s.MetricsCounter[metric.Name] = MetricCounter{metric.Type, metric.Name, incrementValue}
		}
	}

}

func (s *MemStorage) GetMetric(name string, metricType MetricType) (Metric, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var metric Metric
	var metricCounter MetricCounter
	var metricGauge MetricGauge
	var ok bool

	switch metricType {
	case Gauge:
		if metricGauge, ok = s.MetricsGauge[name]; ok {
			metric = Metric{metricGauge.Type, metricGauge.Name, metricGauge.Value}
		}
	case Counter:
		if metricCounter, ok = s.MetricsCounter[name]; ok {
			metric = Metric{metricCounter.Type, metricCounter.Name, metricCounter.Value}
		}
	}
	return metric, ok
}

func IsMetricType(metricType string) bool {
	switch MetricType(metricType) {
	case Gauge, Counter:
		return true
	default:
		return false
	}
}
