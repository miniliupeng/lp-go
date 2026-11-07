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

	// 没有操作数时，每个 case 是一个条件，从上到下试到第一个为真的。
	cmd := "sum"
	switch {
	case cmd == "sum":
		fmt.Println("cmd=sum")
	default:
		fmt.Println("cmd=other")
	}

	// fallthrough 只跳到下一个 case 的语句，不会重新判断那个 case 的值。
	switch n {
	case 2, 3:
		fmt.Print("fall=")
		fallthrough
	default:
		fmt.Println("through")
	}

	// break 带标签时离开的是那一层循环，而不是离它最近的一层。
	found := false
outer:
	for a := 0; a < 3; a++ {
		for b := 0; b < 3; b++ {
			if a == 1 && b == 1 {
				found = true
				break outer
			}
		}
	}
	fmt.Printf("brokeOuter=%t\n", found)
}
