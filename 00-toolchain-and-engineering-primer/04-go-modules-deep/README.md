# 专题 04：Go Modules 现代包管理规范与工程依赖治理

> 在 Go 1.11 之前，Go 依赖全局环境变量 `GOPATH` 管理代码，所有的项目必须挤在 `$GOPATH/src` 目录下，且无法管理同一个库的不同版本。自 Go 1.11 引入、并在现代 Go 中确立为唯一标准的 **Go Modules** 彻底终结了历史混乱。本专题带你掌握现代 Go 项目工程组织与依赖治理的核心规范。

---

## 1. 核心文件体系：`go.mod` 与 `go.sum` 的本质

一个标准的 Go Modules 项目由两个核心文件锚定：

### ① `go.mod`：项目依赖清单（类似前端的 package.json 或 Java 的 pom.xml）
```text
module my-project // 1. 声明本项目的根模块路径

go 1.27.1         // 2. 声明最低要求的 Go 编译器版本

require (         // 3. 显式引入的直接依赖
    github.com/google/uuid v1.6.0
)

require (         // 4. 间接依赖（带有 // indirect 标记）
    golang.org/x/sys v0.20.0 // indirect
)
```

### ② `go.sum`：密码学加密校验哈希表（类似 package-lock.json）
- 记录了当前项目所依赖的每一个包的特定版本的**密码学 SHA-256 校验和**。
- **核心价值**：防止供应链攻击！确保在任何新机器或 CI/CD 流水线上重新拉取的第三方包，其内容与当初开发者锁定的代码 100% 一字不差。
- **规范**：`go.sum` 必须提交到 Git 版本库，严禁加入 `.gitignore`！

---

## 2. 常用依赖治理命令速查表

| 命令 | 行为与用途 | 大厂推荐时机 |
| :--- | :--- | :--- |
| **`go mod init <模块名>`** | 初始化当前目录为 Go 模块，生成 `go.mod` 文件。 | 项目立项第一天执行。 |
| **`go mod tidy`** | **最常用！** 自动扫描所有代码：下载缺失依赖、移除不再引用的僵尸依赖，并更新 `go.sum`。 | 新写了 `import` 之后、提交 Git 之前必跑。 |
| **`go get <包路径>@<版本>`** | 下载或升级指定第三方依赖（如 `go get -u` 升级次版本）。 | 需要安装或更新特定组件时。 |
| **`go mod download`** | 预先下载 `go.mod` 中定义的所有依赖到本地缓存。 | Docker 容器构建时用于分层缓存。 |
| **`go mod vendor`** | 将所有第三方源码直接拷贝到项目内部的 `vendor/` 目录下。 | 离线网络环境交付或合规安全审查。 |

---

## 3. 大厂工程结构推荐：多包模块化拆分

在 Go 语言中，文件夹名通常就是包名（Package Name）。同一个目录下所有 `.go` 文件的第一行必须是相同的包声明。

```text
my-service/
├── go.mod                # 根模块定义 (module my-service)
├── go.sum                # 依赖锁版本
├── cmd/                  # 包含 main.go 的程序入口目录
│   └── app/
│       └── main.go
└── internal/             # 仅限当前模块内部使用的私有包（Go 编译器强制保护）
    ├── service/          # 业务逻辑包
    └── repository/       # 数据访问包
```

> [!IMPORTANT]
> **可见性大写铁律**：在 Go 语言中，标识符（函数名、结构体名、变量名、字段名）的**首字母大写表示公开（Public/Exported）**，可被其他包引用；**首字母小写表示私有（Private）**，仅在当前包内部可见！没有 `public/private` 关键字。

---

## 4. 本地实操体验

```bash
go run ./00-toolchain-and-engineering-primer/04-go-modules-deep/main.go
```
