package syncx

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestMutexOnceAtomic(t *testing.T) {
	var mu sync.Mutex
	n := 0
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			mu.Lock()
			n++
			mu.Unlock()
		})
	}
	wg.Wait()
	if n != 8 {
		t.Fatal(n)
	}

	var once sync.Once
	calls := 0
	for range 4 {
		once.Do(func() { calls++ })
	}
	if calls != 1 {
		t.Fatal(calls)
	}

	var x atomic.Int64
	x.Add(2)
	if x.Load() != 2 {
		t.Fatal(x.Load())
	}
}

func TestPoolNewOncePerItem(t *testing.T) {
	var news atomic.Int32
	p := &sync.Pool{New: func() any {
		news.Add(1)
		return new(int)
	}}
	v := p.Get().(*int)
	p.Put(v)
	_ = p.Get().(*int)
	if news.Load() != 1 {
		t.Fatalf("New calls %d", news.Load())
	}
}

func TestSyncMapDisjoint(t *testing.T) {
	var m sync.Map
	m.Store("a", 1)
	v, ok := m.Load("a")
	if !ok || v.(int) != 1 {
		t.Fatal(v, ok)
	}
}
