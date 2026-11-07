# 05 目录、格式化与竞态

## 这一章你要会什么

知道一个小服务的文件该放哪，以及提交前哪三条命令必须过。

## 目录

本教程的示例为了和章节对应，放在 `docs/go-tutorial/examples` 下。做一个独立服务时，常见布局是：

```text
module/
  go.mod
  cmd/wordstat/main.go    # 只有 main：解析旗标、调用库
  internal/stats/         # 外部模块不能导入
  internal/httpserver/
```

`main` 要薄。可测试的逻辑放在库包里，函数接受 `io.Reader`、`http.Handler` 这种接口，而不是去读全局变量。`internal` 的导入限制见 [包与模块](../part-a/09-packages-and-modules.md)。

一个目录一个包。测试文件和被测文件放在一起。不要为了「分层」再建 `util`、`common` 这种什么都往里放的包。

## 提交前的三条命令

在仓库根目录：

```bash
go fmt ./docs/go-tutorial/...
go vet ./docs/go-tutorial/...
go test ./docs/go-tutorial/...
```

`go fmt` 使用 `go.mod` 指定的工具链。系统里如果还有旧的 `gofmt`，它不认识 Go 1.27 的泛型方法，会直接报语法错误。格式化用 `go fmt`。

`go vet` 跑一组分析器。Go 1.27 的 `go tool vet help` 列出的名字里，和本教程直接相关的有：

- `printf`：格式串和参数对不上。Go 1.27 起，如果 `%w` 的参数是 `*E`，而 `E` 实现了 `error`，分析器会报告「defeats errors.Is」。源码注释说明这也是为了避免 `errors.As`、`errors.Is` 用错；诊断文本写的是 `errors.Is`。
- `stringintconv`：`string(整数)` 得到的是 Unicode 码点，不是十进制文本。
- `copylocks`：按值复制了 `sync.Mutex` 这类锁。
- `lostcancel`：`context.WithCancel` 返回的 `cancel` 没有被调用。
- `stdversion`：用了比当前文件的 Go 版本更新的标准库符号。文件版本由 `go.mod` 的 `go` 行决定，也可以由文件顶部的 `//go:build go1.X` 抬高。这个分析器在工具链还低于 1.22 时会自己关掉；Go 1.22 及以后的 vet 默认会跑它。所以「本机工具链很新、`go.mod` 的 `go` 行很旧」时，新 API 会被拒绝。

`go vet` 的文档写明它用的是启发式规则，过了 vet 不等于程序正确。

竞态检测：

```bash
go test -race ./docs/go-tutorial/...
```

`-race` 能抓住测试正好撞上的、没有同步的并发写。它替代不了评审：没跑到的路径它看不见。

## 依赖

标准库够用就不要加依赖。令牌桶的生产实现常用 `golang.org/x/time/rate`，教程里的 [limit.Bucket](../examples/partd/limit/limit.go) 把算法写开，是为了能用注入的时钟做测试，源文件开头的注释写了这一点。

`go.sum` 是校验和，要提交。升级依赖的方式是改 `go.mod` 里的版本再执行 `go mod tidy`，让它重写 `go.sum`。

## 练习

在仓库根目录执行 `go vet ./docs/go-tutorial/...`。输出为空、退出码为 0，表示这组分析器没有报告。
