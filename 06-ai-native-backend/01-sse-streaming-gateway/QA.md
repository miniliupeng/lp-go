# 专题 01：大模型流式输出与 SSE 网关高频面试题 (QA)

### Q1：大模型流式输出场景下，为什么主流大厂普遍选用 SSE 而不是 WebSocket？
**标准深度解析**：
1. **单向传输契约匹配**：大模型问答场景下，客户端发起一次 Prompt 提问后，交互过程本质上是服务端向客户端的持续单向流式推流（One-way streaming）。SSE（Server-Sent Events）原生就是为单向服务端推流设计的，无需 WebSocket 双向握手与自定义帧协议维护。
2. **协议兼容性与网关穿透**：SSE 走标准的 HTTP 协议（`text/event-stream`），天然复用现有的 HTTP/2 多路复用、负载均衡、鉴权 Cookie/Header 机制，以及 CDN/WAF 防火墙，而 WebSocket 往往需要额外的协议升级（`Upgrade: websocket`）和特殊网关转发策略。
3. **断线重连开箱即用**：浏览器原生 `EventSource` 默认支持断线自动重连与 `Last-Event-ID` 机制，降低前端实现复杂度。

---

### Q2：在使用 Nginx 反向代理 SSE 网关时，前端为什么经常遇到“卡顿很久然后一次性吐出所有文字”？如何解决？
**标准深度解析**：
1. **根本原因（Proxy Buffering 缓冲）**：Nginx 等七层网关默认开启了 `proxy_buffering on`。当网关服务推流时，Nginx 会将数据缓存在内存 buffer 中，直到累积满 4KB/8KB 或者连接关闭才统一刷给浏览器，导致流式打字机效果失效，首字时延极长。
2. **彻底解决方案**：
   - **服务端响应头声明**：在 Go 后端 HTTP 响应头中显式添加 `X-Accel-Buffering: no`，指示 Nginx 立即禁用此连接的缓冲转发；
   - **Nginx 配置调整**：在 `location` 中针对流式路由配置 `proxy_buffering off;` 以及 `proxy_cache off;`；
   - **Go 代码层面**：务必在每次 `w.Write()` 产生 Token 后立即调用 `flusher.Flush()`，将 TCP 缓冲区内的数据即刻发送。

---

### Q3：大模型生成可能长达数分钟，若客户端中途直接关闭了浏览器标签页，Go 网关如何避免资源浪费和继续向上游模型白白扣费？
**标准深度解析**：
1. **监听 `r.Context().Done()`**：Go 标准库 `http.Request` 的 Context 绑定了底层连接的生命周期。当客户端断开 TCP 连接时，HTTP Server 会主动取消该 `Context`。
2. **级联取消上游请求**：网关在向上游 LLM 推理集群（如 OpenAI、vLLM、Ollama）发起请求时，必须传入带该生命周期的 Context（`req, _ := http.NewRequestWithContext(r.Context(), ...)`）。
3. **协程快速退出与上游扣费终止**：一旦感知到 `<-ctx.Done()`，网关本地消费协程立即 return 退出，并关闭请求管道；上游推理集群接收到断开信号后停止后续 Token 采样生成，避免算力浪费与不必要的计费。
