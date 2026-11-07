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

func TestPointer(t *testing.T) {
	n := 10
	p := &n
	*p++
	if n != 11 {
		t.Fatal(n)
	}
}
