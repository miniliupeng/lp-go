package chanx

import (
	"testing"
	"time"
)

func TestClosedReceive(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 7
	close(ch)
	v, ok := <-ch
	if v != 7 || !ok {
		t.Fatalf("buffered %d %t", v, ok)
	}
	v, ok = <-ch
	if v != 0 || ok {
		t.Fatalf("drained %d %t", v, ok)
	}
}

func TestSendOnClosedPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	ch := make(chan int)
	close(ch)
	ch <- 1
}

func TestNilCaseNeverSelected(t *testing.T) {
	var blocked chan int
	ready := make(chan int, 1)
	ready <- 1
	select {
	case <-blocked:
		t.Fatal("nil channel was selected")
	case v := <-ready:
		if v != 1 {
			t.Fatal(v)
		}
	}
}

func TestTimerCapIsZero(t *testing.T) {
	timer := time.NewTimer(time.Hour)
	t.Cleanup(func() { timer.Stop() })
	ticker := time.NewTicker(time.Hour)
	t.Cleanup(ticker.Stop)
	if cap(timer.C) != 0 || cap(ticker.C) != 0 {
		t.Fatalf("timer %d ticker %d", cap(timer.C), cap(ticker.C))
	}
}
