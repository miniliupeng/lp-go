package profile

import (
	"bytes"
	"reflect"
	"runtime/pprof"
	"testing"
	"time"
)

type box struct{ n int }

func (b *box) Inc() { b.n++ }

func TestTypedNilAndCanSet(t *testing.T) {
	var p *box
	var err error = (*missing)(nil)
	if err == nil {
		t.Fatal("typed nil error compared equal to nil")
	}
	_ = p

	n := 3
	v := reflect.ValueOf(n)
	if v.CanSet() {
		t.Fatal("non-addressable value should not be settable")
	}
	rv := reflect.ValueOf(&n).Elem()
	if !rv.CanSet() {
		t.Fatal("elem of pointer should be settable")
	}
	rv.SetInt(4)
	if n != 4 {
		t.Fatal(n)
	}
}

type missing struct{}

func (m *missing) Error() string { return "missing" }

func TestCPUProfileAndLeakName(t *testing.T) {
	if pprof.Lookup("goroutineleak") == nil {
		t.Fatal("goroutineleak profile is not registered")
	}
	var buf bytes.Buffer
	if err := pprof.StartCPUProfile(&buf); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(50 * time.Millisecond)
	acc := 0
	for time.Now().Before(deadline) {
		acc += Stay(acc)
	}
	pprof.StopCPUProfile()
	if buf.Len() == 0 {
		t.Fatal("empty cpu profile")
	}
	if acc == 0 {
		t.Fatal("loop did not run")
	}
}

func Stay(n int) int { return n + 1 }

func BenchmarkAdd(b *testing.B) {
	x := 0
	for i := 0; i < b.N; i++ {
		x += Stay(i)
	}
	if x == 0 && b.N > 0 {
		b.Fatal(x)
	}
}
