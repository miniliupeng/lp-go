package metrics

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// REDCollector 统计网关请求吞吐量 (Rate)、错误数 (Errors)、请求耗时 (Duration) 与首字延迟 (TTFT)
type REDCollector struct {
	TotalRequests   atomic.Uint64
	TotalErrors     atomic.Uint64
	DurationBuckets [5]atomic.Uint64 // <10ms, <50ms, <200ms, <1000ms, >=1000ms
	mu              sync.RWMutex
	durations       []float64 // 毫秒
	ttfts           []float64 // 首字时延毫秒
}

var DefaultCollector = NewREDCollector()

func NewREDCollector() *REDCollector {
	return &REDCollector{
		durations: make([]float64, 0, 5000),
		ttfts:     make([]float64, 0, 5000),
	}
}

// Record 记录一次 HTTP 请求完成指标
func (c *REDCollector) Record(duration time.Duration, isError bool) {
	c.TotalRequests.Add(1)
	if isError {
		c.TotalErrors.Add(1)
	}

	ms := float64(duration.Microseconds()) / 1000.0

	c.mu.Lock()
	if len(c.durations) < 20000 {
		c.durations = append(c.durations, ms)
	}
	c.mu.Unlock()

	switch {
	case ms < 10:
		c.DurationBuckets[0].Add(1)
	case ms < 50:
		c.DurationBuckets[1].Add(1)
	case ms < 200:
		c.DurationBuckets[2].Add(1)
	case ms < 1000:
		c.DurationBuckets[3].Add(1)
	default:
		c.DurationBuckets[4].Add(1)
	}
}

// RecordTTFT 记录大模型流式输出首字耗时 (Time To First Token)
func (c *REDCollector) RecordTTFT(ttft time.Duration) {
	ms := float64(ttft.Microseconds()) / 1000.0
	c.mu.Lock()
	if len(c.ttfts) < 20000 {
		c.ttfts = append(c.ttfts, ms)
	}
	c.mu.Unlock()
}

// ErrorRate 返回当前错误百分比
func (c *REDCollector) ErrorRate() float64 {
	total := c.TotalRequests.Load()
	if total == 0 {
		return 0.0
	}
	return (float64(c.TotalErrors.Load()) / float64(total)) * 100
}

// P99Duration 返回请求总时延的近似 P99 毫秒值
func (c *REDCollector) P99Duration() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	n := len(c.durations)
	if n == 0 {
		return 0
	}
	tmp := make([]float64, n)
	copy(tmp, c.durations)
	sort.Float64s(tmp)
	idx := int(math.Floor(0.99 * float64(n)))
	if idx >= n {
		idx = n - 1
	}
	return tmp[idx]
}

// AvgTTFT 返回平均首字延迟 (毫秒)
func (c *REDCollector) AvgTTFT() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.ttfts) == 0 {
		return 0.0
	}
	var sum float64
	for _, v := range c.ttfts {
		sum += v
	}
	return sum / float64(len(c.ttfts))
}
