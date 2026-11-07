// Package svc 是一个用注入时钟测试的熔断状态机，以及一条 slog 断言。
package svc

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

type Breaker struct {
	mu        sync.Mutex
	threshold int
	cooldown  time.Duration
	failures  int
	state     state
	openUntil time.Time
	probing   bool
}

func NewBreaker(threshold int, cooldown time.Duration) *Breaker {
	if threshold < 1 || cooldown <= 0 {
		panic("svc: threshold and cooldown must be positive")
	}
	return &Breaker{threshold: threshold, cooldown: cooldown}
}

// Allow 在关闭状态放行。打开状态在冷却结束前拒绝；冷却点一到，下一次进入半开并放行一个探测。
func (b *Breaker) Allow(now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == open && !now.Before(b.openUntil) {
		b.state = halfOpen
		b.probing = false
	}
	switch b.state {
	case open:
		return false
	case halfOpen:
		if b.probing {
			return false
		}
		b.probing = true
		return true
	default:
		return true
	}
}

func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = closed
	b.probing = false
}

func (b *Breaker) Failure(now time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == halfOpen {
		b.trip(now)
		return
	}
	b.failures++
	if b.failures >= b.threshold {
		b.trip(now)
	}
}

func (b *Breaker) trip(now time.Time) {
	b.state = open
	b.openUntil = now.Add(b.cooldown)
	b.probing = false
}
