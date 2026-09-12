package main

import (
	"testing"
)

func TestOperatorsAndBitClear(t *testing.T) {
	// 验证算术与取模
	if 17%5 != 2 {
		t.Errorf("17 %% 5 预期为 2，实际为 %d", 17%5)
	}

	// 验证位清空与权限管理
	perm := PermRead | PermWrite | PermExecute // 1 | 2 | 4 = 7 (0111)

	if !HasPermission(perm, PermWrite) {
		t.Errorf("初始状态必须包含写权限")
	}

	// 剥离写权限
	cleared := ClearPermission(perm, PermWrite) // 0111 &^ 0010 = 0101 (5)

	if HasPermission(cleared, PermWrite) {
		t.Errorf("执行位清空后，写权限必须为 false")
	}

	if !HasPermission(cleared, PermRead) || !HasPermission(cleared, PermExecute) {
		t.Errorf("剥离写权限不得影响读与执行权限")
	}
}
