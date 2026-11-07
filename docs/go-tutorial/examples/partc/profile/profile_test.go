package profile

import (
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

func TestCPUProfileNonEmpty(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "cpu.pprof")
	if err != nil {
		t.Fatal(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(50 * time.Millisecond)
	x := 0
	for time.Now().Before(deadline) {
		x += x + 1
	}
	runtime.KeepAlive(x)
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("empty cpu profile")
	}
}

func TestGoroutineLeakProfileRegistered(t *testing.T) {
	p := pprof.Lookup("goroutineleak")
	if p == nil {
		t.Fatal("goroutineleak profile missing")
	}
	if p.Name() != "goroutineleak" {
		t.Fatal(p.Name())
	}
}
