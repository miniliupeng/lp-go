# 05 目录、格式化与竞态

## 一个小程序怎么放

本教程的示例为了和章节对应，放在 `docs/go-tutorial/examples` 下。做一个独立服务时，常见布局是：

```text
module/
  go.mod
  cmd/wordstat/main.go    # 只有 main：解析旗标、调用库
  internal/stats/         # 外部模块不能导入
  internal/httpserver/
```

`main` 要薄。可测试的逻辑放在库包里，接受接口（`io.Reader`、`http.Handler`）而不是全局变量。

`internal` 的规则见 [包与模块](../part-a/09-packages-and-modules.md)。

## 提交前的三条命令

在仓库根目录：

```bash
go fmt ./docs/go-tutorial/...
go vet ./docs/go-tutorial/...
go test ./docs/go-tutorial/...
```

`go fmt` 使用 `go.mod` 指定的工具链。不要调用系统里旧版本的 `gofmt` 去格式化含泛型方法的文件。

`go vet` 会抓一类确定的错误，例如把 `*T`（`T` 实现了 `error`）用 `%w` 包起来，导致 `errors.As` 对不上。Go 1.27 的 printf 分析器会报告这种写法。

竞态：

```bash
go test -race ./docs/go-tutorial/...
```

它替代不了代码评审，但能抓住「测试正好撞上」的共享内存写入。

## 依赖

标准库够用就不要加依赖。令牌桶的生产实现常用 `golang.org/x/time/rate`，教程里的 [limit.Bucket](../examples/partd/limit/limit.go) 是为了把算法测清楚，注释里写了这一点。

`go.sum` 是校验和，要提交。不要手改版本号之后不跑 `go mod tidy`。

## 版本怎么写进代码

语言版本由 `go.mod` 的 `go` 行决定。使用比那一行更新的标准库 API，`go test` 会按 Go 1.27 起默认打开的 `stdversion` vet 检查报错。所以「本地工具链很新、go.mod 很旧」时，新 API 会被拒绝。这是故意的。
