// Package limit 是令牌桶、带 TTL 的缓存，以及一次调用合并。
// 生产环境的令牌桶常用 golang.org/x/time/rate。这里把规则写开，方便用注入的时钟测试。
package limit

import (
	"sync"
	"time"
)

type Bucket struct {
	mu     sync.Mutex
	rate   float64
	burst  float64
	tokens float64
	last   time.Time
}

func NewBucket(rate, burst float64, now time.Time) *Bucket {
	if rate <= 0 || burst <= 0 {
		panic("limit: rate and burst must be positive")
	}
	return &Bucket{rate: rate, burst: burst, tokens: burst, last: now}
}

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

type entry struct {
	val      string
	negative bool
	expires  time.Time
}

type TTLCache struct {
	mu  sync.Mutex
	ttl time.Duration
	now func() time.Time
	m   map[string]entry
}

func NewTTLCache(ttl time.Duration, now func() time.Time) *TTLCache {
	if ttl <= 0 {
		panic("limit: ttl must be positive")
	}
	if now == nil {
		now = time.Now
	}
	return &TTLCache{ttl: ttl, now: now, m: map[string]entry{}}
}

func (c *TTLCache) Set(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = entry{val: val, expires: c.now().Add(c.ttl)}
}

func (c *TTLCache) SetNegative(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = entry{negative: true, expires: c.now().Add(c.ttl)}
}

func (c *TTLCache) Get(key string) (val string, hit, negative bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.m[key]
	if !ok || !c.now().Before(e.expires) {
		delete(c.m, key)
		return "", false, false
	}
	return e.val, true, e.negative
}

type Group[T any] struct {
	mu sync.Mutex
	m  map[string]*call[T]
}

type call[T any] struct {
	wg  sync.WaitGroup
	val T
	err error
}

// Do 合并同一个 key 上还在执行的调用。fn 返回后记录被删掉，下一次会重新执行。
// fn 不能 panic：panic 发生在 Done 之前时，等待者会一直堵着。
func (g *Group[T]) Do(key string, fn func() (T, error)) (T, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = map[string]*call[T]{}
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
