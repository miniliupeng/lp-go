// 类型。运行：
//
//	go run ./docs/go-tutorial/examples/parta/types
package main

import "fmt"

func main() {
	// 转换必须显式写出目标类型。
	var n int = 65
	s := string(rune(n))
	fmt.Printf("rune=%s code=%U\n", s, rune(n))

	// len 对字符串数的是字节。中文「中」是 3 个字节、1 个 rune。
	word := "中"
	fmt.Printf("bytes=%d runes=%d\n", len(word), len([]rune(word)))

	nums := []int{1, 2, 3}
	nums = append(nums, 4)
	fmt.Printf("slice=%v len=%d cap=%d\n", nums, len(nums), cap(nums))

	scores := map[string]int{"ada": 1}
	scores["bea"] = 2
	if v, ok := scores["ada"]; ok {
		fmt.Printf("ada=%d ok=%t\n", v, ok)
	}
	_, ok := scores["missing"]
	fmt.Printf("missing=%t\n", ok)

	x := 10
	p := &x
	*p = 11
	fmt.Printf("ptr=%d\n", x)
}
