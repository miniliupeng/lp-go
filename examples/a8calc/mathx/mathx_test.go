package mathx

import "testing"

func TestAddAndClamp(t *testing.T) {
	if Add(3, 5) != 8 {
		t.Fatal(Add(3, 5))
	}
	n, err := Clamp(10, 0, 5)
	if n != 5 || err != nil {
		t.Fatalf("%d %v", n, err)
	}
	if _, err = Clamp(1, 3, 1); err == nil {
		t.Fatal("expected error")
	}
}
