package main

import (
	"bytes"
	"sync"
	"testing"
)

// TestChannelAndMutexCorrectness 验证 Channel 行为与 Mutex 并发正确性
func TestChannelAndMutexCorrectness(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 99
	close(ch)
	val, ok := <-ch
	if val != 99 || !ok {
		t.Fatalf("expected 99 and true, got %d, %v", val, ok)
	}
	val, ok = <-ch
	if val != 0 || ok {
		t.Fatalf("expected 0 and false, got %d, %v", val, ok)
	}

	var mu sync.Mutex
	var sum int
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			sum++
			mu.Unlock()
		}()
	}
	wg.Wait()
	if sum != 50 {
		t.Fatalf("expected sum 50, got %d", sum)
	}

	// 验证 sync.Once
	inst1 := GetDatabaseInstance()
	inst2 := GetDatabaseInstance()
	if inst1 != inst2 {
		t.Fatal("expected singleton instances to be identical")
	}

	// 验证 sync.Map
	var sm sync.Map
	sm.Store("k1", "v1")
	if val, ok := sm.Load("k1"); !ok || val != "v1" {
		t.Fatalf("expected v1, got %v", val)
	}
}

// 1. 无对象池：每次新建 bytes.Buffer (高频堆逃逸分配)
func BenchmarkDirectAllocBuffer(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		buf := new(bytes.Buffer)
		buf.WriteString("high performance test payload string")
		_ = buf.Bytes()
	}
}

// 2. 有对象池：基于 sync.Pool 复用 bytes.Buffer (接近 0 分配)
func BenchmarkSyncPoolBuffer(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		buf.WriteString("high performance test payload string")
		_ = buf.Bytes()
		bufferPool.Put(buf)
	}
}

// 3. 并发锁争用基准测试
func BenchmarkMutexContention(b *testing.B) {
	var mu sync.Mutex
	var counter int
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

// 4. sync.Once 已初始化完成后的原子 Fast Path 读取基准
func BenchmarkSyncOnceFastPath(b *testing.B) {
	// 确保已初始化
	_ = GetDatabaseInstance()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetDatabaseInstance()
		}
	})
}

// 5. sync.Map 高并发只读命中基准 (无锁原子指针)
func BenchmarkSyncMapReadThroughput(b *testing.B) {
	var sm sync.Map
	sm.Store("auth_token", "jwt_valid_signature_string")

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			val, _ := sm.Load("auth_token")
			_ = val
		}
	})
}

