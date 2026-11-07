package ops

// BitClear 清掉 x 里那些在 y 中为 1 的位。&^ 是一个运算符，对整数来说结果等于 x & (^y)。
func BitClear(x, y int) int {
	return x &^ y
}

// Inc 演示 ++ 是语句，不能写在表达式里。
func Inc(n int) int {
	n++
	return n
}
