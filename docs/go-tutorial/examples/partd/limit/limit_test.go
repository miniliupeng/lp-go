package limit

import (
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestBucket(t *testing.T) {
	start := time.Unix(0, 0)
	b := NewBucket(10, 2, start)
	if !b.Allow(start) || !b.Allow(start) {
		t.Fatal("burst")
	}
	if b.Allow(start) {
		t.Fatal("third token should be denied")
	}
	if !b.Allow(start.Add(100 * time.Millisecond)) {
		t.Fatal("refill of 1 token after 100ms at 10/s")
	}
}

func TestTTLCacheExpires(t *testing.T) {
	now := time.Unix(0, 0)
	clock := now
	c := NewTTLCache(time.Second, func() time.Time { return clock })
	c.Set("a", "1")
	val, hit, neg := c.Get("a")
	if val != "1" || !hit || neg {
		t.Fatalf("%q %t %t", val, hit, neg)
	}
	c.SetNegative("missing")
	_, hit, neg = c.Get("missing")
	if !hit || !neg {
		t.Fatalf("negative hit=%t neg=%t", hit, neg)
	}
	clock = now.Add(time.Second)
	if _, hit, _ = c.Get("a"); hit {
		t.Fatal("expired positive still hits")
	}
	if _, hit, neg = c.Get("missing"); hit || neg {
		t.Fatal("expired negative still hits")
	}
}

func TestSingleflight(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var g Group[int]
		var calls atomic.Int32
		release := make(chan struct{})
		var wg sync.WaitGroup
		for range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				v, err := g.Do("k", func() (int, error) {
					calls.Add(1)
					<-release
					return 7, nil
				})
				if err != nil || v != 7 {
					t.Errorf("got %d %v", v, err)
				}
			}()
		}
		// 气泡里的 WaitGroup.Wait 和通道接收都是可等待的阻塞。
		// 等到这里时，fn 应该只跑了一次，其余调用堵在同一次结果上。
		synctest.Wait()
		if calls.Load() != 1 {
			t.Fatalf("calls=%d", calls.Load())
		}
		close(release)
		wg.Wait()
	})
}
