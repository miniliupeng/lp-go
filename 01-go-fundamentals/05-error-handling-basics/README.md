# 专题 05：错误处理哲学与 defer 资源释放入门

> Go 语言最为人称道的架构哲学之一就是**“将错误视为普通的值（Errors are values）”**。Go 坚决放弃了传统语言（Java/Python/C++）中沉重、隐藏隐式控制流跳转的 `try-catch-finally` 异常机制，坚持通过显式返回值让工程师对每一处可能失败的环节保持敬畏；同时提供轻量优雅的 `defer` 机制保证资源成对安全释放。

---

## 1. Go 为什么没有 try-catch？

1. **显式优于隐式**：在 `try-catch` 体系中，你无法一眼看出一行代码到底会抛出多少种未声明的运行时异常，容易导致异常被空 catch 静默吞噬或失控冒泡；
2. **极佳的可读性与健壮性**：Go 的函数将 `error` 作为多返回值的最后一项，调用方必须立即显式处理：
   ```go
   f, err := os.Open("config.json")
   if err != nil {
       // 立即在此处决定是重试、降级还是向上层传递，代码执行路径完全线性可推导！
       return fmt.Errorf("打开配置文件失败: %w", err)
   }
   ```

---

## 2. 核心原语：`defer` 的基础法则（成对书写，防止泄漏）

`defer` 用于延迟执行某段代码，它会在**外层包裹它的函数即将退出前的那一刻**被调用。

### ① 最经典用途：成对管理资源开启与关闭
```go
func ReadFileContent(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    // 拿到资源后紧跟 defer 释放，无论后续发生什么、从哪一行 return，file 都能百分之百被安全关闭！
    defer file.Close()

    return io.ReadAll(file)
}
```

### ② 执行顺序：后进先出（LIFO，类似压栈）
如果一个函数内注册了多个 `defer`，它们的执行顺序遵循“后写的先执行”：
```go
defer fmt.Println("1")
defer fmt.Println("2")
defer fmt.Println("3")
// 函数退出时依次打印: 3 -> 2 -> 1
```

---

## 3. 极简异常防御：`panic` 与 `recover`

在 Go 语言中：
- **`error` 对应“预期内的业务错误”**（如文件不存在、网络超时、密码错误）；
- **`panic` 对应“致命的严重程序异常”**（如空指针解引用、数组越界、配置缺失导致无法启动）。

使用 `recover()` 可以在 `defer` 函数内部捕获 `panic`，防止主进程崩溃退出（大厂 Web 框架和网关的核心防护盾）：

```go
func SafeExecute() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("🚨 成功捕获并拦截了 Panic: %v，主程序继续安全运行！\n", r)
        }
    }()

    // 模拟触发严重不可控异常
    panic("致命数据库连接崩溃！")
}
```

---

## 4. 本地实操运行

```bash
go run ./01-go-fundamentals/05-error-handling-basics/main.go
```
