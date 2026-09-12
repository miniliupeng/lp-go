package main

import (
	"testing"
)

func TestEvaluateScore(t *testing.T) {
	if EvaluateScore(95) != "优秀 (Grade A)" {
		t.Errorf("95分评估错误")
	}
	if EvaluateScore(82) != "良好 (Grade B)" {
		t.Errorf("82分评估错误")
	}
	if EvaluateScore(65) != "及格 (Grade C)" {
		t.Errorf("65分评估错误")
	}
	if EvaluateScore(40) != "未达标 (Grade D)" {
		t.Errorf("40分评估错误")
	}
}

func TestFindFirstMatrixIntersection(t *testing.T) {
	matrix := [][]int{
		{10, 20},
		{30, 40},
		{50, 60},
	}

	r, c, found := FindFirstMatrixIntersection(matrix, 40)
	if !found || r != 1 || c != 1 {
		t.Errorf("期望命中坐标 (1, 1)，实际得到 (%d, %d), found=%v", r, c, found)
	}

	_, _, notFound := FindFirstMatrixIntersection(matrix, 999)
	if notFound {
		t.Errorf("元素 999 不在矩阵中，预期 found 应为 false")
	}
}
