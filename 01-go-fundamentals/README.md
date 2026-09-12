# 阶段一：Go 基础语法全貌与核心原语筑基

> 本阶段专为建立清晰、地道的 Go 语言基础心智打造。系统化覆盖从数据类型、变量常量、运算表达式、流程控制语句，到函数签名、错误处理哲学与指针/切片/Map 容器的基础日常 CRUD 操作。

---

## 🗺️ 阶段一 7 大基础筑基专题全景

```text
01-go-fundamentals/
├── README.md                          # 阶段一总纲与知识图谱（当前文件）
│
├── 01-data-types-and-declarations/    # 专题 1：数据类型全貌（数值/文本/布尔/字符）与变量常量 iota
│   ├── README.md                      # 数据类型速查、变量声明四部曲、零值体系、常量与 iota
│   ├── main.go                        # 全类型声明、零值回显与 iota 枚举实操
│   └── main_test.go                   # 类型与常量赋值单元测试
│
├── 02-operators-and-expressions/      # 专题 2：运算符全景与独有位清空 (&^) 操作符
│   ├── README.md                      # 算术/关系/逻辑/位运算、位清空 &^、自增自减独立语句铁律
│   ├── main.go                        # 运算规则与 RBAC 权限位清空演示
│   └── main_test.go                   # 运算符计算逻辑单元测试
│
├── 03-control-flow-statements/        # 专题 3：流程控制（条件分支、自动分号 ASI 与四态循环）
│   ├── README.md                      # if 前置短声明、else if 换行禁忌、for 循环四种形态、switch 自动 break
│   ├── main.go                        # 控制流完整实战代码
│   └── main_test.go                   # 流程控制分支与循环逻辑测试
│
├── 04-functions-and-signatures/       # 专题 4：函数签名、多返回值与闭包入门
│   ├── README.md                      # 多返回值、命名返回值、可变参数 ...T、闭包、全盘值拷贝本质
│   ├── main.go                        # 算术多返回值与状态累加闭包演示
│   └── main_test.go                   # 函数调用与闭包状态单测
│
├── 05-error-handling-basics/          # 专题 5：错误处理哲学与 defer 资源释放入门
│   ├── README.md                      # 显式错误返回契约、为什么没有 try-catch、defer 成对释放、panic/recover 兜底
│   ├── main.go                        # 自定义错误返回、defer LIFO 时序、安全 recover 演示
│   └── main_test.go                   # 错误断言与异常捕获单测
│
├── 06-pointers-and-slices-maps/       # 专题 6：指针入门、数组长度推导、切片与字典基础 CRUD
│   ├── README.md                      # & 取地址与 * 解引用、数组推导 [...]T、切片四大操作、Map 随机遍历哈希种子
│   ├── main.go                        # 指针改值、切片扩容与 Map 增删改查实操
│   └── main_test.go                   # 指针与容器基础操作单测
│
└── 07-structs-and-interfaces/         # 专题 7：结构体定义、方法接收者与非侵入式接口多态
    ├── README.md                      # 结构体声明、值接收者 vs 指针接收者深度对比、鸭子类型接口、多态传参实战
    ├── main.go                        # 车辆引擎统一契约与多态调用实操
    ├── main_test.go                   # 接收者方法与接口多态断言单测
    └── QA.md                          # 为什么没有 class、接收者如何选等面试高频真题
```

---

## 🎯 专题速查清单与认知目标

| 专题目录 | 核心掌握要点 | 关键避坑指南 |
| :--- | :--- | :--- |
| **01-data-types** | 掌握所有整型、浮点、字符串（双引号 vs 反引号）、字符 `byte`/`rune` | 严禁隐式类型转换，哪怕同尺寸也必须 `int64(a)` 显式强转 |
| **02-operators** | 掌握算术/位运算与独有的 `&^` 位清空 | 严禁写 `++i` 或 `x = i++`，自增是独立语句不是表达式 |
| **03-control-flow** | 熟练编写 `if` 前置声明、自动分号机制、现代 `for i := range 10` | 严禁将 `else` 换行写在新的一行，否则触发词法分析器 ASI 编译报错 |
| **04-functions** | 掌握多返回值、命名返回值与高阶闭包 | 牢记 Go 语言一切参数传递皆为**值拷贝（Pass-by-value）** |
| **05-error-handling**| 掌握 `if err != nil` 检查与 `defer` 资源清理 | `defer` 执行顺序是**后进先出（LIFO，类似压栈）** |
| **06-pointers-and-slices-maps** | 搞懂指针操作、数组 `[...]` 推导、切片 `append` 与 Map `delete` | `for range` 遍历 Map 顺序随机是运行时刻意注入的防脆弱设计 |
| **07-structs-and-interfaces** | 掌握结构体初始化、值/指针接收者差异、鸭子类型接口与多态 | 结构体包含锁原语时**严禁使用值接收者**，大厂 95% 场景首选指针接收者 |

---

## ⚡ 极速全阶段自动化测试

```bash
go test -v ./01-go-fundamentals/...
```

