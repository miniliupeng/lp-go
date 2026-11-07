// 调用本模块里的另一个包。运行：
//
//	go run ./docs/go-tutorial/examples/parta/calc
package main

import (
	"fmt"

	"lp-go/docs/go-tutorial/examples/parta/calc/mathx"
)

func main() {
	fmt.Printf("add=%d\n", mathx.Add(3, 5))
	n, err := mathx.Clamp(10, 0, 5)
	fmt.Printf("clamp=%d err=%v\n", n, err)
	_, err = mathx.Clamp(1, 3, 1)
	fmt.Printf("bad=%v\n", err)
}
