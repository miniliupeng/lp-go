package jsonfile

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"os"
	"path/filepath"
	"testing"
)

type pair struct {
	A int `json:"a"`
}

func TestDuplicateKey(t *testing.T) {
	raw := []byte(`{"a":1,"a":2}`)
	var v1 pair
	if err := json.Unmarshal(raw, &v1); err != nil || v1.A != 2 {
		t.Fatalf("v1 %+v %v", v1, err)
	}
	var v2 pair
	if err := jsonv2.Unmarshal(raw, &v2); err == nil {
		t.Fatalf("v2 accepted duplicate: %+v", v2)
	}
}

func TestRootStaysInside(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { root.Close() })
	got, err := root.ReadFile("a.txt")
	if err != nil || string(got) != "hi" {
		t.Fatalf("%q %v", got, err)
	}
	if _, err = root.ReadFile("../a.txt"); err == nil {
		t.Fatal("path escaped the root")
	}
}
