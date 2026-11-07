# 01 环境与第一个程序

本章假设你能打开终端。示例在本仓库的 **Go 1.27.1** 模块里运行。语言本身从 Go 1.22 起就够用；跟本仓库走，会用到 1.27 的工具链。

## 这一章你要会什么

装上 Go，确认 `go version`，用模块方式跑通 `你好，Go`。

## 安装

从 [https://go.dev/dl/](https://go.dev/dl/) 安装官方发行版。

- Linux 常见做法是把压缩包解到 `/usr/local/go`，把 `/usr/local/go/bin` 加进 `PATH`。
- macOS 可以用安装包。
- Windows 可以用安装包，装完重开终端。

然后执行：

```bash
go version
go env GOROOT GOPATH GOPROXY GOTOOLCHAIN
```

`GOROOT` 是工具链所在目录。`GOPATH` 现在主要放下载的模块缓存和默认的 `bin`，**新代码不要再放到 `GOPATH/src` 里**。从 Go 1.16 起，模块模式就是默认模式。

`GOPROXY` 的默认值是 `https://proxy.golang.org,direct`。在国内如果代理连不上，很多人会改成 `https://goproxy.cn,direct`。这是环境配置，不是语言的一部分。改用户级配置：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

`go env -w` 写的是你这台机器上的用户配置，不要写进业务仓库。

本仓库根目录的 `go.mod` 第一行附近是：

```text
module lp-go

go 1.27.1
```

在这个目录里跑命令时，`GOTOOLCHAIN=auto` 会使用 Go 1.27.1，哪怕系统里的 `go` 命令更旧。教程里的 `go test`、`go fmt` 都在仓库根目录执行。系统自带的旧 `gofmt` 不认识 Go 1.27 的泛型方法，格式化请用 `go fmt`。

## 编辑器

装语言服务器 `gopls`。保存时运行 `gofmt`（编辑器里一般叫 Format on Save）。VS Code 的 Go 插件和 GoLand 都能做这件事。`gofmt` 没有配置项，团队不讨论花括号风格。

## 第一个程序

文件 [examples/parta/hello/main.go](../examples/parta/hello/main.go)：

```go
package main

import "fmt"

func main() {
    fmt.Println("你好，Go")
}
```

在仓库根目录：

```bash
go run ./docs/go-tutorial/examples/parta/hello
```

输出：

```text
你好，Go
```

`package main` 并且有 `func main()` 的包才能当程序入口。`go run` 编译并执行。`go build -o hello ./docs/go-tutorial/examples/parta/hello` 会在当前目录生成二进制文件。

## 常见问题

- `go: command not found`：`PATH` 里没有 Go 的 `bin`。
- 在仓库外建了一个 `go 1.22` 的模块，再复制第 6 章以后带「Go 1.27」标记的代码：泛型方法和提升字段键会编译失败。跟着本仓库的 `go.mod` 即可。

## 练习

1. 把打印的字符串改成你的名字，再 `go run` 一次。
2. 执行 `go env GOVERSION`，确认在仓库根目录看到 `go1.27.1`。
