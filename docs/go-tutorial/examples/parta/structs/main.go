// 结构体与方法。运行：
//
//	go run ./docs/go-tutorial/examples/parta/structs
package main

import "fmt"

// Stamp 被嵌入 Task 后，Signature 会提升到 Task 上。
type Stamp struct {
	By string
}

func (s Stamp) Signature() string {
	return "by " + s.By
}

type Task struct {
	Title string
	done  bool
	Stamp
}

func NewTask(title, by string) Task {
	return Task{Title: title, Stamp: Stamp{By: by}}
}

// Complete 要改字段，用指针接收者。
func (t *Task) Complete() {
	t.done = true
}

func (t Task) Status() string {
	state := "todo"
	if t.done {
		state = "done"
	}
	return fmt.Sprintf("%s [%s] (%s)", t.Title, state, t.Signature())
}

func main() {
	task := NewTask("write tutorial", "ada")
	fmt.Println(task.Status())
	task.Complete()
	fmt.Println(task.Status())

	// Go 1.27：字面量的键可以是提升字段的名字，不必再包一层 Stamp。
	promoted := Task{Title: "review", Stamp: Stamp{By: "bea"}}
	also := Task{Title: "review", By: "bea"}
	fmt.Printf("sameBy=%t\n", promoted.By == also.By)
}
