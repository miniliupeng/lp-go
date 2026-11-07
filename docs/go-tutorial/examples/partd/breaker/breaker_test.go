package breaker

import (
	"testing"
	"time"
)

func TestBreaker(t *testing.T) {
	now := time.Unix(0, 0)
	clock := now
	b := New(2, time.Second, func() time.Time { return clock })

	if !b.Allow() {
		t.Fatal("closed")
	}
	b.Failure()
	if !b.Allow() {
		t.Fatal("still closed after one failure")
	}
	b.Failure()
	if b.Allow() {
		t.Fatal("should be open")
	}

	clock = clock.Add(time.Second)
	if !b.Allow() {
		t.Fatal("half-open probe")
	}
	if b.Allow() {
		t.Fatal("only one probe")
	}
	b.Failure()
	if b.Allow() {
		t.Fatal("probe failure reopens")
	}

	clock = clock.Add(time.Second)
	if !b.Allow() {
		t.Fatal("second probe")
	}
	b.Success()
	if !b.Allow() {
		t.Fatal("closed again")
	}
}
