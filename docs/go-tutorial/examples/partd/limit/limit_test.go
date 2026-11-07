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
