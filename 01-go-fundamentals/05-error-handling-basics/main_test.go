package main

import (
	"testing"
)

func TestProcessTransaction(t *testing.T) {
	// 正常扣款
	bal, err := ProcessTransaction(40, 100)
	if err != nil || bal != 60 {
		t.Errorf("扣款逻辑异常: bal=%d, err=%v", bal, err)
	}

	// 负数金额错误
	_, errNeg := ProcessTransaction(-10, 100)
	if errNeg == nil {
		t.Errorf("金额为负数时必须返回错误")
	}

	// 余额不足错误
	_, errOver := ProcessTransaction(150, 100)
	if errOver == nil {
		t.Errorf("超额扣款时必须返回错误")
	}
}

func TestSafeRunWorkerRecovery(t *testing.T) {
	// 验证未发生 panic
	caughtNoPanic := SafeRunWorker(false)
	if caughtNoPanic {
		t.Errorf("未触发 panic 时 intercepted 应当为 false")
	}

	// 验证 panic 被安全捕获
	caughtPanic := SafeRunWorker(true)
	if !caughtPanic {
		t.Errorf("触发 panic 时 intercepted 必须为 true，以证明已被 recover 兜底")
	}
}
