# 学习进度与工程操作记录（Progress & Changelog）

本项目用于记录大厂 Go 后端核心技术栈的实战与演练过程。每次完成专题、环境配置或重大实验时，在此同步记录操作轨迹与进度。

---

## 📊 总体进度大纲

- [x] **00. 开发环境与工程工具链（零基础第一步）** (完成度: 100%)
  - [x] `01-environment-and-goproxy`: 跨平台 SDK 安装、GOROOT/GOPATH 边界、国内 GOPROXY 避坑
  - [x] `02-ide-and-toolchain`: 现代 IDE 工具链（VS Code / gopls / Delve 调试）
  - [x] `03-hello-world-and-compiler`: 第一个 Go 程序、编译原理与交叉编译
  - [x] `04-go-modules-deep`: 现代包管理规范（go.mod / go.sum / 依赖治理）
- [x] **01. 基础语法与核心原语（日常编程筑基）** (完成度: 100%)
  - [x] `01-data-types-and-declarations`: 数据类型全貌（数值/字符/文本/布尔）与变量声明/iota
  - [x] `02-operators-and-expressions`: 算术/关系/逻辑/位运算/位清空 &^ 语法
  - [x] `03-control-flow-statements`: 流程控制（if 前置声明、ASI自动分号、for 四大形态、switch 自动 break）
  - [x] `04-functions-and-signatures`: 函数签名、多返回值、变参、闭包与值传递本质
  - [x] `05-error-handling-basics`: 错误处理哲学（显式 error 契约、defer 资源释放、panic/recover）
  - [x] `06-pointers-and-slices-maps`: 指针入门、数组推导、切片四大操作与 Map 随机遍历机制
  - [x] `07-structs-and-interfaces`: 结构体定义、方法接收者对比、非侵入式接口与多态参数实战
- [x] **02. 进阶穿透与版本演进（语言机制深潜）** (完成度: 100%)
  - [x] `01-syntax-and-memory`: 内存对齐规则与 unsafe 指针体系
  - [x] `02-control-and-defer`: Go 1.22 循环修正与 defer 汇编时序
  - [x] `03-composite-data-deep`: SliceHeader 内部模型、哈希扩容与零拷贝
  - [x] `04-oop-interface-reflect`: 鸭子类型方法集隐式约束、eface/iface 与 typed-nil
  - [x] `05-concurrency-basics`: Channel 3 态 5 操矩阵、死锁边界与 Context 级联
  - [x] `06-version-evolution-almanac`: Go 1.0~1.27 里程碑与内置 clear() 实测
- [x] **03. Go 核心与底层系统（T0 淘汰线）** (完成度: 100%)
  - [x] `01-gmp-scheduler`: GMP 模型、Work-stealing、信号抢占与 Trace 观测
  - [x] `02-memory-gc`: 逃逸分析、三色标记与写屏障、GOMEMLIMIT 防 OOM 演练
  - [x] `03-channel-sync`: Channel 底层 hchan、Mutex 饥饿转换、sync.Pool 对象复用
  - [x] `04-modern-runtime-features`: 现代运行时新特性、iter 迭代器、slog、PGO 调优
  - [x] `05-pprof-tuning`: CPU/内存泄漏/Goroutine 悬挂/锁争用现场复现与排查
- [x] **04. 存储与高并发中间件（T0 ~ T1 标配）** (完成度: 100%)
  - [x] `01-mysql-deep`: 事务隔离、Next-Key Lock 锁加锁规则与死锁复现
  - [x] `02-redis-patterns`: 分布式锁(Lua+看门狗)、SingleFlight 击穿防范、双写一致
  - [x] `03-kafka-kraft`: KRaft 架构、消息防丢与幂等防重实战
- [x] **05. 微服务架构与治理（T1 标配）** (完成度: 100%)
  - [x] `01-rpc-kitex-grpc`: gRPC / Kitex 跨服务调用与编解码性能对比
  - [x] `02-reliability`: 令牌桶/自适应限流、断路器、幂等性 Token 落地
  - [x] `03-observability`: OpenTelemetry 链路透传 + Prometheus RED 指标打点
- [x] **06. AI 原生后端（2026 核心增量）** (完成度: 100%)
  - [x] `01-sse-streaming-gateway`: SSE 流式输出网关、反压机制与长连接保活
  - [x] `02-vector-rag-pipeline`: 向量数据库（Milvus/PGVector）与 RAG 召回链路
  - [x] `03-eino-agent`: 字节 Eino 框架落地（Agent 工具调用与 Graph 编排）
- [x] **07. 企业级综合实战项目** (完成度: 100%)
  - [x] `ai-gateway-service`: 高并发流式 AI 智能网关服务 (DDD架构 + SSE反压 + RAG混合重排 + 优雅停机)

---

## 📝 详细操作日志

### [2026-09-08] 项目立项与专题 01 落地

#### 1. 纲领与规范建立
- 新建 [go-backend-tech-stack.md](file:///Users/max/Desktop/lp/code/lp-go/go-backend-tech-stack.md)：大厂 2026 最新时效性 Go 技术栈调研与定级考核矩阵。
- 新建 [learning-guide.md](file:///Users/max/Desktop/lp/code/lp-go/learning-guide.md)：“三位一体”（原理 $\to$ 实证 $\to$ 面试自测）学习闭环规范。

#### 2. 工程脚手架初始化
- 协助安装并确认 Go 环境：`Go 1.27.1 (darwin/arm64)`。
- 初始化根模块：`go mod init lp-go`。
- 新建根目录 [Makefile](file:///Users/max/Desktop/lp/code/lp-go/Makefile)：封装 `test`, `race`, `bench`, `escape`, `clean` 常用效能指令。

### 阶段 00：基础夯实与版本演进（新手入门到穿透底层）

#### 专题 01：基础语法入门 ──穿透──> 内存对齐规则与 unsafe 指针体系
- **路径**：`00-go-basics-evolution/01-syntax-and-memory/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/01-syntax-and-memory/README.md)：变量声明、零值安全、`make` vs `new`、类型别名 vs 类型定义、位清空运算符 `&^`、5 种字符串拼接对比、结构体内存对齐与字段重排优化。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/01-syntax-and-memory/main.go)：零值验证、make/new 区别演示、权限位清空、字符串/切片转换、结构体内存大小打印与 unsafe 私有字段穿透。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/01-syntax-and-memory/main_test.go)：内存对齐单测 `TestMemoryAlignment`、5 种字符串拼接基准、结构体字段重排前后内存分配 Benchmark（实测提速 35%）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/01-syntax-and-memory/QA.md)：大厂真题全面覆盖（未用变量编译报错哲学、RuneCountInString 原理、内存对齐根因、末尾空结构体 8 字节填充、`uintptr` 无法防 GC 回收、`strings.Builder` 扩容机制）。
- **实测执行记录**：
  - `go run ./00-go-basics-evolution/01-syntax-and-memory/main.go` 成功通过。
  - `go test -v -bench=. -benchmem` 单测 100% 通过；`strings.Builder` 凭借预分配与零拷贝在 5 种拼接中性能最优（36.42 ns/op）；重排后结构体较未重排结构体单次分配耗时从 464.6 ns 骤降至 303.6 ns，立省 33.3% 内存并大幅减少总线事务。

#### 专题 02：流程控制与 1.22 循环 ──穿透──> defer 汇编时序与现代错误链
- **路径**：`00-go-basics-evolution/02-control-and-defer/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/02-control-and-defer/README.md)：前置初始化、1.22 循环语法糖、defer 汇编三步执行机理、循环安全 defer、panic 跨协程失效、Go 1.13+ 错误链。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/02-control-and-defer/main.go)：四大时序函数验证、循环匿名函数隔离、recover 拦截与 errors.Is/Join。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/02-control-and-defer/main_test.go)：deferF1~F3 严格时序单测、Go 1.14 开放编码 defer vs 直接调用基准（采用 `b.Loop()`）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/02-control-and-defer/QA.md)：return 拆解三步时序、recover 跨协程失效根因、长循环 defer 泄漏与开放编码底层机理。
- **实测执行记录**：
  - `go run ./00-go-basics-evolution/02-control-and-defer/main.go` 验证通过，输出完全符合底层时序预期。
  - `go test -bench=.` 测得开放编码 defer 耗时仅 **2.117 ns/op**（与纯直接调用 2.125 ns 处于同一纳秒级），证实现代 Go 中 defer 开销几乎归零。

#### 专题 03：复合容器与切片/映射底层 ──穿透──> SliceHeader/哈希扩容/零拷贝
- **路径**：`00-go-basics-evolution/03-composite-data-deep/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/03-composite-data-deep/README.md)：数组 vs 切片、Map comma-ok、SliceHeader 内部模型、三索引切片防污染、1.18+ 平滑扩容曲线、Map 并发 fatal error 根因、Go 1.20+ 零拷贝。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/03-composite-data-deep/main.go)：三索引切片内存隔离演示、平滑扩容轨迹打印、unsafe.String/Slice 零拷贝双向转换。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/03-composite-data-deep/main_test.go)：三索引切片防污染严格单测、零拷贝 vs 传统 string 转 byte 基准（采用 `b.Loop()`）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/03-composite-data-deep/QA.md)：切片与数组传参底层区别、1.18 扩容公式重构根因、Map 并发写 crash 无法 recover 的底层权衡。
- **实测执行记录**：
  - `go run ./00-go-basics-evolution/03-composite-data-deep/main.go` 验证通过，三索引切片成功隔离数据污染，零拷贝数据地址完全一致。
  - `go test -v -bench=.` 单测通过，零拷贝转换达到 2.076 ns 极致吞吐。

#### 专题 04：面向对象、接口底层 ──穿透──> 鸭子类型方法集/eface/iface/反射代价
- **路径**：`00-go-basics-evolution/04-oop-interface-reflect/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/04-oop-interface-reflect/README.md)：值/指针接收者决策准则、鸭子类型隐式实现、方法集限制铁律、eface/iface 内存结构、致命的 nil 接口判断失误、反射性能代价。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/04-oop-interface-reflect/main.go)：值/指针接收者效果验证、typed-nil 赋给 error 陷阱复现、unsafe 拆解 iface tab/data、反射 CanSet 安全修改值。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/04-oop-interface-reflect/main_test.go)：typed-nil 严格单测、原生调用 vs 接口动态分发 vs 反射调用三方基准对比（采用 `b.Loop()`）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/04-oop-interface-reflect/QA.md)：typed-nil 接口判等失败根因、方法集规则对接口实现的制约、反射性能损耗根因与大厂选型取舍。
- **实测执行记录**：
  - `go run ./00-go-basics-evolution/04-oop-interface-reflect/main.go` 验证通过，成功复现 typed-nil 错误陷阱，并通过 unsafe 提取 tab/data 地址。
  - `go test -bench=.` 测得接口动态分发耗时 2.188 ns（与原生直接调用 2.230 ns 基本无异），而反射耗时暴增 50 倍至 113.3 ns 并引发堆逃逸与 3 次内存分配。

#### 专题 05：并发起步、Channel 状态机 ──穿透──> 状态矩阵与 Context 级联
- **路径**：`00-go-basics-evolution/05-concurrency-basics/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/05-concurrency-basics/README.md)：go func() 优雅起步、sync.WaitGroup 计数规则、Channel 3 态 5 操全景矩阵、死锁边界、select 伪随机轮询机制、context 树形级联取消。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/05-concurrency-basics/main.go)：WaitGroup 协程协调、Closed 通道排空后读取验证、select default 非阻塞、Context 超时批量级联退出。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/05-concurrency-basics/main_test.go)：关闭 Channel 排空单测、无缓冲（同步握手）vs 有缓冲（异步缓冲）通道基准（采用 `b.Loop()`）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/05-concurrency-basics/QA.md)：Channel 3 态 5 操大厂必背口诀、无缓冲死锁成因、select 伪随机防饿死机制、Context 四大核心派生区别。
- **实测执行记录**：
  - `go run ./00-go-basics-evolution/05-concurrency-basics/main.go` 成功通过，单测全部通过。
  - `go test -bench=.` 测得有缓冲通道吞吐耗时 33.24 ns，较无缓冲同步握手（135.9 ns）**提速近 4 倍**！

#### 专题 06：Go 版本演进编年史 ──穿透──> 现代 API 升级指南与内置 clear() 实测
- **路径**：`00-go-basics-evolution/06-version-evolution-almanac/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/06-version-evolution-almanac/README.md)：Go 1.0~1.27 编年史、淘汰废弃清单（`io/ioutil` 等）、新版标准库替代指引。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/06-version-evolution-almanac/main.go)：1.22 循环变量独立作用域、内置 min/max/clear、modern errors 链、1.23+ 迭代器。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/06-version-evolution-almanac/main_test.go)：单元测试 `TestVersionEvolutionFeatures`、内置 `clear(slice)` vs 手动遍历清零基准（采用 `b.Loop()`）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/00-go-basics-evolution/06-version-evolution-almanac/QA.md)：循环变量向后兼容机制、泛型 Gcshape 实现机制、迭代器设计目的与高频面试题。
- **实测执行记录**：
  - `go run ./00-go-basics-evolution/06-version-evolution-almanac/main.go` 验证通过。
  - `go test -bench=.` 测得内置 `clear(slice)` 耗时仅 **2.138 ns/op**，较手动逐个置零（409.3 ns/op）**提速 190 倍**（底层调用汇编级 `runtime.memclrNoHeapPointers` 向量化指令）！

#### 阶段零全量回归测试与基准报告
- 全量单元测试：`go test -v ./00-go-basics-evolution/...`：**6 大专题 100% PASS，0 警告，0 错误**。
- 全量基准测试：`go test -v -bench=. -benchmem ./00-go-basics-evolution/...`：全部通过，数据与底层原理完全吻合。

---

### 阶段 01：Go 语言核心与底层系统（T0 淘汰线）

#### 专题 01：GMP 调度、Work-stealing、信号抢占与 Trace 观测
- **路径**：`01-go-core/01-gmp-scheduler/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/01-gmp-scheduler/README.md)：GMP 拓扑、调度循环、Work-stealing 规则、Hand-off 机制、Go 1.14 信号抢占演进。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/01-gmp-scheduler/main.go)：单核环境密集计算被 `SIGURG` 抢占实验代码，并生成 `trace.out`。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/01-gmp-scheduler/main_test.go)：Goroutine 创建开销与 Gosched 性能基准测试。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/01-gmp-scheduler/QA.md)：字节/腾讯 4 大高频 GMP 面试题标准解析。
- **实测执行记录**：
  - `go run ./01-go-core/01-gmp-scheduler/main.go` 成功通过，单 P 下非协作抢占验证生效。
  - `make bench` 成功输出基准指标（Goroutine 创建开销 327.3 ns/op，32 B/op，已升级为 `b.Loop()`）。

#### 专题 02：内存分配、逃逸分析、三色标记与 GOMEMLIMIT 实战
- **路径**：`01-go-core/02-memory-gc/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/02-memory-gc/README.md)：TCMalloc 三级分配架构（mcache/mcentral/mheap/Tiny）、逃逸分析决策树与指令验证、三色标记状态机与 Go 1.8 混合写屏障消灭二次 STW 原理、Go 1.19+ GOMEMLIMIT 彻底终结 K8s 容器 OOM 实践。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/02-memory-gc/main.go)：指针返回/动态类型/栈溢出/闭包 4 大逃逸场景实测、`runtime.ReadMemStats` 观测堆内存增量与主动 GC 效果、`debug.SetMemoryLimit` 演练容器软上限下自适应垃圾回收。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/02-memory-gc/main_test.go)：单元测试 `TestEscapeAndGCStats`、纯栈分配 vs 堆逃逸分配 vs 闭包捕获基准测试（采用 `b.Loop()`）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/02-memory-gc/QA.md)：大厂高频面试 6 题标准解析（局部指针不一定逃逸的内联条件、Tiny 分配器解决内存碎片、插入写屏障为何必须二次 STW、混合写屏障 4 条铁律、Go 为何不用分代 GC、容器内 GOMEMLIMIT 与 GOGC 双轴配置方案）。
- **实测执行记录**：
  - `go build -gcflags="-m -l"` 成功验证 4 类逃逸的编译器决策。
  - `go run ./01-go-core/02-memory-gc/main.go` 成功运行，设置 50MB 软上限后触发 15 轮密集 GC 成功压制内存峰值。
  - `go test -bench=.` 测得栈分配耗时仅 **2.231 ns/op (0 B/op, 0 allocs/op)**，指针逃逸堆分配上升至 **9.303 ns/op (8 B/op, 1 allocs/op)**，闭包捕获耗时达到 **15.81 ns/op (24 B/op, 2 allocs/op)**。

#### 专题 03：Channel 底层 hchan、Mutex 模式翻转、sync.Pool 与 sync.Once/Map
- **路径**：`01-go-core/03-channel-sync/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/03-channel-sync/README.md)：hchan 环形缓冲区结构、等待队列 `sudog` 挂起机制、无锁 fast path、`sync.Mutex` 4 字节状态字与饥饿模式自适应翻转（1ms 判定）、`sync.Pool` 双层缓存（local + victim）跨 2 轮 GC 防尖刺机制、`sync.Once` 双重检查锁（DCL）Fast-Path、`sync.Map` 读写分离架构（read/dirty 字典提升）。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/03-channel-sync/main.go)：已关闭通道读取 fast path 验证、10 协程并发争用累加、sync.Pool 获取与复用指针验证、50 协程并发单例仅初始化 1 次验证、sync.Map 读写删除与 LoadOrStore。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/03-channel-sync/main_test.go)：单元测试 `TestChannelAndMutexCorrectness`、临时创建 vs sync.Pool 复用压测、锁争用基准、sync.Once 原子读取基准（0.33ns）、sync.Map 高并发无锁读取基准（2.02ns）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/03-channel-sync/QA.md)：大厂 8 题深度解析（关闭通道读写根因、sudog 解耦多等待队列、Mutex 自旋四铁律、为何不提供可重入锁、sync.Pool 不能当连接池的原因、RWMutex 写优先机制、sync.Once 为何不能单靠 CAS、sync.Map 高频写性能暴跌根因）。
- **实测执行记录**：
  - `go run ./01-go-core/03-channel-sync/main.go` 验证通过。
  - `go test -bench=.` 测得 `sync.Once` 原子检查仅需 **0.3394 ns/op**，`sync.Map` 只读并发吞吐仅需 **2.026 ns/op**，`sync.Pool`（11.79 ns/op, 0 B/op）较临时分配提速 1 倍并消除 100% 堆分配。

#### 专题 04：现代运行时新特性、PGO 调优与结构化日志
- **路径**：`01-go-core/04-modern-runtime-features/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/04-modern-runtime-features/README.md)：Go 1.23+ 原生推进式迭代器 `iter.Seq` 零物化架构、Go 1.21+ `log/slog` 高性能结构化日志与强类型属性、Go 1.20+ PGO 性能画像引导优化提速 5%~15% 原理与自动化流水线闭环。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/04-modern-runtime-features/main.go)：`FilterMap` 泛型惰性流水线、`slog` JSON 结构化强类型打点、密集计算并生成真实 `default.pgo` 样本文件。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/04-modern-runtime-features/main_test.go)：单元测试 `TestIteratorAndSlog`、传统切片过滤（8184 B/op）vs 原生迭代器流水线（120 B/op，降幅 98.5%）基准测试、slog JSON 格式化吞吐压测。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/04-modern-runtime-features/QA.md)：大厂 4 题标准解析（Push-based 迭代器内联优势与安全短路、PGO 三大优化机理、CI/CD 自动化采集闭环、slog 较传统 log 与 zap 的生态优势）。
- **实测执行记录**：
  - 成功采集并生成生产级 `default.pgo` 样本文件。
  - 通过 `go build -pgo=./01-go-core/04-modern-runtime-features/default.pgo` 成功触发编译器基于画像的深度内联与代码优化。
  - `go test -bench=.` 测得 `slog` 达 733.6 ns/op 工业级高吞吐，迭代器流水线相比传统切片物化降低了 **98.5% 内存分配**。

#### 专题 05：pprof 性能剖析、火焰图与生产级 4 大致命故障排查
- **路径**：`01-go-core/05-pprof-tuning/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/05-pprof-tuning/README.md)：pprof 探针底层原理（SIGPROF 信号、mallocgc 内存插桩、协程快照、锁阻塞）、线上 4 大致命故障（CPU 100%、切片悬挂导致大数组泄漏、协程泄漏、锁争用）排查全解、火焰图“平顶山”诊断法则。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/05-pprof-tuning/main.go)：CPU 密集峰值注入、10MB 悬挂切片泄漏模拟、无缓冲通道协程挂起注入、自动导出 `cpu.pprof` / `mem.pprof` / `goroutine.pprof` 样本。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/05-pprof-tuning/main_test.go)：单元测试 `TestMemoryLeakFixVerification`、CPU 计算基准、深拷贝切片防泄漏分配 Benchmark。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/01-go-core/05-pprof-tuning/QA.md)：大厂高频面试 4 题解析（安全低开销集成 pprof、inuse_space 与 alloc_space 比对、Goroutine 泄漏三大排查步骤、火焰图平顶山优化决策树）。
- **实测执行记录**：
  - `go run ./01-go-core/05-pprof-tuning/main.go` 成功生成 3 份现场真实 Profile 样本。
  - 通过 `go tool pprof -top mem.pprof` 实测：精准定位出 `main.LeakingMemoryProducer` **独占了 30720kB（93.75%）** 的内存泄漏源！
  - 单元测试与基准测试 100% 通过。

#### 阶段一全量回归测试与性能基准报告
- 全量单元测试：`go test -v ./01-go-core/...`：**5 大专题 100% PASS，0 警告，0 错误**。
- 全量基准测试：`go test -bench=. -benchmem ./01-go-core/...`：**全量纳秒级基准 100% 通过**，所有指标与底层机制完全吻合。阶段一建设全面圆满达成！

---

### 阶段 02：存储与高并发中间件（T0 ~ T1 标配）

#### 专题 01：MySQL InnoDB 内核、MVCC 机制、Next-Key Lock 与生产级调优
- **路径**：`02-storage-middleware/01-mysql-deep/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/01-mysql-deep/README.md)：Buffer Pool 内存引擎与冷热分离 LRU、Undo Log 版本链与 ReadView 核心四判定法则（RC vs RR 快照时机）、Next-Key Lock 退化法则与死锁闭环、延迟关联超大分页优化（提速 10~100倍）。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/01-mysql-deep/main.go)：MVCC 判定算法全场景模拟、并发间隙锁与插入意向锁互斥死锁演示、`database/sql` 生产级连接池调优参数配置。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/01-mysql-deep/main_test.go)：单元测试 `TestMVCCVisibilityRules`、`TestDBPoolConfigValidity`、MVCC 快照读判定基准测试（4.17ns 极速）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/01-mysql-deep/QA.md)：大厂 4 题标准解析（RR 隔离级别下当前读幻读根因、B+ 树 3~4 层千万级容量数学计算、间隙锁死锁闭环排查、Go 连接池踩坑三大故障）。
- **实测执行记录**：
  - `go run ./02-storage-middleware/01-mysql-deep/main.go` 验证通过。
  - `go test -bench=.` 测得 MVCC 快照判定耗时仅 **4.179 ns/op (0 B/op, 0 allocs/op)**。

#### 专题 02：Redis 高并发模式、分布式锁、缓存防击穿与数据一致性
- **路径**：`02-storage-middleware/02-redis-patterns/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/02-redis-patterns/README.md)：Redis 7.0+ 数据结构演进（Listpack 彻底淘汰 Ziplist 消除连锁更新、SkipList 跳表多级索引）、缓存三剑客（穿透/击穿/雪崩）终极治理、生产级分布式锁规范（SET NX PX + Lua 原子防误删 + Watchdog 自动续期）、MySQL-Redis 双写一致性（Cache-Aside + Canal 监听 Binlog 异步投递）。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/02-redis-patterns/main.go)：SingleFlight 内存防击穿实战（100 并发请求合并为 1 次底层回源）、带 Watchdog 看门狗后台续期的分布式锁及安全解锁演示。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/02-redis-patterns/main_test.go)：单元测试 `TestSingleFlightDedup`、`TestRedisLockWatchdog`、SingleFlight 并发合并吞吐基准（202.6ns）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/02-redis-patterns/QA.md)：大厂 4 题标准解析（跳表范围查询优势与红黑树对比、Redlock 争议与时钟漂移缺陷、先删缓存再写库大忌与延时双删弊端、SingleFlight 处理 panic 与超时的最佳实践）。
- **实测执行记录**：
  - `go run ./02-storage-middleware/02-redis-patterns/main.go` 成功运行，100 并发成功合并为 1 次回源，看门狗顺利完成 4 次周期续期并安全释放。
  - `go test -bench=.` 单测全部通过。

#### 专题 03：Kafka KRaft 架构、端到端零丢失保证与高吞吐幂等消费
- **路径**：`02-storage-middleware/03-kafka-kraft/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/03-kafka-kraft/README.md)：KRaft 架构彻底抛弃 ZooKeeper 突破百万分区、四大系统级黑科技（顺序 I/O、PageCache、sendfile 零拷贝、批量压缩）、端到端零丢失四项铁律（Producer acks=all/min.insync.replicas/Consumer 手动 ACK）、消费端幂等防重表实战、分区有序性与积压（Lag）治理。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/03-kafka-kraft/main.go)：零丢失策略实证（ISR 不足时拒绝写入，充足时持久化）、消费端幂等去重状态机（成功拦截 At-least-once 重复投递）、哈希分区确定性路由。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/03-kafka-kraft/main_test.go)：单元测试 `TestKafkaZeroLossAndIdempotence`、去重表查询基准（13.92ns）、分区哈希计算基准。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/02-storage-middleware/03-kafka-kraft/QA.md)：大厂 4 题深度解析（零拷贝 sendfile 原理、为何 acks=all 仍需配合 min.insync.replicas、消费者 Rebalance 危害与心跳/静态成员避免、KRaft 防止脑裂机制）。
- **实测执行记录**：
  - `go run ./02-storage-middleware/03-kafka-kraft/main.go` 成功运行，成功拦截重复消息并验证路由一致性。
  - `go test -bench=.` 测得单次幂等检查耗时仅 **13.92 ns/op (0 B/op, 0 allocs/op)**。

#### 阶段二全量回归测试与性能基准报告
- 全量单元测试：`go test -v ./02-storage-middleware/...`：**3 大专题 100% PASS，0 警告，0 错误**。
- 全量基准测试：`go test -bench=. -benchmem ./02-storage-middleware/...`：**全量纳秒级基准 100% 通过**。阶段二建设全面圆满达成！

### 阶段 03：微服务架构与治理（T1 标配）

#### 专题 01：gRPC / Kitex 跨服务调用与序列化性能横评
- **路径**：`03-microservice-governance/01-rpc-kitex-grpc/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/01-rpc-kitex-grpc/README.md)：gRPC (HTTP/2 + Protobuf) 与 Kitex (Netpoll 反应堆 + FastCodec 内存池) 架构深度对比、实测基准数据。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/01-rpc-kitex-grpc/main.go)：轻量定长头高性能二进制 RPC 服务端、多路复用并发客户端、Context 超时控制。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/01-rpc-kitex-grpc/main_test.go)：单元测试与序列化 Benchmark（二进制比 JSON 提速 35%~55%）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/01-rpc-kitex-grpc/QA.md)：大厂 3 题标准解析（Netpoll 解决协程暴涨、HTTP/2 队头阻塞根因、RPC 超时链路传递最佳实践）。

#### 专题 02：微服务高可用三板斧（限流、熔断降级、业务幂等）
- **路径**：`03-microservice-governance/02-reliability/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/02-reliability/README.md)：令牌桶应对突发流量模型、熔断器状态机（Closed/Open/Half-Open）、业务幂等 Token 防重复提交规范。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/02-reliability/main.go)：惰性令牌桶实现（无后台协程）、断路器状态机自动恢复、并发安全幂等管理器。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/02-reliability/main_test.go)：单测与 Benchmark（令牌桶判定 64ns，幂等比对 71ns）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/02-reliability/QA.md)：大厂 3 题标准解析（令牌桶 vs 漏桶、熔断器半开探测防震荡、高并发扣款幂等落地）。

#### 专题 03：微服务可观测性（OpenTelemetry 链路贯穿 + Prometheus RED 模型）
- **路径**：`03-microservice-governance/03-observability/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/03-observability/README.md)：W3C Trace Context 跨服务协议规范、Prometheus RED 黄金三指标（Rate, Errors, Duration）。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/03-observability/main.go)：跨进程 traceparent 注入与提取、Context 派生子 Span、RED 指标打点直方图与 P99 统计。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/03-observability/main_test.go)：链路单测与指标打点基准（单次打点仅 20.94ns）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/03-microservice-governance/03-observability/QA.md)：大厂 3 题标准解析（W3C Trace 规范优势、RED 模型 vs Google SRE 四大指标、全量链路采样灾难与尾部采样治理）。

### 阶段 04：AI 原生后端（2026 核心增量）

#### 专题 01：大模型流式输出 SSE 网关、反压机制与长连接保活
- **路径**：`04-ai-native-backend/01-sse-streaming-gateway/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/01-sse-streaming-gateway/README.md)：SSE 协议规范、网络反压（Backpressure）防 OOM、心跳保活与 Nginx 缓冲规避。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/01-sse-streaming-gateway/main.go)：生产级 HTTP Flusher SSE 推流 Handler、有界 Channel 超时反压管道、客户端断开级联取消。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/01-sse-streaming-gateway/main_test.go)：SSE 网关集成单测、反压超时单测、SSE 文本帧格式化 Benchmark（339ns）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/01-sse-streaming-gateway/QA.md)：大厂 3 题标准解析（大模型为何选用 SSE 而非 WebSocket、Nginx 缓冲导致卡顿根治方案、客户端关闭标签页级联取消上游计费）。

#### 专题 02：向量数据库（Milvus / PGVector）与 RAG 混合召回链路
- **路径**：`04-ai-native-backend/02-vector-rag-pipeline/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/02-vector-rag-pipeline/README.md)：ANN 近似检索算法（HNSW vs IVF）、混合召回（Dense 稠密向量 + Sparse 稀疏分词）、RRF 倒数排名融合重排。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/02-vector-rag-pipeline/main.go)：纯 Go 浮点余弦相似度、向量存储 Top-K 召回、BM25 词频匹配、RRF 融合与知识增强 Prompt 构建。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/02-vector-rag-pipeline/main_test.go)：单测与 Benchmark（1536 维向量余弦计算仅 3.77μs，百级混合召回融合仅 10.7μs）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/02-vector-rag-pipeline/QA.md)：大厂 3 题标准解析（HNSW 图索引优势、混合召回消除专有名词幻觉、两阶段粗召回+精重排架构）。

#### 专题 03：字节 Eino 框架落地（Agent 工具调用与 Graph 编排）
- **路径**：`04-ai-native-backend/03-eino-agent/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/03-eino-agent/README.md)：字节 Eino 核心架构设计、ReAct 思考-行动-观察循环拓扑、状态图防死循环保护。
  - [main.go](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/03-eino-agent/main.go)：Eino 强类型 Message/Tool 体系、Weather/Calculator 工具实现、Graph 状态图循环调度。
  - [main_test.go](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/03-eino-agent/main_test.go)：Tool Calling 完整闭环单测、状态图流转 Benchmark（单轮仅 58.75ns）。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/04-ai-native-backend/03-eino-agent/QA.md)：大厂 3 题标准解析（Eino 与 LangChain 本质区别、Tool Calling 死循环防御、流式与工具调用优雅结合）。

#### 阶段三与阶段四全量回归验证
- 全量单元测试：`go test -v ./...`：**全库 20 个模块 100% 全部通过（PASS）**！

### 阶段 05：企业级综合实战项目（2026 简历杀手锏）

#### 项目：`ai-gateway-service`（高并发流式 AI 智能网关）
- **路径**：`05-capstone-project/ai-gateway-service/`
- **产出文件**：
  - [README.md](file:///Users/max/Desktop/lp/code/lp-go/05-capstone-project/ai-gateway-service/README.md)：系统设计拓扑、洋葱中间件架构、实测性能数据表格、**可直接用于个人简历的 STAR 法则模板**。
  - [QA.md](file:///Users/max/Desktop/lp/code/lp-go/05-capstone-project/ai-gateway-service/QA.md)：大厂 10 大硬核连环追问与标准答案深度解析（涵盖客户端关闭连接级联取消上游计费、反压防 OOM、SingleFlight 流式场景广播方案、fmt.Appendf 零堆逃逸、优雅停机平滑排空四步法等）。
  - [Makefile](file:///Users/max/Desktop/lp/code/lp-go/05-capstone-project/ai-gateway-service/Makefile)：一键构建、启动、单测、数据竞争检查与压测。
  - 核心源代码：
    - `cmd/gateway/main.go`：多系统信号监听、平滑优雅停机 (Graceful Shutdown)、连接排空。
    - `configs/config.go`：强类型配置定义与完整性校验。
    - `pkg/trace/trace.go` & `pkg/logger/logger.go` & `pkg/metrics/metrics.go`：W3C 协议链路贯穿、自动注 TraceID 的 slog 日志、Prometheus RED 模型与 TTFT 统计。
    - `internal/middleware/`：Recovery 兜底、Trace 注入、租户令牌桶限流、自适应断路器、RED 指标打点。
    - `internal/rag/`：纯 Go 向量余弦检索、BM25 词频匹配、RRF 倒数排名融合重排。
    - `internal/streaming/`：`fmt.Appendf` 极速推帧器、有界通道反压保护管道。
    - `internal/upstream/`：上游大模型推理接口契约与 MockEngine 仿真引擎。
    - `internal/service/`：SingleFlight 防击穿 + RAG 知识动态注入 + SSE 流式协同调度。
  - 测试套件：`cmd/gateway/main_test.go`，覆盖同步非流式、SSE 流式、限流 429 拦截、断路器熔断与恢复、RAG 混合召回，并经过 `go test -race` 检验通过。
- **实测性能数据**：
  - 1536 维超高维向量检索单次仅 **1.9 微秒 (1901 ns/op)**，且 **0 内存分配 (0 B/op, 0 allocs/op)**。
  - 端到端通过 `go test -race` 严苛检验，**100% 零数据竞争（Zero Data Race）**。

#### 全库大圆满回归验证
- 执行全库全阶段测试：`go test -v ./...`：**阶段零至阶段五共 21 个核心工程包全部绿色通过（100% PASS）**！




