// Package syntax 放语言里最容易记错的几条运算规则。
package syntax

// BitClear 用 &^ 清掉 y 里为 1 的位。
func BitClear(x, y int) int {
	return x &^ y
}

// Inc 说明 ++ 是语句，不是表达式。
func Inc(n int) int {
	n++
	return n
}
