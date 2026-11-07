# 03 缓存与限流

## 缓存

进程内缓存就是一张 map。普通 map 不能并发写。多个请求一起读时，用 `RWMutex`，或者确认访问模式符合 `sync.Map` 的两种场景（只增长，或键按 goroutine 分开）。见 [sync 一章](../part-c/05-sync-atomic.md)。

缓存未命中时很多请求会一起打到数据库，这是击穿。 [singleflight](../examples/partd/limit/limit.go) 把同一个 key 上正在进行的调用合成一次：第一个调用执行函数，后来的调用等这一次的结果。测试用 `testing/synctest`：8 个 goroutine 堵在同一次调用上时，函数体只计数 1 次。`synctest` 里 `WaitGroup.Wait` 和气泡内创建的 channel 是可等待的阻塞，所以这个断言不靠 `Sleep`。

`singleflight` 的这次实现要求函数不要 panic。panic 会让等在后面的调用一直堵着。生产里的 `golang.org/x/sync/singleflight` 还有 `Forget` 和共享结果的细节；这里的泛型版本用来看合并本身。

不要把失败缓存很久。 [TTLCache](../examples/partd/limit/limit.go) 用注入的时钟保存正缓存和负缓存。负缓存表示「这个键确认没有」。到期条件是 `now` 不再早于 `expires`，和短链一样，到期时刻本身就算过期。测试把时钟拨过 TTL 之后，两种条目都变成未命中。负缓存的 TTL 应该很短：把「数据库挂了」缓存五分钟，会把故障延长。

`sync.Pool` 不是缓存。GC 可以清掉它。

## 令牌桶

`Bucket` 有两个参数：每秒补充的令牌数 `rate`，以及桶容量 `burst`。调用方传入 `now`，测试不用睡觉。

规则：按距离上次调用的时间补充 `elapsed * rate`，不超过 `burst`，然后尝试取走 1 个。测试里 `rate=10`、`burst=2`：同一时刻可以过两次，第三次失败；100ms 后补充 1 个，又可以通过。

这是算法演示。生产限流常用 `golang.org/x/time/rate`，它还处理时钟和并发取多个令牌。漏桶限制的是流出速率，桶满就丢；令牌桶允许积攒一定突发。面试时先说清你实现的是哪一种。

限流要放在认证之后还是之前，取决于你要保护什么。防打爆进程可以放在很前面；按用户配额必须先知道用户是谁。

## 运行

```bash
go test ./docs/go-tutorial/examples/partd/limit -count=1
```

## 面试常问

- 击穿、穿透、雪崩分别指什么：击穿是热点 key 失效瞬间打穿；穿透是查一个一定不存在的 key；雪崩是大量 key 同时失效。本教程的代码解决的是击穿里的「同一时刻合并」，不解决穿透。
- 令牌桶为什么能突发：桶里可以预先存 `burst` 个令牌。
- `singleflight` 和缓存的差别：它不保存上一轮的结果。函数返回后，下一次调用会再执行。

## 练习

把桶的 `burst` 改成 1，预测同一时刻第二次 `Allow` 的结果，再改测试。
