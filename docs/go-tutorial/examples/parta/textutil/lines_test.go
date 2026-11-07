package textutil

import (
	"fmt"
	"slices"
	"testing"
)

func TestNonEmptyLines(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{name: "empty", in: "", want: nil},
		{name: "blanks", in: " \n\t\n", want: nil},
		{name: "trim", in: " a \n\n b \n", want: []string{"a", "b"}},
		{name: "crlf", in: "a\r\n b", want: []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NonEmptyLines(tt.in)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func ExampleNonEmptyLines() {
	for _, line := range NonEmptyLines("a\n\n b \n") {
		fmt.Println(line)
	}
	// Output:
	// a
	// b
}
