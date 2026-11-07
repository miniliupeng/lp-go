// Package mathx 是给包与模块一章用的库。
package mathx

import "fmt"

// Add 返回两数之和。
func Add(a, b int) int { return a + b }

// Clamp 把 n 限制在 [lo, hi]。lo > hi 时返回错误。
func Clamp(n, lo, hi int) (int, error) {
	if lo > hi {
		return 0, fmt.Errorf("mathx: lo %d > hi %d", lo, hi)
	}
	if n < lo {
		return lo, nil
	}
	if n > hi {
		return hi, nil
	}
	return n, nil
}
