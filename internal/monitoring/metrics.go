package monitoring

import (
	"sync"
	"time"
)

// Metrics collector for performance monitoring
type Metrics struct {
	mu                sync.RWMutex
	requestCount      int64
	errorCount        int64
	totalDuration     time.Duration
	activeConnections int64
	cacheHits         int64
	cacheMisses       int64
	workerPoolStats   map[string]WorkerPoolMetrics
}

// WorkerPoolMetrics represents metrics for a worker pool
type WorkerPoolMetrics struct {
	JobsProcessed int64         `json:"jobs_processed"`
	JobsFailed    int64         `json:"jobs_failed"`
	AvgDuration   time.Duration `json:"avg_duration"`
	QueueSize     int           `json:"queue_size"`
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	return &Metrics{
		workerPoolStats: make(map[string]WorkerPoolMetrics),
	}
}

// IncrementRequests increments the request counter
func (m *Metrics) IncrementRequests() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestCount++
}

// IncrementErrors increments the error counter
func (m *Metrics) IncrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCount++
}

// AddDuration adds request duration to total
func (m *Metrics) AddDuration(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalDuration += duration
}

// SetActiveConnections sets the number of active connections
func (m *Metrics) SetActiveConnections(count int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeConnections = count
}

// IncrementCacheHits increments cache hit counter
func (m *Metrics) IncrementCacheHits() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheHits++
}

// IncrementCacheMisses increments cache miss counter
func (m *Metrics) IncrementCacheMisses() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheMisses++
}

// UpdateWorkerPoolStats updates worker pool statistics
func (m *Metrics) UpdateWorkerPoolStats(poolName string, stats WorkerPoolMetrics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.workerPoolStats[poolName] = stats
}

// GetSnapshot returns a snapshot of current metrics
func (m *Metrics) GetSnapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	avgDuration := time.Duration(0)
	if m.requestCount > 0 {
		avgDuration = m.totalDuration / time.Duration(m.requestCount)
	}
	
	cacheHitRate := float64(0)
	totalCacheRequests := m.cacheHits + m.cacheMisses
	if totalCacheRequests > 0 {
		cacheHitRate = float64(m.cacheHits) / float64(totalCacheRequests)
	}
	
	workerPoolStats := make(map[string]WorkerPoolMetrics)
	for k, v := range m.workerPoolStats {
		workerPoolStats[k] = v
	}
	
	return MetricsSnapshot{
		RequestCount:      m.requestCount,
		ErrorCount:        m.errorCount,
		ErrorRate:         float64(m.errorCount) / float64(m.requestCount),
		AvgDuration:       avgDuration,
		ActiveConnections: m.activeConnections,
		CacheHitRate:      cacheHitRate,
		CacheHits:         m.cacheHits,
		CacheMisses:       m.cacheMisses,
		WorkerPoolStats:   workerPoolStats,
		Timestamp:         time.Now(),
	}
}

// MetricsSnapshot represents a point-in-time metrics snapshot
type MetricsSnapshot struct {
	RequestCount      int64                        `json:"request_count"`
	ErrorCount        int64                        `json:"error_count"`
	ErrorRate         float64                      `json:"error_rate"`
	AvgDuration       time.Duration                `json:"avg_duration"`
	ActiveConnections int64                        `json:"active_connections"`
	CacheHitRate      float64                      `json:"cache_hit_rate"`
	CacheHits         int64                        `json:"cache_hits"`
	CacheMisses       int64                        `json:"cache_misses"`
	WorkerPoolStats   map[string]WorkerPoolMetrics `json:"worker_pool_stats"`
	Timestamp         time.Time                    `json:"timestamp"`
}
