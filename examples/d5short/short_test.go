package shorturl

import (
	"errors"
	"testing"
	"time"
)

func TestPutResolveExpire(t *testing.T) {
	now := time.Unix(1_000, 0)
	clock := now
	s := New(time.Minute, func() time.Time { return clock })
	a, err := s.Put("https://example.com/a")
	if err != nil || a != "1" {
		t.Fatal(a, err)
	}
	again, err := s.Put("https://example.com/a")
	if err != nil || again != a {
		t.Fatal(again, err)
	}
	got, err := s.Resolve(a)
	if err != nil || got != "https://example.com/a" {
		t.Fatal(got, err)
	}
	b, err := s.Put("https://example.com/b")
	if err != nil || b != "2" {
		t.Fatal(b, err)
	}
	clock = now.Add(time.Minute)
	if _, err = s.Resolve(a); !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}
	fresh, err := s.Put("https://example.com/a")
	if err != nil || fresh == a {
		t.Fatal(fresh, err)
	}
	if _, err = s.Put(""); err == nil {
		t.Fatal("empty")
	}
}
