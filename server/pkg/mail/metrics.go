package mail

import "sync/atomic"

// MetricsTracker tracks email sending metrics
type MetricsTracker struct {
	sentCount     atomic.Int64
	failedCount   atomic.Int64
	templateCache map[string]bool
}

// NewMetricsTracker creates a new metrics tracker
func NewMetricsTracker() *MetricsTracker {
	return &MetricsTracker{
		templateCache: make(map[string]bool),
	}
}

// RecordSent increments the sent counter
func (m *MetricsTracker) RecordSent() {
	m.sentCount.Add(1)
}

// RecordFailed increments the failed counter
func (m *MetricsTracker) RecordFailed() {
	m.failedCount.Add(1)
}

// GetSent returns the total sent count
func (m *MetricsTracker) GetSent() int64 {
	return m.sentCount.Load()
}

// GetFailed returns the total failed count
func (m *MetricsTracker) GetFailed() int64 {
	return m.failedCount.Load()
}

// GetSuccessRate returns the success rate as a percentage
func (m *MetricsTracker) GetSuccessRate() float64 {
	sent := m.sentCount.Load()
	failed := m.failedCount.Load()
	total := sent + failed

	if total == 0 {
		return 100.0
	}

	return float64(sent) / float64(total) * 100
}

// Reset resets all counters
func (m *MetricsTracker) Reset() {
	m.sentCount.Store(0)
	m.failedCount.Store(0)
}
