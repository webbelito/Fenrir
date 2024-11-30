package metricinterfaces

import "time"

// PerformanceMetricsCollector is an interface for collecting performance metrics
type PerformanceMetricsCollector interface {
	GetMetrics() PerformanceMetrics
}

// PerformanceMetrics is a struct that holds performance metrics
type PerformanceMetrics struct {
	SystemName      string
	UpdateStartTime time.Time
	UpdateDuration  time.Duration
	RenderStartTime time.Time
	RenderDuration  time.Duration
	TotalDuration   time.Duration
}
