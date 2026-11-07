package generics

import "testing"

func TestSet(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")
	s.Add("a")
	if !s.Has("a") || s.Has("b") || s.Len() != 1 {
		t.Fatalf("set len=%d", s.Len())
	}
}

type myInts []int

func TestMapSlice(t *testing.T) {
	got := MapSlice(myInts{1, 2, 3}, func(n int) int { return n * 2 })
	if len(got) != 3 || got[0] != 2 || got[2] != 6 {
		t.Fatalf("%v", got)
	}
}

func TestGenericMethod(t *testing.T) {
	var b Box
	if b.Identity(3) != 3 || b.Identity("go") != "go" {
		t.Fatal(b)
	}
}

func TestCountDownStopsEarly(t *testing.T) {
	var got []int
	for n := range CountDown(5) {
		got = append(got, n)
		if n == 3 {
			break
		}
	}
	if len(got) != 3 || got[0] != 5 || got[2] != 3 {
		t.Fatalf("%v", got)
	}
}

func TestPromotedFieldKey(t *testing.T) {
	type inner struct{ X int }
	type outer struct{ inner }
	o := outer{X: 7}
	if o.X != 7 {
		t.Fatal(o.X)
	}
}
