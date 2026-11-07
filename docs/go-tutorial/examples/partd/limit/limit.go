// Package limit 给出两个可测试的小工具：令牌桶和 singleflight。
// 生产环境的令牌桶常用 golang.org/x/time/rate。这里把算法写开，方便对照。
package limit

import (
	"sync"
	"time"
)

// Bucket 是一个令牌桶。rate 是每秒补充的令牌数，burst 是桶的容量。
// 调用方传入 now，测试就可以不睡觉。
type Bucket struct {
	mu     sync.Mutex
	rate   float64
	burst  float64
	tokens float64
	last   time.Time
}

// NewBucket 创建一个装满的桶。rate 和 burst 都必须为正。
func NewBucket(rate, burst float64, now time.Time) *Bucket {
	if rate <= 0 || burst <= 0 {
		panic("limit: rate and burst must be positive")
	}
	return &Bucket{rate: rate, burst: burst, tokens: burst, last: now}
}

// Allow 在 now 这个时刻尝试取走 1 个令牌。
func (b *Bucket) Allow(now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens += elapsed * b.rate
		if b.tokens > b.burst {
			b.tokens = b.burst
		}
		b.last = now
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// Group 把同一个 key 上并发的调用合并成一次。
type Group[T any] struct {
	mu sync.Mutex
	m  map[string]*call[T]
}

type call[T any] struct {
	wg  sync.WaitGroup
	val T
	err error
}

// Do 执行 fn。同一个 key 上，fn 还在跑时，后来的调用等待同一次结果。
func (g *Group[T]) Do(key string, fn func() (T, error)) (T, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[T])
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call[T])
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}
