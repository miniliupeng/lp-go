package task

import "testing"

func TestTaskAndPromotedKey(t *testing.T) {
	item := New("tutorial", "ada")
	if item.Label() != "tutorial [todo] (by ada)" {
		t.Fatal(item.Label())
	}
	item.Finish()
	if item.Label() != "tutorial [done] (by ada)" {
		t.Fatal(item.Label())
	}
	// Go 1.27：字面量键可以是提升字段的名字。
	other := Task{By: "bea", Title: "note"}
	if other.By != "bea" {
		t.Fatal(other.By)
	}
}

func TestDivSumClosureDefer(t *testing.T) {
	q, r, err := Div(7, 2)
	if q != 3 || r != 1 || err != nil {
		t.Fatalf("%d %d %v", q, r, err)
	}
	if _, _, err = Div(1, 0); err == nil {
		t.Fatal("expected divide by zero")
	}
	if Sum(1, 2, 3, 4) != 10 {
		t.Fatal(Sum(1, 2, 3, 4))
	}
	next := Counter()
	if next() != 1 || next() != 2 {
		t.Fatal("counter")
	}
	got := DeferOrder()
	if len(got) != 3 || got[0] != "body" || got[1] != "defer 2" || got[2] != "defer 1" {
		t.Fatal(got)
	}
	arg := DeferArg()
	if len(arg) != 1 || arg[0] != 1 {
		t.Fatal(arg)
	}
	if Recovered() != "boom" {
		t.Fatal(Recovered())
	}
}
