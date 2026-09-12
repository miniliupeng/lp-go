# 专题 02：Redis 高并发模式、分布式锁、缓存防击穿与数据一致性

> 本专题深入穿透内存王者 **Redis** 在超高并发场景下的架构实践与核心机理。涵盖 **Redis 7.0+ 数据结构演进（Listpack 彻底淘汰 Ziplist）**、**Go `singleflight` 内存级防击穿实战**、**生产级分布式锁与 Watchdog 看门狗自动续期** 以及 **MySQL-Redis 双写一致性架构**。

---

## 一、 Redis 7.0+ 核心数据结构演进与底层内核

### 1. 为什么 Redis 7.0+ 使用 Listpack 彻底取代了 Ziplist？
- **Ziplist（压缩列表）的致命缺陷 —— 连锁更新（Cascading Updates）**：
  - Ziplist 中每个 entry 节点记录了前一个节点的长度 `prevlen`（如果小于 254 字节占 1 字节，大于等于 254 字节占 5 字节）。
  - **雪崩效应**：若有连续多个长度处于 250~253 字节的节点，一旦插入或修改了其中一个节点导致其超过 254 字节，紧随其后的节点的 `prevlen` 必须从 1 字节扩容为 5 字节；这又导致该节点总长超过 254 字节，连锁触发后续所有节点重新申请内存并拷贝，最坏时间复杂度直接退化为 $O(N^2)$，瞬间导致 Redis 主线程卡死！
- **Listpack（紧凑列表）的破解之道**：
  - Listpack 的每个节点**只记录当前节点自身的长度，不再记录前一个节点的长度**！
  - 从后向前遍历时，通过当前节点末尾的自身长度即可直接跳到当前节点头部，**从物理结构上彻底消除了连锁更新的可能**。

---

## 二、 缓存高并发三剑客终极实战

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. 缓存穿透: 请求一个根本不存在的 Key (直接穿透打死 DB)           │
│    -> 终极解法: 布隆过滤器 (Bloom Filter) + 空值短期缓存 (TTL 30s) │
├─────────────────────────────────────────────────────────────────┤
│ 2. 缓存击穿: 热点 Key 突发过期瞬间，海量并发同时涌入打穿 DB       │
│    -> 终极解法: Go 原生 singleflight 请求合并 (实测合并 99% 压力) │
├─────────────────────────────────────────────────────────────────┤
│ 3. 缓存雪崩: 大量 Key 在同一时间集中过期，或 Redis 集群宕机       │
│    -> 终极解法: 过期时间加随机抖动 (TTL + Rand) + 多级 LocalCache │
└─────────────────────────────────────────────────────────────────┘
```

### 1. Go `singleflight` 内存级防击穿实战
- 当高并发流量涌入请求热点数据时，若缓存刚好失效，如果直接放行，数千个协程会同时向底层 MySQL 发送相同 SQL，瞬间引发数据库连接池耗尽死锁。
- **`SingleFlightGroup` 核心机制**：
  - 基于 `map[string]*call` + `sync.WaitGroup`；
  - 只有**第一个协程**真正发起 DB 查库，其余 999 个并发协程全部挂起等待；
  - 第一个协程拿到结果后，直接将返回值广播共享给所有等待协程，**底层 DB 实际只承受 1 次查询开销**！

---

## 三、 生产级 Redis 分布式锁规范

单纯的 `SET resource key NX PX 30000` 存在两个致命工业级隐患：
1. **误删他人持有的锁（Token 校验）**：
   - 协程 A 业务耗时超过 30s，锁自动过期被释放；协程 B 成功抢到锁；
   - 此时协程 A 执行完毕调用 `DEL`，**会把协程 B 持有的锁误删**！
   - **解法**：Value 必须存入每个客户端独一无二的随机 UUID（Token），解锁时必须通过 **Lua 脚本** 原子核验：
     ```lua
     if redis.call("get", KEYS[1]) == ARGV[1] then
         return redis.call("del", KEYS[1])
     else
         return 0
     end
     ```
2. **长业务锁提前过期（Watchdog 看门狗自动续期）**：
   - 在锁持有期间，启动一个独立的后台后台协程（看门狗），每隔 $\text{TTL} / 3$ 周期性向 Redis 执行续期操作；
   - 业务正常结束或抛出 panic 时，通过 `context.CancelFunc` 停止看门狗并释放锁。

---

## 四、 MySQL 与 Redis 双写一致性架构抉择

| 方案 | 流程 | 致命弊端 | 推荐度 |
| :--- | :--- | :--- | :--- |
| **先删缓存，再写库** | 1. DEL Cache $\to$ 2. UPDATE DB | 并发读请求会在步骤 1 和 2 之间读取旧 DB 并把脏数据重新写入缓存，**产生永久脏数据**！ | ❌ 严禁使用 |
| **延迟双删** | 1. DEL $\to$ 2. UPDATE DB $\to$ 3. 休眠 N 秒 $\to$ 4. 再次 DEL | 延迟时间 $N$ 无法精准估算；且主从复制延迟可能导致再次写入旧值。 | ⚠️ 不推荐 |
| **先写库，再删缓存 (Cache-Aside)** | 1. UPDATE DB $\to$ 2. DEL Cache | 绝大多数场景下足够安全（读操作并发更新缓存概率极低）。但若步骤 2 删除失败，仍可能短期不一致。 | ✅ 常规业务推荐 |
| **Canal 监听 Binlog 异步投递** | UPDATE DB $\to$ Binlog $\to$ Canal/Kafka $\to$ 消费端重试删除缓存 | 业务代码与缓存解耦，消息队列提供重试与 ACK 保障，**最终一致性（Eventual Consistency）的终极工业标准**！ | 🏆 核心高可用系统推荐 |

---

## 五、 本机实测运行输出与基准量化报告

### ① 主程序实测运行输出
```bash
go run ./02-storage-middleware/02-redis-patterns/main.go
```
**实测控制台输出（Apple M1 Pro / darwin-arm64）**：
```text
=== 1. Redis 高并发防击穿：SingleFlight 并发请求合并 ===
[SingleFlight 防击穿实测] 100 个高并发请求涌入，底层数据库实际被查询次数: 1 (成功合并 99% 并发压力！)

=== 2. Redis 分布式锁与看门狗 (Watchdog) 自动续期实战 ===
  [Watchdog 续期] 自动为锁 (Token: client_uuid_abc_123) 刷新 TTL 至 60ms
  [Watchdog 续期] 自动为锁 (Token: client_uuid_abc_123) 刷新 TTL 至 60ms
  [Watchdog 续期] 自动为锁 (Token: client_uuid_abc_123) 刷新 TTL 至 60ms
  [Watchdog 续期] 自动为锁 (Token: client_uuid_abc_123) 刷新 TTL 至 60ms
[RedisLock 解锁] 通过 Lua 比对 Token (client_uuid_abc_123) 校验一致，成功安全释放锁
```

---

### ② 性能基准测试报告
```bash
go test -v -bench=. -benchmem ./02-storage-middleware/02-redis-patterns/...
```
**实测性能数据（现代 `b.Loop()` 规范）**：
```text
=== RUN   TestSingleFlightDedup
--- PASS: TestSingleFlightDedup (0.01s)
=== RUN   TestRedisLockWatchdog
--- PASS: TestRedisLockWatchdog (0.05s)
goos: darwin
goarch: arm64
pkg: lp-go/02-storage-middleware/02-redis-patterns
cpu: Apple M1 Pro
BenchmarkWithoutSingleFlight-8   	546002202	         2.081 ns/op	       0 B/op	       0 allocs/op
BenchmarkSingleFlightDo-8        	  5721300	       204.9 ns/op	      22 B/op	       0 allocs/op
PASS
ok  	lp-go/02-storage-middleware/02-redis-patterns	2.871s
```
> **量化结论**：`SingleFlight` 仅消耗 **204.9 ns** 即可完成并发合并与结果广播，在牺牲微量 CPU 的前提下彻底保护了下游数据库免遭瞬时打穿。
