package gcinfo

import (
	"runtime"
	"runtime/debug"
	"testing"
)

func TestMemoryLimitRoundTrip(t *testing.T) {
	prev := debug.SetMemoryLimit(-1)
	t.Cleanup(func() { debug.SetMemoryLimit(prev) })
	const want = 512 << 20
	if got := debug.SetMemoryLimit(want); got != prev {
		t.Fatalf("old limit %d, stored %d", prev, got)
	}
	if got := debug.SetMemoryLimit(-1); got != want {
		t.Fatalf("read back %d", got)
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	if ms.Sys == 0 {
		t.Fatal("Sys is 0")
	}
}
