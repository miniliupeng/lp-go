// Package task 演示结构体、方法和 Go 1.27 的提升字段字面量键。
package task

import "fmt"

type Author struct {
	By string
}

type Task struct {
	Author
	Title string
	Done  bool
}

func New(title, by string) Task {
	return Task{Author: Author{By: by}, Title: title}
}

func (t Task) Label() string {
	state := "todo"
	if t.Done {
		state = "done"
	}
	return fmt.Sprintf("%s [%s] (by %s)", t.Title, state, t.By)
}

func (t *Task) Finish() {
	t.Done = true
}

func Div(a, b int) (int, int, error) {
	if b == 0 {
		return 0, 0, fmt.Errorf("divide by zero")
	}
	return a / b, a % b, nil
}

func Sum(nums ...int) int {
	n := 0
	for _, v := range nums {
		n += v
	}
	return n
}

func Counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// DeferOrder 用命名返回值，defer 里对 out 的修改调用方能看见。
// 多个 defer 后进先出，所以顺序是 body、defer 2、defer 1。
func DeferOrder() (out []string) {
	out = append(out, "body")
	defer func() { out = append(out, "defer 1") }()
	defer func() { out = append(out, "defer 2") }()
	return out
}

// DeferArg 的参数在执行到 defer 这一行时求值。
func DeferArg() (out []int) {
	n := 1
	defer func(v int) { out = append(out, v) }(n)
	n = 2
	return
}

// Recovered 在 defer 里直接调用 recover，并把命名返回值交给调用方。
func Recovered() (msg string) {
	defer func() {
		if v := recover(); v != nil {
			msg = fmt.Sprint(v)
		}
	}()
	panic("boom")
}
