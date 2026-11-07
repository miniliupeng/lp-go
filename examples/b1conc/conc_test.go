package conc

import "testing"

func TestSumSquares(t *testing.T) {
	if SumSquares(4) != 30 || SumSquaresLocked(4) != 30 {
		t.Fatalf("chan=%d lock=%d", SumSquares(4), SumSquaresLocked(4))
	}
}

func TestFirstOnlyOneReady(t *testing.T) {
	data := make(chan string, 1)
	stop := make(chan struct{})
	data <- "ok"
	v, ok := First(data, stop)
	if !ok || v != "ok" {
		t.Fatalf("%q %t", v, ok)
	}
	close(stop)
	v, ok = First(data, stop)
	if ok || v != "" {
		t.Fatalf("%q %t", v, ok)
	}
}
