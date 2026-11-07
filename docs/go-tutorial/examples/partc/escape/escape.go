package escape

// Leak 把参数的地址返回去，编译器必须让 n 活得比函数更久。
func Leak(n int) *int {
	return &n
}

// Stay 只返回一个整数，局部量不需要逃到堆上。
func Stay(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}
