package escape

// Leak 把局部变量的地址返回去，编译器必须把它放到堆上。
func Leak() *int {
	n := 1
	return &n
}

// Stay 的参数和结果都不逃出函数。
func Stay(a, b int) int { return a + b }
