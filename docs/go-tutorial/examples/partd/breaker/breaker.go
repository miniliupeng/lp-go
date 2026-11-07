// Package breaker 是一个教学用熔断器：关闭、打开、半开。
// 它没有错误分类、没有并发探针上限以外的保护，不能直接当生产组件。
package breaker

import (
	"sync"
	"time"
)

type state int

const (
	closed state = iota
	open
	halfOpen
)

// Breaker 连续失败达到 threshold 后打开，冷却 cooldown 后允许一次探测。
type Breaker struct {
	mu        sync.Mutex
	state     state
	failures  int
	threshold int
	openUntil time.Time
	cooldown  time.Duration
	now       func() time.Time
	halfInUse bool
}

// New 返回关闭状态的熔断器。now 由调用方注入，测试可以停住时间。
func New(threshold int, cooldown time.Duration, now func() time.Time) *Breaker {
	if threshold < 1 {
		panic("breaker: threshold < 1")
	}
	if now == nil {
		now = time.Now
	}
	return &Breaker{threshold: threshold, cooldown: cooldown, now: now}
}

// Allow 报告现在能不能发出一次调用。半开状态同时只允许一个探测。
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case closed:
		return true
	case open:
		if b.now().Before(b.openUntil) {
			return false
		}
		b.state = halfOpen
		b.halfInUse = true
		return true
	default:
		if b.halfInUse {
			return false
		}
		b.halfInUse = true
		return true
	}
}

// Success 记下一次成功。半开探测成功后回到关闭并清零失败次数。
func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.halfInUse = false
	b.state = closed
}

// Failure 记下一次失败。关闭状态到达阈值就打开；半开探测失败会重新打开。
func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.halfInUse = false
	if b.state == halfOpen || b.state == open {
		b.state = open
		b.openUntil = b.now().Add(b.cooldown)
		return
	}
	b.failures++
	if b.failures >= b.threshold {
		b.state = open
		b.openUntil = b.now().Add(b.cooldown)
	}
}
