package syncx

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestMutexCounter(t *testing.T) {
	var mu sync.Mutex
	n := 0
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100 {
				mu.Lock()
				n++
				mu.Unlock()
			}
		})
	}
	wg.Wait()
	if n != 800 {
		t.Fatal(n)
	}
}

func TestOnce(t *testing.T) {
	var once sync.Once
	n := 0
	for range 5 {
		once.Do(func() { n++ })
	}
	if n != 1 {
		t.Fatal(n)
	}
}

func TestAtomicAdd(t *testing.T) {
	var n atomic.Int64
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100 {
				n.Add(1)
			}
		})
	}
	wg.Wait()
	if n.Load() != 800 {
		t.Fatal(n.Load())
	}
}

func TestPoolNew(t *testing.T) {
	var calls atomic.Int32
	pool := &sync.Pool{New: func() any {
		calls.Add(1)
		return new(int)
	}}
	if pool.Get() == nil || calls.Load() != 1 {
		t.Fatalf("calls=%d", calls.Load())
	}
}

func TestMapStoreLoad(t *testing.T) {
	var m sync.Map
	m.Store("k", 1)
	v, ok := m.Load("k")
	if !ok || v.(int) != 1 {
		t.Fatal(v, ok)
	}
}
