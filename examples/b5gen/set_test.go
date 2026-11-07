package gen

import "testing"

func TestSetMapMethodIter(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")
	if !s.Has("a") || s.Has("b") {
		t.Fatal("set")
	}
	got := MapSlice([]int{1, 2}, func(n int) int { return n * 2 })
	if len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Fatal(got)
	}
	if (Box{}).Identity(3) != 3 || (Box{}).Identity("x") != "x" {
		t.Fatal("identity")
	}
	if Promoted().X != 7 {
		t.Fatal(Promoted().X)
	}
	var seq []int
	for n := range CountDown(5) {
		seq = append(seq, n)
		if n == 3 {
			break
		}
	}
	if len(seq) != 3 || seq[0] != 5 || seq[2] != 3 {
		t.Fatal(seq)
	}
}
