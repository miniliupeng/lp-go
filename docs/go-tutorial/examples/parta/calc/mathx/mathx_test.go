package mathx

import "testing"

func TestAdd(t *testing.T) {
	if Add(3, 5) != 8 {
		t.Fatal(Add(3, 5))
	}
}

func TestClamp(t *testing.T) {
	got, err := Clamp(10, 0, 5)
	if err != nil || got != 5 {
		t.Fatalf("clamp high = %d, %v", got, err)
	}
	got, err = Clamp(-1, 0, 5)
	if err != nil || got != 0 {
		t.Fatalf("clamp low = %d, %v", got, err)
	}
	if _, err = Clamp(1, 3, 1); err == nil {
		t.Fatal("expected error")
	}
}
