package metricinterfaces

import "time"

type PerformanceMetricsCollector interface {
	GetMetrics() PerformanceMetrics
}

type PerformanceMetrics struct {
	SystemName      string
	UpdateStartTime time.Time
	UpdateDuration  time.Duration
	RenderStartTime time.Time
	RenderDuration  time.Duration
	TotalDuration   time.Duration
}
