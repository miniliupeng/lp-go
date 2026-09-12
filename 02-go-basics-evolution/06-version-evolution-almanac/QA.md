# 专题：Go 版本演进与废弃大厂高频面试题

---

### Q1: 为什么 Go 1.16 废弃了 `io/ioutil` 包？旧函数对应迁移到哪里？

**标准回答**：
- **设计缺陷**：`io/ioutil` 历史上是一个“杂物箱（Catch-all）”包，内部函数命名和职责不清晰（例如 `ioutil.ReadFile` 本质是对 OS 文件系统的操作，放在 `io` 抽象层违背了单一职责原则；而 `ioutil.ReadAll` 又是纯 I/O 流操作）。
- **迁移对应方案**：
  1. `ioutil.ReadFile` $\to$ **`os.ReadFile`**
  2. `ioutil.WriteFile` $\to$ **`os.WriteFile`**
  3. `ioutil.ReadAll` $\to$ **`io.ReadAll`**
  4. `ioutil.Discard` $\to$ **`io.Discard`**
  5. `ioutil.NopCloser` $\to$ **`io.NopCloser`**
  6. `ioutil.TempDir / TempFile` $\to$ **`os.MkdirTemp` / `os.CreateTemp`**

---

### Q2: 经典面试题：Go 1.22 对 `for range` 循环变量的作用域做了什么改变？

**标准回答**：
- **历史问题（1.21 及更早）**：`for k, v := range slice` 中，变量 `v` 在整场循环中**只分配一次地址**，每轮循环只是将新值覆盖写入该地址。如果在循环体内启动 Goroutine 并引用了 `v`，由于并发执行延迟，协程最终读到的往往是最后一次循环的值（即著名的闭包抓取同一引用 Bug）。旧解决办法是必须在循环体内手动写一行 `v := v` 进行局部变量遮蔽（Shadowing）。
- **1.22 演进**：Go 1.22 彻底修改了语言规范，`for range` 的每次迭代都会**创建全新且生命周期独立的变量实例**。从此不再需要 `v := v`，直接引用变量即并发安全，消除了过去十年 Go 最具争议的语言陷阱。

---

### Q3: Go 1.13 引入的 `errors.Is` 和 `errors.As` 解决了什么痛点？

**标准回答**：
- **旧痛点**：过去错误被包装后（如 `fmt.Errorf("do something: %v", err)`），由于字符串被拼接改变，外部无法再通过 `err == io.EOF` 判定底层具体错误类型，开发者只能无奈地进行危险的字符串包含判断（`strings.Contains(err.Error(), "EOF")`），易错且脆弱。
- **现代化解法**：
  1. 通过 `%w` 动词包装错误（`fmt.Errorf("failed: %w", err)`），形成错误嵌套链（Unwrap 树）。
  2. **`errors.Is(err, target)`**：递归调用 `Unwrap()` 判定链路中是否包含目标哨兵错误（Sentinel Error）。
  3. **`errors.As(err, &targetStruct)`**：递归查找链路上是否有匹配目标指针类型的自定义结构体错误，并将其还原赋值，便于提取结构体内部字段。

---

### Q4: Go 的泛型（Generics）在底层是如何实现的？与 Java/C++ 有何区别？

**标准回答**：
- **C++ 模式（Monomorphization 全单态化）**：为每个具体的类型实例化一套机器码，执行最快但编译慢、生成二进制文件体积巨大。
- **Java 模式（Type Erasure 类型擦除）**：将类型全部擦除为 Object 并强转，由于涉及大量基本类型包装与反射装箱，运行时性能较差。
- **Go 1.18+ 模式（Gcshape Stenciling + Dictionary Passing）**：
  - **指针与引用类型共享代码（gcshape）**：所有具有相同内存布局的类型（如所有的指针类型、interface）共享同一份生成的机器指令，通过运行时传递“类型字典（Dictionary）”来动态获取具体方法和大小。
  - **值类型独立生成**：对于 int、float64 等具有不同大小和寄存器调用约定的值类型，分别生成特化机器码。
  - **平衡哲学**：在编译速度、二进制文件体积与运行时高性能之间取得了极佳平衡。
