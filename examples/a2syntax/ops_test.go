package syntax

import "testing"

func TestBitClear(t *testing.T) {
	if BitClear(0b1110, 0b0110) != 0b1000 {
		t.Fatal(BitClear(0b1110, 0b0110))
	}
}

func TestDivisionTowardZero(t *testing.T) {
	if -1/2 != 0 || -1%2 != -1 {
		t.Fatalf("quot=%d rem=%d", -1/2, -1%2)
	}
}

func TestShortCircuit(t *testing.T) {
	called := false
	side := func() bool {
		called = true
		return true
	}
	if false && side() {
		t.Fatal("and")
	}
	if called {
		t.Fatal("&& ran the right side")
	}
	if !(true || side()) {
		t.Fatal("or")
	}
	if called {
		t.Fatal("|| ran the right side")
	}
}

func TestInc(t *testing.T) {
	if Inc(1) != 2 {
		t.Fatal(Inc(1))
	}
}

func TestIota(t *testing.T) {
	const (
		Red = iota
		Green
		Blue
	)
	if Red != 0 || Green != 1 || Blue != 2 {
		t.Fatalf("%d %d %d", Red, Green, Blue)
	}
}

func TestZeroValues(t *testing.T) {
	var (
		n int
		s string
		b bool
		p *int
	)
	if n != 0 || s != "" || b || p != nil {
		t.Fatalf("%d %q %t %v", n, s, b, p)
	}
}
