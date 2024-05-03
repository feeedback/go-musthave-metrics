package storage

import "sync"

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	Type  MetricType
	Name  string
	Value float64
}

type MemStorage struct {
	Metrics map[string]Metric
	mutex   sync.Mutex
}

func (s *MemStorage) UpdateMetric(metric Metric) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if existingMetric, ok := s.Metrics[metric.Name]; ok {
		switch metric.Type {
		case Gauge:
			existingMetric.Value = metric.Value
		case Counter:
			existingMetric.Value += metric.Value
		}
		s.Metrics[metric.Name] = existingMetric
	} else {
		s.Metrics[metric.Name] = metric
	}
}

func (s *MemStorage) GetMetric(name string) (Metric, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	metric, ok := s.Metrics[name]
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
