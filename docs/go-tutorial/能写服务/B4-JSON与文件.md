# B4 JSON、文件与 os.Root

示例在 [examples/b4json](../../../examples/b4json/json_test.go)。语言版本是 Go 1.27.1。标签、`omitempty` 和时间布局在 [标准库与测试](../地基/A9-标准库与测试.md)。

## 这一章你要会什么

说清 `encoding/json` 和 `encoding/json/v2` 对重复键的差别，会用 `os.Root` 把文件访问限制在一个目录里。

## 两个 JSON 实现

`encoding/json` 仍是 v1 的 API。Go 1.27 起它的实现换成 v2，但 v1 的可见行为还在。对象里重复的键，`json.Unmarshal` 让后面的值覆盖前面的值：`{"a":1,"a":2}` 得到 `A == 2`。缺少的键保持零值，`Unmarshal` 不报错。类型不对，例如把字符串解进 `int`，返回错误。

`encoding/json/v2` 的 `Unmarshal` 默认拒绝重复键，测试断言它返回错误。需要旧实现时，构建期可以设 `GOEXPERIMENT=nojsonv2`。这是过渡开关。新代码按你选择的包来写，不要同时依赖两种重复键行为。

两边都认 `json:"a"` 这种结构体标签。字段仍然必须导出。未导出字段在编码和解码时都被忽略。

数字的默认目标如果是 `float64`（解进 `any` 或 `map[string]any`），很大的整数会失去精度。标识符、雪花 ID 用 `string` 或 `json.Number`，不要解进 `float64` 再转回整数。

`Decoder` 可以从 `io.Reader` 里连续解多个值，并配合 `MaxBytesReader` 限制大小。`Unmarshal` 要求手里已经是完整的 `[]byte`。请求体用 `Decoder`，配置文件这种已有字节切片再用 `Unmarshal`。

未知字段默认被忽略。如果协议要求严格，v1 可以在 `Decoder` 上 `DisallowUnknownFields`。打开之后，对方多写一个键就会失败。这是兼容性和严格性之间的选择，两种都要在接口文档里写明。

## 文件名和 os.Root

用户提供的文件名不能直接拼进路径。`filepath.Join(dir, userInput)` 在 `userInput` 含有 `..` 时可以离开 `dir`。`filepath.Clean` 只是整理字符串，不会把访问限制在目录里。

Go 1.24 起 `os.OpenRoot(dir)` 返回一棵限定在该目录里的文件树。`ReadFile("a.txt")` 能读到目录里的文件。`ReadFile("../a.txt")` 返回错误，路径出不去。符号链接如果指向目录外面，同样出不去。处理上传文件名、解压条目时用它。

`OpenRoot` 之后要 `Close`。测试用 `t.TempDir()` 和 `t.Cleanup`。

`os.Open` 返回的文件实现了 `io.Reader`。读到末尾是 `io.EOF`。`ReadFile` 把 EOF 处理掉，调用方看到的是字节和打开或读取时的错误。

写文件要考虑写到一半失败。写临时文件，刷盘并关闭成功之后，再 `os.Rename` 到目标名字。`Rename` 在同一文件系统里是替换目标的常用办法。忽略 `Close` 的错误会把还在缓冲区里的数据当成已经落盘，见 [错误](../地基/A7-错误.md)。

## 对比陷阱

| 写法 | 实际结果 |
| --- | --- |
| 以为 v1 和 v2 对重复键一样 | v1 后写覆盖，v2 默认拒绝。测试锁住了这两条 |
| 缺少字段就认为 `Unmarshal` 会失败 | 缺字段得到零值，不报错 |
| 把大整数解进 `any` | 目标是 `float64`，精度会丢 |
| `filepath.Join` 之后直接 `os.Open` 用户输入 | `..` 可以离开目标目录 |
| 只调用 `Clean` 就认为路径安全 | `Clean` 不限制根目录 |
| 写文件只检查 `Write` 的错误 | `Close` 才把缓冲刷出去 |

## 运行

```bash
go test ./examples/b4json -count=1
```

## 练习

1. 给 v1 再写一组字节切片：缺了 `a` 这个键。预测 `A` 是 0，并且 `Unmarshal` 不返回错误。
2. 在 `os.Root` 下再读一次 `a.txt/../../etc/passwd` 这种拼出来的路径，确认返回错误。不要真的去读系统密码文件，断言的是错误本身。
3. 把同一个对象用 v1 编码再解回来，确认导出字段还在，未导出字段变成零值。

## 面试延伸

重复键：v1 后写覆盖并成功返回，v2 默认返回错误。Go 1.27 的 v1 API 由 v2 实现承接，这条可见行为仍以测试为准。文件题：用户路径用 `os.Root`，`Join` 和 `Clean` 只处理字符串。大整数不要进 `float64`。写文件要检查 `Close`，并用临时文件加 `Rename` 避免留下半截内容。

## 问答题

### 问：`{"a":1,"a":2}` 用 v1 解出来是几？

`A == 2`。后写覆盖前写，并且 `Unmarshal` 返回 nil。缺了 `a` 时 `A` 是 0，也不报错。类型不对，例如字符串解进 `int`，才返回错误。

### 问：v2 对同一份字节做什么？

`encoding/json/v2` 的 `Unmarshal` 默认拒绝重复键，测试断言它返回错误。`GOEXPERIMENT=nojsonv2` 是过渡开关。新代码选定一个包，不要同时依赖两种重复键行为。

### 问：大整数为什么不能解进 `any`？

默认数字目标是 `float64`。超过 53 位尾数的整数会变。标识符用字符串或 `json.Number`。

### 问：`DisallowUnknownFields` 改变什么？

多出来的键会让解码失败。默认是忽略。协议要严格就打开，并写进接口说明。这是 v1 `Decoder` 上的开关。

### 问：`filepath.Join(dir, user)` 为什么挡不住 `..`？

`Join` 和 `Clean` 只整理字符串。用户输入 `../a.txt` 可以离开 `dir`。Go 1.24 起用 `os.OpenRoot(dir)`，`ReadFile("../a.txt")` 返回错误。符号链接指向目录外也同样出不去。用完 `Close`。

### 问：`ReadFile` 和 `os.Open` 对 EOF 的差别？

`Read` 到末尾返回 `io.EOF`。`ReadFile` 把 EOF 处理掉，调用方看到字节和真正的错误。按行读用 `bufio.Scanner`，默认一行最多 64 KiB。

### 问：写一半失败怎么避免留下半截文件？

写到临时文件，`Close` 成功后再 `Rename` 到目标名。`Close` 的错误用 `errors.Join` 留住，因为缓冲区可能还没刷下去。

### 问：未导出字段能被 JSON 看见吗？

不能。编码和解码都忽略。标签 `json:"a"` 只对导出字段生效。`omitempty` 的空值规则在 [标准库与测试](../地基/A9-标准库与测试.md)。
