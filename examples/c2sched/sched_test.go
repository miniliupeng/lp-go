package sched

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

// TestBusyGoroutineYieldsP 在 GOMAXPROCS=1 时，第二个 goroutine 仍能跑完。
// 忙循环调用了 atomic.Bool.Load，这是一次函数调用，因而是同步安全点。
// 它不能证明抢占信号是 SIGURG。Unix 上异步抢占用的信号在 runtime/signal_unix.go 里是 sigPreempt = _SIGURG。
func TestBusyGoroutineYieldsP(t *testing.T) {
	old := runtime.GOMAXPROCS(1)
	t.Cleanup(func() { runtime.GOMAXPROCS(old) })
	var stop atomic.Bool
	done := make(chan struct{})
	go func() {
		for !stop.Load() {
		}
	}()
	go func() {
		stop.Store(true)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("second goroutine did not run")
	}
}
