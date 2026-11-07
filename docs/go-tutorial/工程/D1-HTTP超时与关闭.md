# D1 HTTP 超时与关闭

路由写法见 [HTTP 服务](../能写服务/B3-HTTP服务.md)。这一章讲 `http.Server` 的超时、`Shutdown`，以及 SSE 在客户端断开时怎么停。字段注释来自 Go 1.27.1 的 `net/http/server.go`。示例在 [examples/d1http](../../../examples/d1http/httpx_test.go)。

## 这一章你要会什么

能说出四个超时字段的零值分别回退到谁，能说明 `Shutdown` 不打断正在跑的 Handler，能写一个会 `Flush` 的 SSE，并在 `r.Context()` 结束时停下来。

## 四个时间字段

| 字段 | 文档里的含义 |
| --- | --- |
| `ReadTimeout` | 读完整份请求，包含 body。从接受连接到 body 读完。0 或负数表示不设超时 |
| `ReadHeaderTimeout` | 只限制读请求头。0 时回退到 `ReadTimeout`。两者都是 0 或负数则没有超时 |
| `WriteTimeout` | 写响应的上限。读到下一个请求头时会重置。0 或负数表示不设超时 |
| `IdleTimeout` | keep-alive 等下一个请求。0 时回退到 `ReadTimeout`。如果 `ReadTimeout` 也是 0，则没有空闲超时 |

测试构造的 `Server` 四个字段都是 0，也就是还没设置超时。生产服务不要留着四个 0。慢速的客户端可以一直占着连接。

头和 body 分开限制。`ReadHeaderTimeout` 设得短，挡住迟迟不发送头的连接。body 的上限交给 `http.MaxBytesReader` 和处理器自己的逻辑。一个很短的 `ReadTimeout` 会把合法的大文件上传一起掐掉，因为它把读 body 的时间也算进去。

`WriteTimeout` 从读完请求头开始算写响应的时间。对 SSE 这种要写很久的响应，一个短的 `WriteTimeout` 会把连接掐断。长连接要单独评估这个字段，并仍然观察 `r.Context()`。

`GET /files/{path...}` 会吃掉剩余路径。请求 `/files/a/b` 时 `PathValue("path")` 是 `a/b`。

## Shutdown

`Shutdown` 先关掉所有监听，再关掉空闲连接，然后等待还在处理的连接回到空闲。它不打断正在运行的 Handler。context 先到期时，返回的是 context 的错误。测试里 Handler 堵在 channel 上，100 毫秒的 context 得到 `context.DeadlineExceeded`。放开 Handler 之后，`Serve` 返回 `http.ErrServerClosed`。

文档还写明：调用 `Shutdown` 时，`Serve` 会返回 `ErrServerClosed`。程序要等 `Shutdown` 返回，不要在它还没结束时让进程退出。被 Hijack 的连接（例如自己接管的 WebSocket）不在等待范围内，要另外通知那些处理器退出。`Shutdown` 之后这台 `Server` 不能再 `Serve`。

收到进程信号时的顺序是：`signal.Notify` 等到 `SIGINT` 或 `SIGTERM`，用一个有截止时间的 context 调用 `Shutdown`，同时把这个取消传给后台任务。处理器里如果堵在自己的 channel 上而不看 `r.Context()`，`Shutdown` 只能干等，直到你给的 context 到期。

客户端一侧，`http.Client` 也要有超时。没有 `Timeout`、请求也没有 context 时，对端不响应，调用方会一直等。`Client.Timeout` 包含连接、重定向和读 body。更细的控制用 `Transport` 的 `DialContext`、`TLSHandshakeTimeout`、`ResponseHeaderTimeout`。请求级截止时间用 `NewRequestWithContext`。

## SSE

服务器推送事件（SSE）是一条 `Content-Type: text/event-stream` 的响应。每条事件是 `data: ...` 再加一个空行。中间代理和 `http.ResponseWriter` 的缓冲会把小块数据留住，调用方要断言写出的接口实现了 `http.Flusher`，写完一条就 `Flush`。

客户端断开时，`r.Context()` 会结束。循环里在 `Flush` 之后观察这个 context，结束就 `return`。测试 `TestSSEStopsWhenClientCancels` 写了一条 `data: one`，读到之后取消请求，处理器必须从 `<-r.Context().Done()` 返回。

反压是指下游消费得比你生产得慢。SSE 的 `Write` 会在连接的发送缓冲区满时阻塞，这是传输层给你的反压。不要再开一个无界 channel 把事件堆在内存里。生产者往一个有界 channel 里送，channel 满了就等待或丢弃，并在 context 取消时离开。有界 channel 当信号量的 happens-before 在 [内存与逃逸](../运行时/C1-内存与逃逸.md)。

`Flush` 只能把数据推出 `net/http` 这一层。中间如果还有反向代理，代理自己的缓冲也要关。那是部署配置，不是 `Flusher` 能保证的。

## 对比陷阱

| 说法 | 实际情况 |
| --- | --- |
| `ReadTimeout` 只限制请求头 | 它包含 body。只限制头用 `ReadHeaderTimeout` |
| 零值表示用一个默认的几秒 | 0 或负数表示不设这个超时。`ReadHeaderTimeout` 和 `IdleTimeout` 会回退 |
| `Shutdown` 会取消正在跑的 Handler | 不取消。Handler 自己看 `r.Context()`，`Shutdown` 等它返回 |
| `Serve` 在 `Shutdown` 之后还要再关一次 | `Shutdown` 被调用后 `Serve` 返回 `ErrServerClosed`。等的是 `Shutdown` 本身 |
| Hijack 的连接会被 `Shutdown` 等完 | 文档写明不等。要另开退出路径 |
| SSE 写了 `Write` 客户端就能立刻看见 | 还要 `Flush`。代理仍可能再缓冲 |
| 客户端断开后循环可以继续生成 | 生成会堆内存。看 `r.Context()`，队列要有界 |

## 运行

```bash
go test ./examples/d1http -count=1
```

## 练习

1. 给测试里的 `Server` 设置 `ReadHeaderTimeout` 为 1 秒，断言它不再是零值。不要在测试里真的等这个超时。
2. 读 `TestSSEStopsWhenClientCancels`，指出哪一行把事件推给客户端，哪一行让处理器停下来。
3. 把 SSE 处理器改成忽略 `r.Context()`、只在一个永远等不到的 channel 上阻塞。预测 `Shutdown` 会等到 context 超时，然后把代码改回去。

## 面试延伸

超时题先画回退关系：`ReadHeaderTimeout` 为 0 时用 `ReadTimeout`，`IdleTimeout` 为 0 时也用 `ReadTimeout`，`ReadTimeout` 含 body。关闭题说 `Shutdown` 关监听和空闲连接，不打断 Handler，context 到期返回 context 的错误，`Serve` 得到 `ErrServerClosed`，Hijack 的连接要另说。SSE 题说 `text/event-stream`、`Flush`、`r.Context()`，以及用有界队列做反压，而不是无界地堆事件。

## 问答题

### 问：四个字段都是 0 时，服务器有超时吗？

没有。`TestWildcardAndZeroTimeouts` 断言这四个字段都是 0。0 或负数表示不设这个限制。`ReadHeaderTimeout` 为 0 时才回退到 `ReadTimeout`，而 `ReadTimeout` 也是 0，所以读头也没有超时。`IdleTimeout` 同样回退。生产不要留着四个 0。

### 问：`ReadTimeout` 为什么会误伤大文件上传？

它从接受连接一直算到 body 读完。合法的大文件可能超过一个很短的 `ReadTimeout`。头用 `ReadHeaderTimeout` 挡住迟迟不发送的客户端，body 的上限用 `MaxBytesReader`。

### 问：`WriteTimeout` 对 SSE 有什么影响？

它限制写响应的时间，读到下一个请求头时才重置。一条要写很久的 SSE 会被短的 `WriteTimeout` 掐断。长连接要单独评估这个字段，并仍然看 `r.Context()`。

### 问：`Shutdown` 会取消 Handler 吗？

不会。它关监听和空闲连接，然后等活跃连接回到空闲。测试里 Handler 堵在 channel 上，100 毫秒的 context 得到 `DeadlineExceeded`。放开 Handler 之后，`Serve` 返回 `ErrServerClosed`。程序要等 `Shutdown` 返回，不能在它还没结束时退出进程。

### 问：Hijack 的连接谁来关？

`Shutdown` 不等它们。WebSocket 这类自己接管的连接要另有退出信号。`Shutdown` 之后这台 `Server` 不能再 `Serve`。

### 问：客户端断开时 SSE 测试在断言什么？

`TestSSEStopsWhenClientCancels` 先 `Flush` 出 `data: one\n\n`，客户端读到这 11 个字节后取消请求。处理器堵在 `<-r.Context().Done()`，必须在两秒内返回。不看 context 的循环会把 `Shutdown` 拖到超时。

### 问：反压为什么不能用无界 channel？

生产比消费快时，无界队列把事件堆在内存里。`Write` 在发送缓冲区满时会阻塞，这是传输层的反压。应用侧再用有界 channel，满了就等或丢，并在 context 取消时离开。代理自己的缓冲不是 `Flusher` 能关掉的。

### 问：`http.Client` 不加超时会怎样？

对端不响应时，调用方一直等。`Client.Timeout` 包含连接、重定向和读 body。单个请求的截止时间用 `NewRequestWithContext`，并把入口剩下来的时间传下去。
