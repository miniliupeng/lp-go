# B4 JSON、文件与 os.Root

示例在 [examples/b4json](../../../examples/b4json/json_test.go)。语言版本是 Go 1.27.1。

## 两个 JSON 实现

`encoding/json` 仍是 v1 的 API。Go 1.27 起它的实现换成 v2，但 v1 的可见行为还在。对象里重复的键，`json.Unmarshal` 让后面的值覆盖前面的值：`{"a":1,"a":2}` 得到 `A == 2`。

`encoding/json/v2` 的 `Unmarshal` 默认拒绝重复键，测试断言它返回错误。需要旧实现时，构建期可以设 `GOEXPERIMENT=nojsonv2`。这是过渡开关，新代码按你选择的包来写，不要同时依赖两种重复键行为。

两边都认 `json:"a"` 这种结构体标签。字段仍然必须导出。

## os.Root

Go 1.24 起 `os.OpenRoot(dir)` 返回一棵限定在该目录里的文件树。`ReadFile("a.txt")` 能读到目录里的文件。`ReadFile("../a.txt")` 返回错误，路径出不去。处理上传文件名时用它，不要自己拼接 `filepath.Join(dir, userInput)` 之后直接 `os.Open`：用户输入里的 `..` 可以离开目标目录。

`OpenRoot` 之后要 `Close`。测试用 `t.TempDir()` 和 `t.Cleanup`。

```bash
go test ./examples/b4json -count=1
```

## 练习

给 v1 再写一组字节切片：缺了 `a` 这个键。预测 `A` 是 0，并且 `Unmarshal` 不返回错误。
