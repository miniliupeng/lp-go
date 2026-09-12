package main

import (
	"testing"
)

func TestSafeDivide(t *testing.T) {
	res, err := SafeDivide(100, 4)
	if err != nil || res != 25.0 {
		t.Errorf("100 / 4 预期得到 25.0 无错误，实际得到 %.2f, err=%v", res, err)
	}

	_, errZero := SafeDivide(5, 0)
	if errZero == nil {
		t.Errorf("除数为 0 时必须返回非空错误")
	}
}

func TestCalculateStatistics(t *testing.T) {
	min, max, avg := CalculateStatistics(10, 20, 30)
	if min != 10 || max != 30 || avg != 20 {
		t.Errorf("统计结果异常: min=%.1f, max=%.1f, avg=%.1f", min, max, avg)
	}
}

func TestMakeSequenceGenerator(t *testing.T) {
	seq := MakeSequenceGenerator(10, 5)
	if seq() != 10 {
		t.Errorf("第 1 轮步进预期为 10")
	}
	if seq() != 15 {
		t.Errorf("第 2 轮步进预期为 15")
	}
	if seq() != 20 {
		t.Errorf("第 3 轮步进预期为 20")
	}
}
