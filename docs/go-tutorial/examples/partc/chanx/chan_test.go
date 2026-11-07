package chanx

import (
	"testing"
)

func TestClosedReceive(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 7
	close(ch)
	v, ok := <-ch
	if v != 7 || !ok {
		t.Fatalf("buffered before close: %d %t", v, ok)
	}
	v, ok = <-ch
	if v != 0 || ok {
		t.Fatalf("after drain: %d %t", v, ok)
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

func TestCloseNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	var ch chan int
	close(ch)
}

func TestNilChannelInSelect(t *testing.T) {
	var ch chan int
	select {
	case <-ch:
		t.Fatal("nil channel must not be selected")
	default:
	}
}

func TestSelectPicksAReadyCase(t *testing.T) {
	a := make(chan string, 1)
	b := make(chan string, 1)
	a <- "a"
	b <- "b"
	select {
	case v := <-a:
		if v != "a" {
			t.Fatal(v)
		}
	case v := <-b:
		if v != "b" {
			t.Fatal(v)
		}
	}
}

func TestRangeStopsAfterClose(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	if sum != 3 {
		t.Fatal(sum)
	}
}
