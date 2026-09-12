package main

import (
	"testing"
)

func TestCalculateFibonacci(t *testing.T) {
	cases := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{8, 21},
	}

	for _, c := range cases {
		actual := CalculateFibonacci(c.n)
		if actual != c.expected {
			t.Errorf("Fibonacci(%d) 期望为 %d，实际得到 %d", c.n, c.expected, actual)
		}
	}
}
