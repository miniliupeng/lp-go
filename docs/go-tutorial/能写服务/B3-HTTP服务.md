# B3 HTTP 服务

Go 1.22 起 `http.NewServeMux` 的模式可以写方法和路径参数。本教程不设置 `GODEBUG=httpmuxgo121=1`。示例在 [examples/b3http](../../../examples/b3http/server.go)。超时字段、`Shutdown` 和 SSE 在 [HTTP 超时与关闭](../工程/D1-HTTP超时与关闭.md)。

## 这一章你要会什么

会用方法和路径参数注册路由，会写不听端口的 `httptest`，会把请求号中间件放在 `WriteHeader` 之前，知道 JSON 编码会多一个换行。

## ServeMux

| 模式 | 行为 |
| --- | --- |
| `GET /healthz` | 200，正文 `ok` |
| `GET /items/{id}` | JSON，`id` 来自 `r.PathValue("id")` |
| `POST /echo` | 解码 JSON 再编码；不合法则 400，正文是固定的 `bad json` |

模式里写了方法之后，方法不对得到 405，并带 `Allow` 头。`POST /healthz` 在测试里就是 405。不写方法的模式 `"/healthz"` 匹配所有方法，这是旧写法。新代码把方法写在模式里。

`{id}` 匹配一段路径，不含斜杠。`{path...}` 吃掉剩余路径，包括斜杠。请求 `/files/a/b` 时 `PathValue("path")` 是 `a/b`。通配符的例子在 D1。名字冲突、模式盖住另一条更具体的模式时，注册会 panic。启动时把路由注册完，让进程在接流量之前失败。

`http.ListenAndServe(":8080", mux)` 使用默认服务器，超时全是零。零的含义是「不限制」，见 D1。能写自己的 `http.Server` 时就把四个时间字段设上。

`Handler` 是：

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

`HandleFunc` 接受 `func(ResponseWriter, *Request)`。函数和方法只要签名对得上，就能当处理器。

## 请求体

`Request.Body` 是 `io.ReadCloser`。读完要 `Close`。服务器会在处理器返回后关掉还没关的 body，自己的代码如果把 body 交给别的函数，仍由那一层负责关。

对 body 加大小限制：`http.MaxBytesReader(w, r.Body, 1<<20)`。不限制时，一个客户端可以把进程的内存打满。JSON 用 `json.NewDecoder(limited).Decode(&dst)`。`Decode` 只解一个值，后面多余的字节要自己决定拒不拒绝。

不要把 `err.Error()` 原样回给外部用户。解析失败返回固定的 `bad json`。内部错误写进日志，响应用 500 和一句固定的话。

## JSON 编码

`json.NewEncoder(w).Encode` 会在末尾加换行。比对正文时把换行算进去，或者只检查包含 `"id":"abc"`。`Encode` 失败往往是因为连接已经断了，这时再 `WriteHeader` 已经来不及，头在第一次 `Write` 时就已经发出。

字段要导出，标签见 [标准库与测试](../地基/A9-标准库与测试.md)。重复键的差别在 [JSON 与文件](B4-JSON与文件.md)。

## 中间件

中间件是 `func(http.Handler) http.Handler`。它在调用 `next` 之前改请求，在 `next` 返回之后改不了已经写出的正文。

[WithRequestID](../../../examples/b3http/server.go) 不启动端口：`httptest.NewRecorder` 加 `NewRequest`。它把请求头 `X-Request-Id` 抄到响应头，没有这个头时写成 `none`。响应头要在 `WriteHeader` 或第一次 `Write` 之前设置。中间件先 `Header().Set`，再调用 `next`。

`NewRecorder` 的 `Code` 在没人写头时是 200。`Body` 是缓冲区，可以直接比较字符串。需要真实连接的行为（刷新、客户端取消）再用 `httptest.NewServer` 或自己 `Listen`，见 D1 的 SSE 测试。

## 处理器里的 context

`r.Context()` 在客户端断开、或服务器的超时取消请求时结束。长时间的循环和下游调用都带上它：

```go
req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
```

下游用 `context.Background()` 会把用户已经取消的请求继续做完。超时预算怎么从入口传到下游，在 [熔断、重试与可观测性](../工程/D4-熔断重试与可观测性.md)。

## 对比陷阱

| 写法 | 实际结果 |
| --- | --- |
| 模式写成 `GET /items/{id}` 却用 `POST` 访问 | 405，不是 404 |
| `{id}` 想匹配 `a/b` | 匹配不到。多段用 `{path...}` |
| 在 `next` 已经 `Write` 之后设置响应头 | 头已经发出，设置无效 |
| 用 `Encode` 的结果和没有换行的字符串全等比较 | 差一个 `\n` |
| 把解码错误的 `Error()` 写进响应 | 内部细节漏给调用方 |
| 处理器里用 `context.Background()` 调下游 | 客户端断开后下游仍继续 |

## 运行

```bash
go test ./examples/b3http -count=1
```

## 练习

1. 加 `DELETE /items/{id}`，返回 204。用 `NewMux().ServeHTTP(rec, req)` 写一个不听端口的测试。
2. 对 `POST /echo` 送一个超过 `MaxBytesReader` 限制的 body，断言状态码不是 200。
3. 在中间件里把 `WriteHeader` 放到 `Header().Set` 之前，看测试里的请求号为什么丢了，再改回来。

## 面试延伸

路由题说 Go 1.22 起模式可以带方法，方法不对是 405，`{id}` 是一段，`{path...}` 是剩余路径。中间件题说头必须在写正文之前。body 题说限制大小、关闭、错误句子不回给用户。context 题说 `r.Context()` 随客户端断开而取消，下游要用它派生请求。

## 问答题

### 问：`POST /healthz` 为什么是 405 不是 404？

模式写成了 `GET /healthz`，路径存在，方法不对。ServeMux 回 405，并带 `Allow`。不写方法的旧模式会匹配所有方法。本教程不设置 `GODEBUG=httpmuxgo121=1`。

### 问：`{id}` 和 `{path...}` 各匹配什么？

`{id}` 是一段，不含斜杠。`{path...}` 吃掉剩余路径。`/files/a/b` 的 `PathValue("path")` 是 `a/b`，测试在 [D1](../工程/D1-HTTP超时与关闭.md)。模式冲突在注册时 panic，让进程在接流量前失败。

### 问：`Encode` 的正文为什么和手写的 JSON 差一个字符？

`json.NewEncoder(w).Encode` 在末尾加换行。测试可以比较整段包含换行的字符串，或者只检查里面有 `"id":"abc"`。第一次 `Write` 之后就不能再改头。

### 问：中间件为什么要先 `Header().Set` 再调用 `next`？

`next` 里一旦 `Write` 或 `WriteHeader`，头就发出了。`WithRequestID` 没有请求头时写成 `none`。这个测试用 `httptest.NewRecorder`，不占端口。

### 问：请求体不限制会怎样？

客户端可以一直发送，进程内存跟着涨。用 `MaxBytesReader` 包住 `r.Body`，再交给 `Decoder`。解码失败回固定的 `bad json`，不回 `err.Error()`。

### 问：处理器里用 `context.Background()` 调下游，问题在哪？

客户端已经断开、入口的 `r.Context()` 已经取消，下游仍会做完。用 `NewRequestWithContext(r.Context(), ...)`。超时字段本身见 D1。

### 问：`NewRecorder` 的 `Code` 默认是多少？

还没有人写状态码时是 200。要断言 204 或 400，必须真的调用 `WriteHeader` 或让 `http.Error` 去写。

### 问：`ListenAndServe` 的超时是多少？

默认服务器四个时间字段都是零，也就是不限制。能构造 `http.Server` 时把超时设上。零值的精确回退在 D1。
