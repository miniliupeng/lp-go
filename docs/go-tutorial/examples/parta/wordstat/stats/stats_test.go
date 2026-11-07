package stats

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountAndTop(t *testing.T) {
	freq, err := Count(strings.NewReader("Go, go!\n"))
	if err != nil {
		t.Fatal(err)
	}
	if freq["go"] != 2 {
		t.Fatalf("freq = %#v", freq)
	}
	top := Top(freq, 1)
	if len(top) != 1 || top[0].Word != "go" || top[0].Count != 2 {
		t.Fatalf("top = %#v", top)
	}
	if len(Top(freq, 0)) != 0 {
		t.Fatal("n=0")
	}
}

func TestSampleFile(t *testing.T) {
	f, err := os.Open(filepath.Join("..", "testdata", "sample.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	freq, err := Count(f)
	if err != nil {
		t.Fatal(err)
	}
	top := Top(freq, 3)
	want := []Item{
		{Word: "go", Count: 2},
		{Word: "is", Count: 2},
		{Word: "practical", Count: 2},
	}
	if len(top) != len(want) {
		t.Fatalf("top = %#v", top)
	}
	for i := range want {
		if top[i] != want[i] {
			t.Fatalf("top[%d] = %#v want %#v", i, top[i], want[i])
		}
	}
}
