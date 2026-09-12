# 阶段二：Go 进阶穿透与版本演进

> 本阶段建立在已掌握基础语法的基石之上，聚焦**语言机制的底层穿透与现代版本演进**：深入物理内存对齐、defer 汇编时序、Slice/Map 内部模型与零拷贝、接口 eface/iface 内存结构，以及 Go 1.0~1.27 现代特性编年史。

---

## 🗺️ 阶段二 6 大进阶穿透专题全景

```text
02-go-basics-evolution/
├── README.md                      # 阶段二总纲（当前文件）
│
├── 01-syntax-and-memory/          # 专题 1：内存对齐规则与 unsafe 指针体系
├── 02-control-and-defer/          # 专题 2：Go 1.22 循环修正 ──穿透──> defer 汇编时序与错误链
├── 03-composite-data-deep/        # 专题 3：SliceHeader 内部模型、哈希扩容机制与零拷贝
├── 04-oop-interface-reflect/      # 专题 4：鸭子类型方法集隐式约束、eface/iface 与 typed-nil 陷阱
├── 05-concurrency-basics/         # 专题 5：Channel 3 态 5 操状态矩阵、死锁边界与 Context 级联
└── 06-version-evolution-almanac/  # 专题 6：Go 1.0 ~ 1.27 里程碑 ──穿透──> b.Loop()/clear() 汇编实测
```

---

## 📚 专题对照与底层知识覆盖清单

| 专题目录 | 进阶应用层 | 硬核底层穿透层 |
| :--- | :--- | :--- |
| 📁 [01-syntax-and-memory](./01-syntax-and-memory/) | make vs new、类型别名 vs 类型定义、5 种字符串拼接 | **内存对齐规则、Padding 计算、字段重排省 33% 内存、空结构体 0 内存原理、uintptr vs unsafe.Pointer GC 风险** |
| 📁 [02-control-and-defer](./02-control-and-defer/) | Go 1.22 循环变量独立、现代 errors 链 | **defer 的 RET 汇编时序（三步机理）、命名返回值修改陷阱、panic 跨协程失效、开放编码 2.1ns 开销** |
| 📁 [03-composite-data-deep](./03-composite-data-deep/) | 三索引切片防污染、Map 并发崩溃分析 | **SliceHeader 内部模型、扩容平滑演进源码、Map 哈希桶扩容与并发写 fatal error、零拷贝字符串转换 (2.07ns)** |
| 📁 [04-oop-interface-reflect](./04-oop-interface-reflect/) | 值接收者 vs 指针接收者方法决策、类型断言 | **鸭子类型方法集判定规则、`eface` / `iface` 底层差异、`typed-nil` 接口不等于 nil 陷阱、反射性能代价与 CanSet** |
| 📁 [05-concurrency-basics](./05-concurrency-basics/) | sync.WaitGroup 协调、select 伪随机防饥饿 | **Channel 3 态 5 操完整矩阵、死锁产生边界条件、`select` 伪随机轮询机制、`context.Context` 树形级联** |
| 📁 [06-version-evolution-almanac](./06-version-evolution-almanac/) | Go 1.0 ~ 1.27 里程碑历史、io/ioutil 废弃清单 | **`b.Loop()` 现代测试规范、内置 `clear()` 比传统清零快 190 倍基准实测、泛型 Gcshape 底层混合模型** |
