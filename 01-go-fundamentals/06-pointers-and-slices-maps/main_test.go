package main

import (
	"testing"
)

func TestUpdateValueByPointer(t *testing.T) {
	num := 20
	UpdateValueByPointer(&num, 5)
	if num != 100 {
		t.Errorf("通过指针修改后预期为 100，实际为 %d", num)
	}
}

func TestFilterEvenNumbers(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	output := FilterEvenNumbers(input)

	if len(output) != 3 {
		t.Fatalf("偶数切片长度预期为 3，实际为 %d", len(output))
	}
	if output[0] != 2 || output[1] != 4 || output[2] != 6 {
		t.Errorf("偶数提取结果不符: %v", output)
	}
}

func TestMapOperations(t *testing.T) {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2

	if val, ok := m["a"]; !ok || val != 1 {
		t.Fatalf("comma-ok 读取 a 失败")
	}

	delete(m, "a")
	if _, ok := m["a"]; ok {
		t.Fatalf("delete 删除后 a 依然存在")
	}
}

