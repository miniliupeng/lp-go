# Go 教程：从零基础到 2026 大厂后端面试

这是一份**并行新增**的线性教程，目录是 `docs/go-tutorial/`。它不修改、不替换仓库里已有的 `00-` 到 `07-` 专题。那些专题仍然留在原来的位置。本教程正文自洽；学完 Part C、Part D 之后，可以把对应专题当作加练。

读者可以从没写过 Go 开始，一直读到能够按 2026 年一线互联网 / 云厂商的 Go 后端面试来回答：语言语义、运行时、并发、性能、网络、数据一致性和服务端设计。

## 版本假设

示例都在本仓库模块 `lp-go` 里编译。`go.mod` 写的是 **Go 1.27.1**。在这个模块中执行 `go test` / `go run` 时，工具链会按该版本工作（`GOTOOLCHAIN=auto` 时，本机 `go` 命令更旧也会下载 1.27.1）。

阅读时按下面三条区分：

- **Go 1.22 起就成立、本教程默认使用的语义**：`for` 每次迭代有自己的循环变量；`for i := range n`；`net/http.ServeMux` 支持方法和 `{id}` 路径参数；`math/rand/v2`；`log/slog`；`slices` / `maps`；`errors.Join`；`context.WithCancelCause`。
- **本教程用 Go 1.27 写、并在 1.27.1 上跑过的部分**：泛型方法；结构体字面量可以用提升字段的名字当键；`encoding/json/v2`；`testing/synctest`；`runtime/pprof` 的 `goroutineleak`；`time` 的计时器通道固定为同步通道（`cap` 为 0，`asynctimerchan` 这个 GODEBUG 已删除）。
- **默认关闭的实验，面试知道边界即可**：`GOEXPERIMENT=greenteagc`（Green Tea GC）、`GOEXPERIMENT=simd`。正文讲的 GC 是默认的并发三色标记-清除，不是 Green Tea。

文中凡是「实现细节」都对照的是 Go 1.27.1 源码里的注释和类型（`runtime/proc.go`、`runtime/mgc.go`、`runtime/mbarrier.go`、`runtime/chan.go`、`internal/sync/mutex.go`、`internal/runtime/maps`）。实现会随小版本调整，语言规范里的行为不会。不要把某次 benchmark 的纳秒数背成事实。

## 读者与读法

按文件名顺序读。每一章都先读正文，再在**仓库根目录**运行该章命令，最后做章末练习。

| 你现在的情况 | 从哪里读 |
| --- | --- |
| 没写过 Go，或只看过别的语言 | Part A 第 1 章 |
| 能写 Go 业务，没系统看过运行时 | Part B 扫一遍，Part C 细读 |
| 准备 2026 后端面试 | Part C、Part D 细读，Part E 闭卷口述，再回头补答不上的章 |

Part A 不省略安装和第一个程序。Part C 起每章都包含：核心原理、实现要点、面试常问、对比陷阱、可运行示例。

## 目录

### Part A 入门

| 文件 | 要点 | 示例 |
| --- | --- | --- |
| [part-a/01-setup.md](part-a/01-setup.md) | 安装、`PATH`、`go env`、模块模式、`gofmt` / `gopls` | `examples/parta/hello` |
| [part-a/02-program-and-syntax.md](part-a/02-program-and-syntax.md) | 包、`main`、声明、短变量、常量、`iota`、零值、格式化 | `examples/parta/syntax` |
| [part-a/03-types.md](part-a/03-types.md) | 基本类型、转换、数组、切片、map、字符串、指针 | `examples/parta/types` |
| [part-a/04-control-flow.md](part-a/04-control-flow.md) | `if`、`for`、`range`、`switch` | `examples/parta/control` |
| [part-a/05-functions.md](part-a/05-functions.md) | 多返回值、变参、函数值、闭包、`defer` | `examples/parta/functions` |
| [part-a/06-structs-and-methods.md](part-a/06-structs-and-methods.md) | 结构体、方法、接收者、嵌入、构造函数 | `examples/parta/structs` |
| [part-a/07-interfaces.md](part-a/07-interfaces.md) | 隐式实现、小接口、`io.Reader`、类型断言 | `examples/parta/interfaces` |
| [part-a/08-errors.md](part-a/08-errors.md) | `error`、包装、`Is` / `As`、`Join`、`panic` 的边界 | `examples/parta/check` |
| [part-a/09-packages-and-modules.md](part-a/09-packages-and-modules.md) | 导出、`go.mod`、本模块里的子包 | `examples/parta/calc` |
| [part-a/10-stdlib.md](part-a/10-stdlib.md) | `fmt`、字符串、文件、`flag`、JSON、时间、路径 | `examples/parta/stdlib` |
| [part-a/11-concurrency-intro.md](part-a/11-concurrency-intro.md) | goroutine、channel、`WaitGroup` 的正确入门用法 | `examples/parta/concurrency` |
| [part-a/12-testing.md](part-a/12-testing.md) | 表驱动测试、`go test`、示例测试、覆盖率 | `examples/parta/textutil` |
| [part-a/13-mini-tool.md](part-a/13-mini-tool.md) | 把前面串成词频命令行工具 | `examples/parta/wordstat` |

### Part B 进阶实战

| 文件 | 要点 | 示例 |
| --- | --- | --- |
| [part-b/01-cli-files-json.md](part-b/01-cli-files-json.md) | `flag`、`os.Root`、JSON v1 与 Go 1.27 的 `encoding/json/v2` | `examples/partb/httpserver` 的 JSON 与 `examples/partb/generics` 旁的 v2 测试在 `examples/partb/jsondiff` |
| [part-b/02-http-service.md](part-b/02-http-service.md) | `ServeMux`、`httptest`、JSON API、超时字段的位置 | `examples/partb/httpserver` |
| [part-b/03-generics.md](part-b/03-generics.md) | 类型参数、约束、`comparable`、泛型方法及其限制 | `examples/partb/generics` |
| [part-b/04-context-logging.md](part-b/04-context-logging.md) | 取消、超时、`Cause`、`slog`、请求号 | `examples/partb/logctx` |
| [part-b/05-layout-and-quality.md](part-b/05-layout-and-quality.md) | 目录、`gofmt`、`go vet`、`-race`、依赖 | 用前面的示例目录说明 |

### Part C 运行时、并发与性能

| 文件 | 要点 | 示例 |
| --- | --- | --- |
| [part-c/01-memory-escape.md](part-c/01-memory-escape.md) | 栈与堆、逃逸分析、内存模型里的 happens-before | `examples/partc/escape` |
| [part-c/02-scheduler-gmp.md](part-c/02-scheduler-gmp.md) | G/M/P、窃取、netpoller、异步抢占 | `examples/partc/sched` |
| [part-c/03-gc.md](part-c/03-gc.md) | 三色标记、混合写屏障、`GOGC`、`GOMEMLIMIT` | `examples/partc/gcinfo` |
| [part-c/04-channels.md](part-c/04-channels.md) | `hchan`、select 语义、关闭规则、泄漏 | `examples/partc/chanx` |
| [part-c/05-sync-atomic.md](part-c/05-sync-atomic.md) | Mutex 正常/饥饿、`WaitGroup`、`Pool`、原子顺序一致性 | `examples/partc/syncx` |
| [part-c/06-iface-reflect.md](part-c/06-iface-reflect.md) | 方法集、typed nil、反射三定律 | `examples/partc/iface` |
| [part-c/07-pprof-trace-bench.md](part-c/07-pprof-trace-bench.md) | CPU/heap profile、trace、bench、`goroutineleak` | `examples/partc/profile` |

### Part D 工程与系统设计

| 文件 | 要点 | 示例 |
| --- | --- | --- |
| [part-d/01-net-http.md](part-d/01-net-http.md) | 超时字段、连接复用、`Shutdown`、HTTP/2 优先级开关 | `examples/partd/httpx`、`examples/partd/shutdown` |
| [part-d/02-database-transactions.md](part-d/02-database-transactions.md) | `database/sql` 连接池、事务边界、教学用内存驱动的真实限制 | `examples/partd/sqltutorial` |
| [part-d/03-cache-ratelimit.md](part-d/03-cache-ratelimit.md) | singleflight、令牌桶、`sync.Map` 的适用面 | `examples/partd/limit` |
| [part-d/04-service-patterns.md](part-d/04-service-patterns.md) | 超时预算、幂等重试、熔断、优雅退出 | `examples/partd/breaker` |
| [part-d/05-observability.md](part-d/05-observability.md) | `slog`、`runtime/metrics`、pprof 端点 | `examples/partd/obs` |
| [part-d/06-system-design.md](part-d/06-system-design.md) | 短链、限流、任务队列在 Go 里怎么落地 | `examples/partd/shorturl` |

### Part E 面试题库

[part-e/01-interview.md](part-e/01-interview.md)：按主题列题，每题有考察点和采分点。先自己答，再对照。采分点只收本书正文里站得住的内容。

## 怎么运行

在仓库根目录：

```bash
go test ./docs/go-tutorial/...
go run ./docs/go-tutorial/examples/parta/hello
go run ./docs/go-tutorial/examples/parta/wordstat -file docs/go-tutorial/examples/parta/wordstat/testdata/sample.txt -top 3
```

只跑基准：

```bash
go test -bench=. -benchmem -run=^$ ./docs/go-tutorial/examples/partc/profile
```

看逃逸分析（输出措辞可能随编译器小版本变化，测试只检查稳定的 “escapes to heap”）：

```bash
go build -gcflags=-m -o /dev/null ./docs/go-tutorial/examples/partc/escape
```

## 示例清单

| 路径 | 作用 |
| --- | --- |
| `examples/parta/hello` | 第一个程序 |
| `examples/parta/syntax` | 声明、常量、零值 |
| `examples/parta/types` | 切片、map、字符串、指针 |
| `examples/parta/control` | 分支和循环 |
| `examples/parta/functions` | 多返回值、变参、闭包、`defer` |
| `examples/parta/structs` | 方法与嵌入 |
| `examples/parta/interfaces` | 小接口与 `io.Reader` |
| `examples/parta/check` | 错误值、`errors.Is` / `As` |
| `examples/parta/calc` | 子包与测试 |
| `examples/parta/stdlib` | flag、JSON、时间、路径 |
| `examples/parta/concurrency` | goroutine、channel、`WaitGroup` |
| `examples/parta/textutil` | 表驱动测试 |
| `examples/parta/wordstat` | 词频小工具 |
| `examples/partb/httpserver` | HTTP 路由与 `httptest` |
| `examples/partb/generics` | 泛型集合与泛型方法 |
| `examples/partb/jsondiff` | `encoding/json` 与 `encoding/json/v2` 的差异 |
| `examples/partb/logctx` | context 与 `slog` |
| `examples/partc/escape` | 逃逸分析 |
| `examples/partc/sched` | `GOMAXPROCS=1` 时的抢占 |
| `examples/partc/gcinfo` | `SetMemoryLimit` 往返、`MemStats` |
| `examples/partc/chanx` | channel 与 select 的规范行为 |
| `examples/partc/syncx` | Mutex、Once、atomic、Pool、`sync.Map` 注释级用法 |
| `examples/partc/iface` | typed nil、方法集、反射 `CanSet` |
| `examples/partc/profile` | CPU profile、`goroutineleak` 名称、基准 |
| `examples/partd/httpx` | 1.22+ 路由与超时字段的零值含义 |
| `examples/partd/shutdown` | `Server.Shutdown` |
| `examples/partd/sqltutorial` | 可运行的 `database/sql` 事务 |
| `examples/partd/limit` | 令牌桶与 singleflight |
| `examples/partd/breaker` | 熔断状态机 |
| `examples/partd/obs` | JSON `slog` |
| `examples/partd/shorturl` | 带过期时间的短链 |

## 不在默认实现里的东西

面试被问到时，直接说边界：

- Green Tea GC：本工具链用 `GOEXPERIMENT=greenteagc` 才打开，默认文件 `internal/goexperiment/exp_greenteagc_off.go` 里 `GreenTeaGC` 为 false。
- `simd` / `simd/archsimd`：`GOEXPERIMENT=simd`，API 还不稳定。
- `encoding/json` 的 v1 函数仍然可用。Go 1.27 起它的实现换成 v2，但 v1 的语义保留；新代码可以显式用 `encoding/json/v2`。需要旧实现时，构建期设置 `GOEXPERIMENT=nojsonv2`（发布说明写明这是过渡开关）。
