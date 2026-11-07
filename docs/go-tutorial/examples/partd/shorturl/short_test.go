package shorturl

import (
	"errors"
	"testing"
	"time"
)

func TestPutIsIdempotentUntilExpiry(t *testing.T) {
	now := time.Unix(1_000, 0)
	clock := now
	s := New(time.Minute, func() time.Time { return clock })

	a, err := s.Put("https://example.com/a")
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.Put("https://example.com/a")
	if err != nil {
		t.Fatal(err)
	}
	if a != b || a != "1" {
		t.Fatalf("codes %s %s", a, b)
	}
	got, err := s.Resolve(a)
	if err != nil || got != "https://example.com/a" {
		t.Fatalf("resolve %s %v", got, err)
	}

	other, err := s.Put("https://example.com/b")
	if err != nil || other != "2" {
		t.Fatalf("other %s %v", other, err)
	}

	clock = now.Add(time.Minute)
	if _, err = s.Resolve(a); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expired: %v", err)
	}
	again, err := s.Put("https://example.com/a")
	if err != nil {
		t.Fatal(err)
	}
	if again == a {
		t.Fatal("expired url should get a new code")
	}
}
