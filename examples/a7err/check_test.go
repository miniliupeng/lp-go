package check

import (
	"errors"
	"fmt"
	"testing"
)

func TestParseAge(t *testing.T) {
	n, err := ParseAge("42")
	if n != 42 || err != nil {
		t.Fatalf("%d %v", n, err)
	}
	_, err = ParseAge("")
	if !errors.Is(err, ErrEmpty) {
		t.Fatal(err)
	}
	_, err = ParseAge("4a")
	var fe *FieldError
	if !errors.As(err, &fe) || fe.Field != "age" {
		t.Fatal(err)
	}
}

func TestJoin(t *testing.T) {
	err := Together(ErrEmpty, fmt.Errorf("other"))
	if !errors.Is(err, ErrEmpty) {
		t.Fatal(err)
	}
}
