# 02 HTTP 服务

## 这一章你要会什么

用 Go 1.22 及以后的 `ServeMux` 注册带方法和路径参数的路由，并用 `httptest` 测它。超时字段的含义放在 [网络一章](../part-d/01-net-http.md)，这里先把服务跑起来。

## 路由

[httpserver.NewMux](../examples/partb/httpserver/server.go) 注册了三条：

| 模式 | 行为 |
| --- | --- |
| `GET /healthz` | 200，正文 `ok` |
| `GET /items/{id}` | JSON `{"id":"..."}`，`id` 来自 `r.PathValue("id")` |
| `POST /echo` | 按 JSON 解码再原样编码；JSON 不合法则 400 |

模式里写了方法之后，同一个路径用错误的方法会得到 **405**。测试对 `POST /healthz` 检查了这一点。`http.NewServeMux` 在 Go 1.22 之前没有方法和通配符；旧行为要用 `GODEBUG=httpmuxgo121=1` 才会回去。本教程不使用那个开关。

`{path...}` 会吃掉剩余路径，示例在 [partd/httpx](../examples/partd/httpx/route_test.go)：`/files/a/b` 的 `path` 是 `a/b`。

## 测试不要听端口

`httptest.NewServer` 会在本地起一个真的服务器，结束时 `Close`。它比自己 `ListenAndServe` 更容易在测试里收干净。

```bash
go test ./docs/go-tutorial/examples/partb/httpserver -count=1
```

## 写 JSON 的顺序

先 `Header().Set`，再 `WriteHeader`，再写正文。`WriteHeader` 之后再改头，改动不会生效。`json.NewEncoder(w).Encode` 会在末尾加换行。测试如果比对整段正文，把换行算进去，或者只检查包含 `"hi"`。

## 中间件

中间件的类型是 `func(http.Handler) http.Handler`。它在调用 `next.ServeHTTP` 之前或之后包一层自己的逻辑。

[WithRequestID](../examples/partb/httpserver/server.go) 把请求头 `X-Request-Id` 抄到响应头；没有这个头时写成 `none`。测试用 `httptest.NewRecorder` 和 `NewRequest`，不启动监听端口。响应头要在 `WriteHeader` 之前设好。中间件里先 `w.Header().Set`，再调用 `next`，`next` 里再写状态码，这个顺序是对的。

## 请求体

生产代码对 body 加大小限制，并在读完后关闭。解码失败返回 400，不要把 `err.Error()` 原样抛给外部用户：里面可能有内部细节。示例为了短，400 的正文是固定的 `bad json`。

## 练习

加 `DELETE /items/{id}`，返回 204。用 `httptest.NewRequest` 写一个不启动服务器的测试：`NewMux().ServeHTTP(rec, req)`。
