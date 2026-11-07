package gcinfo

import (
	"runtime"
	"runtime/debug"
	"testing"
)

func TestMemoryLimitRoundTrip(t *testing.T) {
	prev := debug.SetMemoryLimit(-1)
	t.Cleanup(func() { debug.SetMemoryLimit(prev) })

	const want int64 = 512 << 20
	debug.SetMemoryLimit(want)
	if got := debug.SetMemoryLimit(-1); got != want {
		t.Fatalf("limit = %d", got)
	}
}

func TestReadMemStats(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.Sys == 0 {
		t.Fatal("Sys == 0")
	}
}
