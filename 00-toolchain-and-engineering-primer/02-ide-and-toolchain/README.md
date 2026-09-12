# 专题 02：现代 IDE 工具链配置与 Delve 源码调试实战

> 一个顺手的现代 IDE 环境能让编码效率提升数倍。本专题手把手指导如何配置官方推荐的 VS Code / GoLand 开发环境、配置语言服务器 `gopls` 获得毫秒级代码补全，并掌握使用 Delve 进行断点调试。

---

## 1. 现代编辑器选型与生态

大厂工业界目前主要有两种主流 Go 开发环境：
1. **VS Code + Go 官方扩展（轻量级、开源主流）**：由 Go 官方维护的扩展，具备完善的语法高亮、定义跳转、保存自动格式化与断点调试。
2. **GoLand / IntelliJ IDEA（开箱即用、工业级全能）**：JetBrains 商业 IDE，具有极强的大型工程重构与调用链分析能力。

对于初学者，强烈推荐 **VS Code** 起步。

---

## 2. VS Code 核心插件与配置优化

### ① 安装核心插件
在 VS Code 插件市场搜索并安装：
- **`Go`**（扩展 ID：`golang.Go`，由 Go 官方团队维护）
- **`Error Lens`**（可选但强烈推荐：直接在代码行末高亮显示语法错误，排错极快）

### ② 安装 Go 辅助开发工具链
安装好 `Go` 扩展后，打开任意 `.go` 文件，按快捷键 `Ctrl+Shift+P`（macOS 为 `Cmd+Shift+P`），输入并执行：
```text
> Go: Install/Update Tools
```
勾选所有工具（包括 `gopls`、`dlv`、`staticcheck` 等），点击确定安装。

### ③ 生产级 `settings.json` 配置清单
按 `Ctrl+Shift+P` 打开 `Preferences: Open User Settings (JSON)`，追加以下配置：

```json
{
  // 保存时自动调用 gofmt 和 goimports 整理代码与引用包
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.snippetSuggestions": "top",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "always"
    }
  },
  // 开启 gopls 现代化语法补全与类型推断
  "gopls": {
    "ui.semanticTokens": true,
    "ui.completion.usePlaceholders": true
  }
}
```

---

## 3. 使用 Delve (dlv) 进行断点单步调试

Go 官方标准的调试器是 **Delve (`dlv`)**。

### ① VS Code 图形化断点调试步骤
1. 打开 `main.go` 文件；
2. 在代码行号左侧点击鼠标左键，生成**红点断点**；
3. 按下键盘 **`F5`**（或点击左侧调试面板的“运行和调试”）；
4. 程序将在断点处暂停，你可以：
   - **F10（单步跳过）**：执行下一行代码；
   - **F11（单步进入）**：进入当前调用的函数内部；
   - **Shift+F11（单步跳出）**：执行完当前函数并返回调用处；
   - 在左侧“变量”面板中实时查看变量在内存中的实时取值！

---

## 4. 本地实操体验

进入本专题目录，运行示例代码，观察代码规范与调试支持：

```bash
go run ./00-toolchain-and-engineering-primer/02-ide-and-toolchain/main.go
```
