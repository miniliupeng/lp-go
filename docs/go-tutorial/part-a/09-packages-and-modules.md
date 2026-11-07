# 09 包与模块

## 这一章你要会什么

把函数放到另一个包里，用模块路径导入它，并说清谁能看见大写名字、`internal` 挡的是谁、`init` 什么时候跑。

## 包和导出

一个目录一个包。包名通常就是目录名，全小写，中间没有空格。`mathx` 这个包在 [examples/parta/calc/mathx](../examples/parta/calc/mathx/mathx.go)，`main` 在上一级目录，所以它们是两个包。

其他包只能使用大写开头的名字。`Add` 和 `Clamp` 可以调用，小写开头的函数、类型和字段从包外面看不见。同一个包里的测试文件（`package mathx`）可以测未导出的名字；`package mathx_test` 则只能走外部能用的 API。本教程的 `mathx` 测试和实现在同一个包。

导入写的是路径，代码里用的是包名：

```go
import "lp-go/docs/go-tutorial/examples/parta/calc/mathx"

fmt.Printf("add=%d\n", mathx.Add(3, 5))
```

## 模块路径从哪来

仓库根目录的 `go.mod` 是：

```text
module lp-go

go 1.27.1
```

`module` 后面的路径是这个模块的前缀。导入路径等于「模块路径 + 从模块根开始的目录」。`go` 那一行是这个模块声明自己使用的语言版本。在本仓库里跑 `go test` 时，`stdversion` 检查会按这一行判断标准库符号算不算太新。这个检查从 Go 1.22 的 vet 起就有，文件顶部的 `//go:build go1.X` 也可以单独抬高该文件的版本。

`go.sum` 记录依赖的校验和。本教程的示例只用标准库，`go mod tidy` 不会给 `go.mod` 增加 `require`。以后加了第三方库，把 `go.sum` 和 `go.mod` 一起提交，不要手改版本号之后不跑 `go mod tidy`。

从 Go 1.16 起，在模块目录里默认就是模块模式。新代码放在含有 `go.mod` 的目录下，不要再放到 `GOPATH/src`。

## 命令

在仓库根目录：

```bash
go test ./docs/go-tutorial/examples/parta/calc/mathx
go run ./docs/go-tutorial/examples/parta/calc
```

`go test` 输出 `ok`。`go run` 输出：

```text
add=8
clamp=5 err=<nil>
bad=mathx: lo 3 > hi 1
```

`Clamp(10, 0, 5)` 把 10 收到上限 5。`Clamp(1, 3, 1)` 因为 `lo > hi` 返回错误，不夹紧。

## internal

导入路径里出现路径元素 `internal` 时，只有 `internal` 的父目录这棵树里的代码能导入它。例如 `lp-go/internal/foo` 只能被模块 `lp-go` 里的包导入；`lp-go/cmd/tool/internal/foo` 只能被 `lp-go/cmd/tool` 下面的包导入。这是编译器执行的规则。服务的库放在 `internal` 里，外部模块就拷不走这个导入路径。

## 初始化顺序

规范规定：一个包先初始化它导入的包，每个包只初始化一次。包自己初始化时，先按依赖给包级变量赋初值，再调用全部 `init`。多个文件时，构建系统按文件名的字典序把它们交给编译器，`init` 就按这个顺序跑。不要把业务正确性押在「这两个文件的 `init` 谁先」上。

Go 内存模型进一步规定：如果包 `p` 导入包 `q`，`q` 的 `init` 全部完成 happens before `p` 的任何初始化开始。`import _ "image/png"` 这种空导入只为了执行对方的 `init`（`image/png` 用它向 `image` 注册解码器）。`init` 里做这种注册可以。不要在 `init` 里启动一个停不下来的 goroutine，也不要在里面连数据库：失败时程序还没进入 `main`，调用方无法处理这个错误。

## 练习

在 `mathx` 里加 `Sub(a, b int) int`，给它写表驱动测试，再在 `main` 里打印一次。测试的写法见 [测试](12-testing.md)。
