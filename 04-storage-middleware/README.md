# 阶段四: 存储与高并发中间件

高并发场景下，后端架构的极限瓶颈几乎 100% 发生在存储与分布式中间件层。本阶段深入 MySQL InnoDB 锁机制、Redis 高性能并发模式与 Kafka 分布式流式治理，掌握企业级稳定性核心命脉。

---

## 💡 核心工程心智与底层机制

1. **MySQL InnoDB 事务与锁机制深潜**：
   - 彻底吃透 RR 隔离级别下的快照读（MVCC ReadView + undo log）与当前读（Current Read）；
   - 掌握记录锁（Record Lock）、间隙锁（Gap Lock）与临键锁（Next-Key Lock）在主键、唯一索引与非唯一索引上的加锁判定规则；
   - 深入分析并发插入冲突引发的死锁现场，掌握死锁日志提取与索引覆盖优化。
2. **Redis 高并发分布式设计模式**：
   - 掌握基于 SETNX + Lua 脚本原子校验释放的分布式互斥锁，实现协程异步看门狗自动续期防提早释放；
   - 运用 `golang.org/x/sync/singleflight` 在应用进程内聚合并发请求，实现热点 Key 失效时零击穿数据库；
   - 深入 Cache-Aside 旁路缓存模式下的双写一致性（延时双删与 Canal 监听 binlog 异步刷新权衡）。
3. **Kafka KRaft 分布式消息引擎**：
   - 拥抱去 ZooKeeper 化的现代 KRaft 共识架构（基于 Raft 协议的 Quorum 控制器）；
   - 打造端到端消息“零丢失（Zero-Loss）”铁律：生产者 `acks=all` + 幂等开启（`enable.idempotence=true`）、Broker `min.insync.replicas=2`、消费者手动提交偏移量；
   - 建立业务级消息去重表与滑动窗口幂等防重机制。

---

## 🗺️ 阶段专题导航

| 专题目录 | 核心原语 | 关键掌握目标 |
| :--- | :--- | :--- |
| **[01-mysql-deep](./01-mysql-deep/)** | InnoDB 锁模型与死锁复盘 | 掌握 MVCC、Next-Key Lock 加锁区间判定、死锁日志分析与索引优化 |
| **[02-redis-patterns](./02-redis-patterns/)** | 分布式锁与高并发防击穿 | 掌握 Lua 原子锁+看门狗续期、SingleFlight 聚合与缓存双写一致性 |
| **[03-kafka-kraft](./03-kafka-kraft/)** | KRaft 架构与消息防丢失防重 | 深入 KRaft 架构模型、端到端消息零丢失配置与消费端幂等防重实战 |

---

## 🚀 统一运行验证

```bash
# 阶段全量单测与验证程序串行执行
go test -v ./04-storage-middleware/...
```
