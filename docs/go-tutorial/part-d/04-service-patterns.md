# 04 服务端常见模式

## 超时预算

入口给请求一个 `context.WithTimeout`。往下游调用时传同一个 `ctx`，或者用剩余时间再切一个更短的子 context。下游不看 context，超时就不会发生。见 [context 一章](../part-b/04-context-logging.md)。

不要在每个中间件里各自 `Background()` 再设一个更长的超时，那会把入口的截止时间盖掉。

## 重试

只重试幂等操作：查询、带幂等键的创建、自己可以安全重放的写。POST 创建订单如果没有幂等键，重试会双写。

退避加上抖动，避免所有客户端在同一秒重试。重试预算算进入口超时里：超时到了就停，不要为了「再试一次」新开一个 `context.Background()`。

## 熔断

[breaker](../examples/partd/breaker/breaker.go) 有三个状态：

- **关闭**：请求通过。连续失败达到阈值后打开。
- **打开**：在冷却时间之前拒绝。冷却点一到，下一次 `Allow` 进入半开。
- **半开**：同时只允许一个探测。探测成功就回到关闭并清零失败；探测失败就重新打开，冷却时间从现在再算。

测试用注入的时钟把时间往前拨，不使用 `Sleep`。阈值 2、冷却 1 秒的转移在 `TestBreaker` 里。

这不是生产熔断器。它没有错误分类（4xx 通常不该算失败）、没有半开并发以外的指标、也没有按后端分桶。讲的时候把「状态机」和「还缺什么」一起说。

## 优雅退出

进程收到 `SIGINT` / `SIGTERM` 时：

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()
err := srv.Shutdown(ctx)
```

`Shutdown` 的语义见 [net/http 一章](01-net-http.md)。`NotifyContext` 在信号到来时取消 `ctx`。还要有一个更长的硬超时，避免卡死的 Handler 让进程永远退不掉；那个超时用 `Shutdown` 的 context，而不是 `os.Exit` 跳过收尾。

Hijack 的长连接要单独通知。后台 worker 使用同一个 ctx，退出时把队列里已经拿到的任务处理完或写回。

## 健康检查

存活检查回答「进程还在吗」。就绪检查回答「现在能不能接流量」，例如数据库连通、缓存预热完成。把数据库抖动绑到存活检查上，编排系统会把进程杀掉再拉起，故障被放大。

## 运行

```bash
go test ./docs/go-tutorial/examples/partd/breaker -count=1
```

## 面试常问

- 为什么半开只放一个请求。
- 重试和熔断谁在外层：通常是熔断在外，打开时不要再重试；关闭时对幂等请求做有限次重试。
- `Shutdown` 和 `Close`：`Close` 立即断开，`Shutdown` 先停听、再等活跃连接回到空闲。

## 练习

把熔断阈值改成 1，写出第一次失败之后、冷却结束之前 `Allow` 的返回值。
