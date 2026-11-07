package typesdemo

import "testing"

func TestRuneAndByte(t *testing.T) {
	s := "中"
	if len(s) != 3 || len([]rune(s)) != 1 {
		t.Fatalf("bytes=%d runes=%d", len(s), len([]rune(s)))
	}
	if string(rune(65)) != "A" {
		t.Fatal(string(rune(65)))
	}
}

func TestAppendSharesBackingArray(t *testing.T) {
	base := make([]int, 3, 4)
	base[0], base[1], base[2] = 1, 2, 3
	view := base[:2]
	base = append(base, 9)
	if view[0] != 1 || base[0] != 1 || len(base) != 4 {
		t.Fatalf("view=%v base=%v", view, base)
	}
	if &view[0] != &base[0] {
		t.Fatal("expected the same backing array while cap allowed append")
	}
}

func TestMapCommaOk(t *testing.T) {
	m := map[string]int{"ada": 1}
	v, ok := m["ada"]
	if v != 1 || !ok {
		t.Fatalf("%d %t", v, ok)
	}
	v, ok = m["missing"]
	if v != 0 || ok {
		t.Fatalf("missing %d %t", v, ok)
	}
}

func TestThreeIndexAndCopy(t *testing.T) {
	base := []int{1, 2, 3, 4, 5}
	view := base[1:3:4]
	if len(view) != 2 || cap(view) != 3 || view[0] != 2 {
		t.Fatalf("len=%d cap=%d %v", len(view), cap(view), view)
	}
	dst := make([]int, 2)
	n := copy(dst, []int{9, 8, 7})
	if n != 2 || dst[0] != 9 || dst[1] != 8 {
		t.Fatalf("n=%d dst=%v", n, dst)
	}
	var nilSlice []int
	nilSlice = append(nilSlice, 1)
	if len(nilSlice) != 1 {
		t.Fatal(nilSlice)
	}
}

func TestPointer(t *testing.T) {
	n := 10
	p := &n
	*p++
	if n != 11 {
		t.Fatal(n)
	}
}
