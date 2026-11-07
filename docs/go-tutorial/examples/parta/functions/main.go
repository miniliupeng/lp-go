// 函数。运行：
//
//	go run ./docs/go-tutorial/examples/parta/functions
package main

import (
	"errors"
	"fmt"
)

func main() {
	q, r, err := div(10, 3)
	fmt.Printf("div=%d rem=%d err=%v\n", q, r, err)

	_, _, err = div(1, 0)
	fmt.Printf("divZero=%v\n", err)

	fmt.Printf("sum=%d\n", sum(1, 2, 3, 4))

	next := counter()
	fmt.Printf("closure=%d\n", next())
	fmt.Printf("closure=%d\n", next())

	fmt.Println("body")
	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
}

func div(a, b int) (int, int, error) {
	if b == 0 {
		return 0, 0, errors.New("divide by zero")
	}
	return a / b, a % b, nil
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}
