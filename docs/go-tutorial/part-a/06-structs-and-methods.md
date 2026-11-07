# 06 结构体与方法

## 这一章你要会什么

用结构体组织数据，用方法挂上行为，分清值接收者和指针接收者。

## 定义和字面量

```go
type Stamp struct {
    By string
}

type Task struct {
    Title string
    done  bool
    Stamp
}
```

字段名大写才能被其他包看到。小写字段仍可在本包使用，外部只能通过你提供的方法改它。

字面量用字段名：`Task{Title: "write tutorial"}`。没写到的字段是零值。

## 方法

```go
func (t Task) Status() string { /* 只读，值接收者 */ }
func (t *Task) Complete()     { t.done = true }
```

接收者就是多出来的第一个参数。要修改接收者，用指针。只读的小结构体可以用值接收者。同一个类型的方法集合里，接收者不要混得没有规律：一旦有方法需要指针，导出的方法往往都用指针，避免调用方困惑。

值 `t` 如果可以取地址，`t.Complete()` 会被编译器改写成 `(&t).Complete()`。放进接口时没有这种自动取地址，见 [接口一章](07-interfaces.md) 和 [接口底层一章](../part-c/06-iface-reflect.md)。

## 嵌入

`Task` 里写下 `Stamp` 而不是 `Stamp Stamp`，这是嵌入。`Task` 会提升 `Stamp` 的字段和方法，所以 `task.Signature()` 和 `task.By` 都能用。这是组合，不是子类。外层如果有同名字段或方法，外层优先。

构造函数是普通函数，惯例叫 `NewTask`。没有语言内置的构造器。

## Go 1.27：提升字段可以当字面量的键

Go 1.27 之前，嵌入字段要写成 `Task{Stamp: Stamp{By: "bea"}}`。从 Go 1.27 起，没有歧义时可以直接写提升出来的名字：

```go
also := Task{Title: "review", By: "bea"}
```

键仍然是一个标识符，不是 `Stamp.By` 这种选择器。经过指针才能到达的字段不能这样写。这是 [examples/parta/structs/main.go](../examples/parta/structs/main.go) 里 `sameBy=true` 的原因。

## 运行

```bash
go run ./docs/go-tutorial/examples/parta/structs
```

输出：

```text
write tutorial [todo] (by ada)
write tutorial [done] (by ada)
sameBy=true
```

## 练习

给 `Task` 加一个 `Reset()` 指针方法，把 `done` 设回 `false`。说明它为什么不能是值接收者。
