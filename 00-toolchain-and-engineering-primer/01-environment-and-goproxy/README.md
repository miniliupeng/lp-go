# 专题 01：Go SDK 安装、环境变量体系与国内代理避坑指南

> 搭建开发环境是成为一名合格工程师的第 0 步。本专题帮助新手彻底扫清 SDK 安装、环境变量、以及由于网络原因导致的依赖拉取失败问题。

---

## 1. 跨平台 Go SDK 安装指南

请选择你的宿主操作系统，执行官方推荐的安装流程：

::: code-group
```powershell [Windows]
# 推荐使用官方 MSI 安装包或 winget 包管理器安装
winget install GoLang.Go

# 安装完成后重新打开终端，验证安装状态
go version
```

```bash [macOS]
# 推荐使用 Homebrew 安装最新稳定版本
brew install go

# 安装完成后验证
go version
```

```bash [Linux (Ubuntu/Debian)]
# 1. 下载最新二进制归档（以 1.27.1 为例）
wget https://golang.google.cn/dl/go1.27.1.linux-amd64.tar.gz

# 2. 解压至 /usr/local
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz

# 3. 追加 PATH 环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 4. 验证
go version
```
:::

---

## 2. 核心环境变量解密：GOROOT vs GOPATH

很多新手经常分不清 `GOROOT` 和 `GOPATH`，甚至网上很多陈旧教程误导大家去硬配全局 GOPATH，导致项目报错。必须分清两者的本质边界：

| 环境变量 | 含义与定位 | 是否需要手动配置？ |
| :--- | :--- | :--- |
| **`GOROOT`** | **Go SDK 的安装根目录**。存放 Go 标准库源码（如 `fmt`、`net/http`）以及编译器二进制文件（`go`、`gofmt` 等）。 | **通常无需手动配置**。安装程序会自动识别。 |
| **`GOPATH`** | **Go 历史工作区目录**。在现代 Go Modules 时代，它仅作为**第三方依赖包的本地全局只读缓存目录**（`$GOPATH/pkg/mod`）和编译可执行工具目录（`$GOPATH/bin`）。 | **保持默认即可**（Windows 默认在 `%USERPROFILE%\go`，macOS/Linux 默认在 `~/go`）。**切勿把自己的代码项目强行扔到 GOPATH/src 里！** |
| **`GOPROXY`** | **依赖代理服务器地址**。指定 `go get` 拉取第三方开源模块时的中继代理。 | **国内开发强烈建议手动配置！** |
| **`GO111MODULE`**| **模块化开关**。现代 Go（Go 1.16+ 至今）默认均为 `on`。 | 保持默认 `on`。 |

---

## 3. 国内开发生命线：配置 GOPROXY 代理加速

由于默认的官方模块代理 `proxy.golang.org` 在国内网络环境下常出现超时或阻断，必须配置国内合规的高速镜像源。

执行以下两条跨平台通用命令：

```bash
# 1. 开启 Go 模块支持
go env -w GO111MODULE=on

# 2. 设置国内高速镜像代理（推荐七牛云 goproxy.cn，direct 表示私有包回源）
go env -w GOPROXY=https://goproxy.cn,direct

# 3. 验证配置是否成功写入
go env GOPROXY
```

> [!TIP]
> `go env -w` 会将配置持久化保存在用户级的环境配置文件中（Windows 保存在 `%APPDATA%\go\env`，Linux/macOS 保存在 `~/.config/go/env`），重启电脑或重开终端依然有效。

---

## 4. 本地极速自检与代码实操

进入本专题目录，运行自检程序快速查看当前本机的 Go 环境健康状态：

```bash
go run ./00-toolchain-and-engineering-primer/01-environment-and-goproxy/main.go
```
