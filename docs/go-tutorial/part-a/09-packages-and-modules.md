# 09 包与模块

## 这一章你要会什么

把代码拆到另一个包，用 `go test` 测它，知道 `go.mod` 里的模块路径就是导入路径的前缀。

## 包

一个目录一个包。包名通常等于目录名，全小写。其他包只能使用大写开头的名字。

`mathx.Add` 和 `mathx.Clamp` 在 [examples/parta/calc/mathx](../examples/parta/calc/mathx/mathx.go)。`main` 包导入它们：

```go
import "lp-go/docs/go-tutorial/examples/parta/calc/mathx"
```

`lp-go` 来自根目录 `go.mod` 的 `module lp-go`。导入路径是「模块路径 + 目录」。

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

`go mod tidy` 会整理 `require`。本教程的示例只用标准库，不会改 `go.mod` 的依赖。

## internal

路径里有 `internal` 目录时，只有它的父目录树能导入它。例如 `lp-go/internal/foo` 只能被模块 `lp-go` 里的代码导入。这是编译器执行的边界，不是约定。

## 初始化顺序

一个包的依赖包先初始化。包内多个文件的 `init` 都跑完，包才算初始化完。如果包 `p` 导入 `q`，`q` 的 `init` 完成 happens before `p` 的 `init` 开始（见 Go 内存模型）。`init` 里做注册可以，不要在里面启动停不下来的 goroutine。

## 练习

在 `mathx` 里加 `Sub`，给它写表驱动测试，再在 `main` 里打印一次。测试写法见下一章的前一章：[测试](12-testing.md)。
