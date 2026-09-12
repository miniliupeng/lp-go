# 专题 03：Kafka KRaft 架构、零丢失与高性能消费大厂高频面试题

---

### Q1: Kafka 是如何利用操作系统特性做到每秒处理数百万条消息的高性能吞吐的？

**标准回答**：
Kafka 的超高吞吐依赖于四大底层系统级黑科技：
1. **顺序写磁盘（Sequential I/O）**：Kafka 将消息以追加写（Append-only）的方式写入 Partition 的 Segment 日志文件。现代操作系统与磁盘硬件对顺序写进行了深度预读与调度优化，机械硬盘的顺序写性能可达数百 MB/s，甚至逼近内存随机写的速度；
2. **页缓存（PageCache）**：Kafka 几乎不依赖 JVM 堆内存来缓存消息，而是直接利用 Linux 操作系统的 PageCache。数据写入和读取都发生在内核 PageCache 中，避免了 JVM GC 的停顿与对象开销；
3. **零拷贝技术（Zero-Copy `sendfile`）**：
   - 传统文件传输需要经历 4 次上下文切换与 4 次数据拷贝（磁盘 -> 内核缓冲区 -> 用户态内存 -> Socket 缓冲区 -> 网卡）；
   - Kafka 在向 Consumer 发送消息时，直接调用 Linux 系统的 `sendfile()` 系统调用，**数据直接从 PageCache 经由 DMA 控制器传输到网卡（NIC）缓冲区，全程 0 次 CPU 数据复制，无需经过用户态**！
4. **批量发送与压缩（Batching & Compression）**：Producer 端通过 `batch.size` 与 `linger.ms` 积累一批数据，使用 `zstd` / `snappy` 整体压缩后单次网络传输，极大减少网络 I/O 次数。

---

### Q2: 为什么仅在生产者设置 `acks=all` 依然不能 100% 保证消息不丢失？

**标准回答**：
- **致命陷阱**：`acks=all`（或 `-1`）的定义是“等待当前所有处于 **ISR（In-Sync Replicas）集合中的副本** 全部写入确认”。
- **极端丢失场景**：
  - 假设某个 Topic 的副本因子是 3（1 个 Leader + 2 个 Follower）；
  - 后来 2 个 Follower 由于网络抖动或负载过高，被从 ISR 列表中踢出，此时 **ISR 集合中只剩下了 Leader 自己**！
  - 此时若发送端带着 `acks=all` 发送消息，只要 Leader 自身写入成功就会立即返回 ACK！
  - 紧接着该 Leader 节点发生硬件故障断电，新副本启动由于没有同步到刚刚那条数据，**该消息永久丢失！**
- **终极防御**：必须配合服务端参数 **`min.insync.replicas = 2`**。它强制要求：当处于 ISR 中的副本数小于 2 时，Broker 直接拒绝 Producer 的写入请求并抛出异常，配合 `unclean.leader.election.enable=false` 彻底杜绝丢失。

---

### Q3: Kafka 消费者的 Rebalance（重平衡）机制有什么危害？如何有效避免？

**标准回答**：
- **重平衡的危害**：
  - **全员停顿（Stop The World）**：在重平衡发生期间，整个 Consumer Group 的所有消费者会暂停消费，业务延迟急剧攀升；
  - **重复消费风险**：如果在重平衡前由于超时未提交 Offset，重平衡后其他消费者从旧 Offset 重新拉取，会触发大面积重复消费。
- **引发 Rebalance 的常见场景与避免方案**：
  1. **消费处理太慢导致的心跳超时**：
     - 单批次拉取的消息过多，本地处理时间超过了 `max.poll.interval.ms`（默认 5 分钟），Coordinator 误认为该消费者已死；
     - **解法**：适当调小 `max.poll.records`（如单次仅拉取 50 条），或将耗时繁重的业务处理放入后台 Goroutine WorkerPool 中异步处理。
  2. **偶发网络抖动导致的心跳丢失**：
     - 心跳线程未能在 `session.timeout.ms`（如 45 秒）内向 Broker 发送心跳；
     - **解法**：合理调大 `session.timeout.ms`，并将心跳发送间隔 `heartbeat.interval.ms` 设置为前者的 1/3。
  3. **开启静态成员机制（Static Membership）**：
     - 配置 `group.instance.id`，当容器 Pod 发生发布滚动重启时，在设定时间内不会触发重平衡，避免无谓的集群震荡。

---

### Q4: 现代 KRaft 模式相比于传统 ZooKeeper 模式，如何防范“脑裂（Split-Brain）”？

**标准回答**：
- 传统 ZooKeeper 依赖临时节点和心跳维护 Leader，在偶发网络分区时，旧 Controller 可能未能及时感知自己被踢出，而新 Controller 已经产生，可能造成短时间的指令冲突。
- **KRaft 的防脑裂机制**：
  - 严格依托 **Raft 强一致性共识算法**；
  - 每个控制周期必须拥有单调递增的 **`Leader Epoch`（任期号）**；
  - 任何元数据写入和配置变更必须得到 **Quorum（多数派：超过半数节点）** 的确认才算生效；
  - 如果一个孤立的网络分区节点自认为是 Leader，由于无法联络到超过半数的节点，它的任何写操作都会被强行拒绝，从数学和协议层面彻底杜绝了脑裂。
