// 控制流。运行：
//
//	go run ./docs/go-tutorial/examples/parta/control
package main

import "fmt"

func main() {
	n := 3
	if n%2 == 1 {
		fmt.Println("odd")
	}

	// if 可以带一条短声明，作用域只在 if/else。
	if v := n * 2; v > 5 {
		fmt.Printf("double=%d\n", v)
	}

	sum := 0
	for i := 1; i <= 3; i++ {
		sum += i
	}
	fmt.Printf("sum=%d\n", sum)

	// Go 1.22 起可以对整数 range。i 依次为 0,1,2。
	var seq []int
	for i := range 3 {
		seq = append(seq, i)
	}
	fmt.Printf("rangeN=%v\n", seq)

	switch n {
	case 1:
		fmt.Println("one")
	case 2, 3:
		fmt.Println("two-or-three")
	default:
		fmt.Println("other")
	}
}
