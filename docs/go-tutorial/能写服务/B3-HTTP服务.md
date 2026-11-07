# B3 HTTP 服务

Go 1.22 起 `http.NewServeMux` 的模式可以写方法和路径参数。本教程不设置 `GODEBUG=httpmuxgo121=1`。示例在 [examples/b3http](../../../examples/b3http/server.go)。

| 模式 | 行为 |
| --- | --- |
| `GET /healthz` | 200，正文 `ok` |
| `GET /items/{id}` | JSON，`id` 来自 `r.PathValue("id")` |
| `POST /echo` | 解码 JSON 再编码；不合法则 400，正文是固定的 `bad json` |

模式里写了方法之后，`POST /healthz` 得到 405。`{path...}` 吃掉剩余路径，示例在 [D1](../工程/D1-HTTP超时与关闭.md)。

测试用 `httptest.NewServer`，结束时 `Close`。中间件 [WithRequestID](../../../examples/b3http/server.go) 不启动端口：`httptest.NewRecorder` 加 `NewRequest`。它把请求头 `X-Request-Id` 抄到响应头，没有这个头时写成 `none`。响应头要在 `WriteHeader` 之前设置。中间件先 `Header().Set`，再调用 `next`。

`json.NewEncoder(w).Encode` 会在末尾加换行。比对正文时把换行算进去，或者只检查包含 `"id":"abc"`。

生产代码对 body 加大小限制，读完后关闭。不要把 `err.Error()` 原样回给外部用户。超时字段的零值含义在 [D1](../工程/D1-HTTP超时与关闭.md)。

```bash
go test ./examples/b3http -count=1
```

## 练习

加 `DELETE /items/{id}`，返回 204。用 `NewMux().ServeHTTP(rec, req)` 写一个不听端口的测试。
