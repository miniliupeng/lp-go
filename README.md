# Go 教程（orphan 预览）

这是一条**没有 main 历史**的分支。仓库里只有这份从零写的教程，不包含 main 上的 `00`–`07` 专题。

**不要把本分支合并进 main。** 它和 main 没有共同祖先，合并会用这棵树换掉 main 上的全部文件。

从 [docs/go-tutorial/README.md](docs/go-tutorial/README.md) 按中文目录读：`地基`、`能写服务`、`运行时`、`工程`、`面试`。示例目录 `examples/` 仍用英文。模块路径是 `tutorial`，语言版本以 `go.mod` 的 **Go 1.27.1** 为准。

在仓库根目录：

```bash
go test ./examples/...
go vet ./examples/...
go run ./examples/a1hello
```
