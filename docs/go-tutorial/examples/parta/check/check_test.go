package check

import (
	"errors"
	"strconv"
	"testing"
)

func TestParseAge(t *testing.T) {
	n, err := ParseAge("20")
	if err != nil || n != 20 {
		t.Fatalf("ParseAge(20) = %d, %v", n, err)
	}

	_, err = ParseAge("  ")
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("empty: got %v", err)
	}

	_, err = ParseAge("abc")
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fatalf("abc: errors.As = %v", err)
	}

	_, err = ParseAge("200")
	var field *FieldError
	if !errors.As(err, &field) {
		t.Fatalf("200: got %v", err)
	}
	if field.Field != "age" || field.Reason != "out of range" {
		t.Fatalf("field = %+v", field)
	}
}
