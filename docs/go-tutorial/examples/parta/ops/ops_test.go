package ops

import "testing"

func TestBitClear(t *testing.T) {
	// 0b1110 清掉 0b0110 的位，剩下 0b1000。
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
		t.Fatal("&& evaluated the right side")
	}
	if !(true || side()) {
		t.Fatal("or")
	}
	if called {
		t.Fatal("|| evaluated the right side")
	}
}

func TestInc(t *testing.T) {
	if Inc(1) != 2 {
		t.Fatal(Inc(1))
	}
}
