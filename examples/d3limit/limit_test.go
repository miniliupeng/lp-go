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
	if !b.Allow(start) || !b.Allow(start) || b.Allow(start) {
		t.Fatal("burst of 2")
	}
	if !b.Allow(start.Add(100 * time.Millisecond)) {
		t.Fatal("expected one refilled token")
	}
}

func TestTTL(t *testing.T) {
	now := time.Unix(0, 0)
	clock := now
	c := NewTTLCache(time.Second, func() time.Time { return clock })
	c.Set("a", "1")
	val, hit, neg := c.Get("a")
	if val != "1" || !hit || neg {
		t.Fatalf("%q %t %t", val, hit, neg)
	}
	c.SetNegative("missing")
	if _, hit, neg = c.Get("missing"); !hit || !neg {
		t.Fatal("negative")
	}
	clock = now.Add(time.Second)
	if _, hit, _ = c.Get("a"); hit {
		t.Fatal("expired")
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
					t.Errorf("%d %v", v, err)
				}
			}()
		}
		synctest.Wait()
		if calls.Load() != 1 {
			t.Fatalf("calls=%d", calls.Load())
		}
		close(release)
		wg.Wait()
		if _, err := g.Do("k", func() (int, error) {
			calls.Add(1)
			return 8, nil
		}); err != nil {
			t.Fatal(err)
		}
		if calls.Load() != 2 {
			t.Fatalf("second call did not run: %d", calls.Load())
		}
	})
}
