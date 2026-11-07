# Go：从地基到 2026 后端面试

按中文目录和中文文件名读。示例仍在仓库根的 `examples/`，目录名保持英文，用模块 `tutorial` 编译。`go.mod` 写的是 **Go 1.27.1**。在本目录的父级（仓库根）执行 `go test` 时，`GOTOOLCHAIN=auto` 会使用这个版本。

三条版本界线：

- **Go 1.22 起就成立、正文默认使用**：每次循环迭代有自己的循环变量；`for i := range n`；`net/http.ServeMux` 的方法与 `{id}`；`log/slog`；`slices` / `maps`；`errors.Join`；`context.WithCancelCause`。
- **在 Go 1.27.1 上写过并跑过**：泛型方法；结构体字面量用提升字段的名字当键；`encoding/json/v2`；`testing/synctest`；`runtime/pprof` 的 `goroutineleak`；`time` 计时器通道容量为 0。
- **默认关闭**：`GOEXPERIMENT=greenteagc` 才打开 Green Tea GC。本工具链的 `internal/goexperiment/exp_greenteagc_off.go` 里 `GreenTeaGC` 为 false。正文讲的是默认的并发三色标记-清除。

实现细节对照的是 Go 1.27.1 源码。语言规范里的行为不会随小版本改掉；运行时内部字段会。不要把某次 benchmark 的纳秒数背下来。

## A 地基

| 章 | 文件 | 要点 | 示例 |
| --- | --- | --- | --- |
| A1 | [环境与第一个程序](地基/A1-环境与第一个程序.md) | 安装、模块、第一个程序 | `examples/a1hello` |
| A2 | [声明、常量与运算符](地基/A2-声明常量与运算符.md) | 零值、`iota`、`&^` | `examples/a2syntax` |
| A3 | [类型](地基/A3-类型.md) | 字符串、切片、map、指针 | `examples/a3types` |
| A4 | [控制流](地基/A4-控制流.md) | `if`、`for`、`switch` | `examples/a4control` |
| A5 | [函数、结构体与方法](地基/A5-函数结构体与方法.md) | `defer`、方法、嵌入 | `examples/a5struct` |
| A6 | [接口](地基/A6-接口.md) | 隐式实现、方法集、类型断言 | `examples/a6iface` |
| A7 | [错误](地基/A7-错误.md) | `error`、`Is` / `As`、`Join` | `examples/a7err` |
| A8 | [包与模块](地基/A8-包与模块.md) | 导出、`go.mod`、初始化 | `examples/a8calc` |
| A9 | [标准库与测试](地基/A9-标准库与测试.md) | 表驱动测试、`flag`、JSON、时间 | `examples/a9text`、`examples/a9stdlib` |

## B 能写服务

| 章 | 文件 | 要点 | 示例 |
| --- | --- | --- | --- |
| B1 | [并发入门](能写服务/B1-并发入门.md) | goroutine、channel、`WaitGroup` | `examples/b1conc` |
| B2 | [上下文与日志](能写服务/B2-上下文与日志.md) | 取消、`Cause`、`slog` | `examples/b2ctx` |
| B3 | [HTTP 服务](能写服务/B3-HTTP服务.md) | `ServeMux`、`httptest`、中间件 | `examples/b3http` |
| B4 | [JSON 与文件](能写服务/B4-JSON与文件.md) | JSON v1/v2、`os.Root` | `examples/b4json` |
| B5 | [泛型与迭代器](能写服务/B5-泛型与迭代器.md) | 约束、泛型方法、`iter.Seq` | `examples/b5gen` |
| B6 | [目录与检查](能写服务/B6-目录与检查.md) | 目录、`go vet`、`-race` | 用前面的示例 |

## C 运行时

| 章 | 文件 | 要点 | 示例 |
| --- | --- | --- | --- |
| C1 | [内存与逃逸](运行时/C1-内存与逃逸.md) | 逃逸、内存模型、切片扩容 | `examples/c1escape` |
| C2 | [调度器](运行时/C2-调度器.md) | G/M/P、状态、抢占 | `examples/c2sched` |
| C3 | [垃圾回收](运行时/C3-垃圾回收.md) | 三色标记、写屏障、内存上限 | `examples/c3gc` |
| C4 | [通道与选择](运行时/C4-通道与选择.md) | 关闭规则、`select`、计时器通道 | `examples/c4chan` |
| C5 | [同步与原子](运行时/C5-同步与原子.md) | Mutex、原子操作、`Pool`、`sync.Map` | `examples/c5sync` |
| C6 | [接口、反射与剖析](运行时/C6-接口反射与剖析.md) | 方法集、反射、profile | `examples/c6prof` |

## D 工程

| 章 | 文件 | 要点 | 示例 |
| --- | --- | --- | --- |
| D1 | [HTTP 超时与关闭](工程/D1-HTTP超时与关闭.md) | 超时字段、`Shutdown`、通配路径 | `examples/d1http` |
| D2 | [数据库与事务](工程/D2-数据库与事务.md) | 连接池、事务、教学驱动的边界 | `examples/d2sql` |
| D3 | [缓存与限流](工程/D3-缓存与限流.md) | 击穿、穿透、雪崩、令牌桶 | `examples/d3limit` |
| D4 | [熔断、重试与可观测性](工程/D4-熔断重试与可观测性.md) | 熔断、重试预算、指标 | `examples/d4svc` |
| D5 | [短链设计](工程/D5-短链设计.md) | 短链怎么落到数据和锁 | `examples/d5short` |

## E 面试

| 章 | 文件 | 内容 |
| --- | --- | --- |
| E1 | [语言与标准库](面试/E1-语言与标准库.md) | 采分点 |
| E2 | [运行时与并发](面试/E2-运行时与并发.md) | 采分点 |
| E3 | [工程与设计](面试/E3-工程与设计.md) | 采分点 |

零基础从 A1 开始。已经能写业务的人从 B1 扫到 B6，再细读 C 和 D，最后闭卷做 E。

## 命令

在仓库根目录：

```bash
go test ./examples/...
go vet ./examples/...
go run ./examples/a1hello
go run ./examples/a9stdlib -name Ada -upper
```
