// 声明、常量、零值。运行：
//
//	go run ./docs/go-tutorial/examples/parta/syntax
package main

import "fmt"

// 同一 const 块里，iota 从 0 开始，每行加 1。
const (
	Red int = iota
	Green
	Blue
)

func main() {
	var language string
	year := 2009
	language = "Go"

	const title = "tutorial"

	fmt.Printf("language=%s year=%d title=%s\n", language, year, title)
	fmt.Printf("zeroInt=%d zeroString=%q zeroBool=%t\n", 0, zeroString(), false)
	fmt.Printf("Red=%d Green=%d Blue=%d\n", Red, Green, Blue)
	fmt.Printf("sum=%d\n", add(1, 2))
}

func zeroString() string {
	var s string
	return s
}

func add(a, b int) int {
	return a + b
}
