package sched

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

// TestBusyGoroutineYieldsP 把 GOMAXPROCS 设为 1。
// 循环里的 atomic 读取是函数调用，调用点是同步安全点，调度器可以在这里换上另一个 G。
// 完全没有调用的死循环要靠异步抢占（Unix 上是 SIGURG），那部分见教程正文，这个测试不假装自己证明了信号。
func TestBusyGoroutineYieldsP(t *testing.T) {
	prev := runtime.GOMAXPROCS(1)
	t.Cleanup(func() { runtime.GOMAXPROCS(prev) })

	var stop atomic.Bool
	started := make(chan struct{})
	go func() {
		close(started)
		var x int32
		for !stop.Load() {
			x++
		}
		runtime.KeepAlive(x)
	}()
	<-started

	done := make(chan struct{})
	go func() {
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("second goroutine did not run; preemption did not happen")
	}
	stop.Store(true)
}
