package jsondiff

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"testing"
)

func TestV1KeepsLastDuplicate(t *testing.T) {
	var m map[string]int
	if err := json.Unmarshal([]byte(`{"a":1,"a":2}`), &m); err != nil {
		t.Fatal(err)
	}
	if m["a"] != 2 {
		t.Fatalf("v1 duplicate = %v", m)
	}
}

func TestV2RejectsDuplicate(t *testing.T) {
	var m map[string]int
	err := jsonv2.Unmarshal([]byte(`{"a":1,"a":2}`), &m)
	if err == nil {
		t.Fatalf("v2 accepted duplicates: %#v", m)
	}
}

func TestV2Marshal(t *testing.T) {
	b, err := jsonv2.Marshal(struct {
		A int `json:"a"`
	}{A: 1})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"a":1}` && string(b) != "{\"a\":1}\n" {
		t.Fatalf("marshal %q", b)
	}
}
