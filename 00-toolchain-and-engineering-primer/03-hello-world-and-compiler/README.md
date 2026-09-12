# 专题 03：第一个 Go 程序与编译底层流转全解

> “Hello World” 是每个程序员的启蒙。在 Go 语言中，这段看似简单的代码背后隐藏着 Go 语言的程序入口规范、包管理结构、静态链接机制与跨平台交叉编译黑科技。

---

## 1. 经典 Hello World 逐行解剖

```go
package main // 1. 声明当前文件所属的包为 main

import (
    "fmt"    // 2. 引入标准库 fmt（Format 格式化 I/O 包）
)

// 3. 程序的唯一主入口函数
func main() {
    fmt.Println("Hello, Go World!") // 4. 向控制台标准输出打印一行文字
}
```

### 核心语法规则铁律：
1. **`package main`**：
   - 每一个可独立执行的 Go 程序，**必须且只能有一个 `main` 包**。
   - 如果声明为其他包名（如 `package user`），则表示这是一个库（Library），编译器只会将它打包成库文件，无法直接生成可执行程序。
2. **`func main()`**：
   - 这是 Go 运行时的固定入口函数名。**它不接受任何参数，也不返回任何错误码**（若要获取命令行参数，使用 `os.Args`；若要退出并返回非 0 状态码，使用 `os.Exit(code)`）。
3. **分号自动补全**：
   - Go 代码行末**不需要写分号 `;`**，Go 词法分析器会在每一行的末尾自动追加分号。这也导致了一个重要语法规定：**函数的左大括号 `{` 必须紧跟在函数声明末尾，绝对不能另起一行！**

---

## 2. 编译命令三剑客：`go run` vs `go build` vs `go install`

| 命令 | 行为机制 | 产物形态 | 适用场景 |
| :--- | :--- | :--- | :--- |
| **`go run main.go`** | 编译并在系统临时目录生成临时二进制，立即执行，随后自动清理。 | **内存/临时文件**，不污染当前目录。 | 本地日常极速调试、体验小段代码。 |
| **`go build main.go`** | 完整执行词法、语法分析、SSA 静态单赋值优化，将所有依赖静态打包为一个独立文件。 | **可执行二进制文件**（Windows 下为 `.exe`，Linux/Mac 为 ELF/Mach-O 二进制）。 | 生产环境打包、发布与交付。 |
| **`go install`** | 编译并直接将产物安装到 `$GOPATH/bin` 目录中。 | **系统全局可调用命令行工具**。 | 安装全局 CLI 工具（如 `dlv`, `protoc-gen-go`）。 |

---

## 3. Go 杀手锏特性：跨平台交叉编译（Cross-Compilation）

在其他语言（如 C/C++）中，在 Windows 上要编译一个在 Linux 上跑的程序往往需要繁琐的交叉编译工具链；而 Go 语言的编译器原生自带交叉编译能力，**只需切换两个环境变量即可一键跨平台打包**：

::: code-group
```powershell [Windows 下编译 Linux 运行包]
# 在 Windows PowerShell 中执行：
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o server_linux ./00-toolchain-and-engineering-primer/03-hello-world-and-compiler/main.go

# 编译出的 server_linux 可直接通过 scp 丢到 Linux 服务器上 chmod +x 运行，无需安装 Go 环境！
```

```bash [macOS 下编译 Windows .exe]
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o app.exe ./00-toolchain-and-engineering-primer/03-hello-world-and-compiler/main.go
```

```bash [Linux 下编译 macOS ARM64 (Apple Silicon)]
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o app_mac ./00-toolchain-and-engineering-primer/03-hello-world-and-compiler/main.go
```
:::

> [!NOTE]
> 设置 `CGO_ENABLED=0` 表示禁用 C 语言动态链接，生成**纯静态编译的单二进制文件**，能够零依赖运行在任何极简 Linux 容器（如 Scratch / Alpine）中！

---

## 4. 本地实操体验

```bash
go run ./00-toolchain-and-engineering-primer/03-hello-world-and-compiler/main.go "Go 工程师"
```
