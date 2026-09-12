package main

import (
	"testing"
)

func TestGeneratePseudoUUID(t *testing.T) {
	uuid1 := GeneratePseudoUUID()
	uuid2 := GeneratePseudoUUID()

	if len(uuid1) != 36 {
		t.Errorf("UUID 格式长度预期为 36，实际长度为 %d: %s", len(uuid1), uuid1)
	}

	if uuid1 == uuid2 {
		t.Errorf("两次独立生成的 UUID 预期不相同，实际均为 %s", uuid1)
	}
}
