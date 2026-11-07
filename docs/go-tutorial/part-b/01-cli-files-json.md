# 01 命令行、文件与 JSON

读者已经会写 [Part A](../part-a/10-stdlib.md) 里的旗标和 `encoding/json`。这里补上边界，以及 Go 1.27 的 `encoding/json/v2`。

## JSON 的两个包

| | `encoding/json` | `encoding/json/v2` |
| --- | --- | --- |
| 状态 | 原有 API。Go 1.27 起实现换成 v2，**v1 语义保留** | Go 1.27 起正式提供 |
| 重复的对象键 | 后面的值覆盖前面的，`Unmarshal` 成功 | `Unmarshal` 返回错误 |
| 非法 UTF-8 | 沿用 v1 | v2 拒绝 |
| 本仓库测试 | [jsondiff](../examples/partb/jsondiff/json_test.go) 断言覆盖后 `a == 2` | 同一输入必须得到错误 |

对这个结构：

```go
struct {
    A int `json:"a"`
}{A: 1}
```

两个包的 `Marshal` 在 Go 1.27.1 上都得到 `{"a":1}`。不要把「有没有空格和换行」背成规范；`Encoder` 默认会再加换行，`Marshal` 不加。

v1 函数继续可用。发布说明写明：若实现不兼容，构建期可以设 `GOEXPERIMENT=nojsonv2` 退回旧实现，这个开关预计以后会去掉。新代码要更严的默认行为时，直接导入 `encoding/json/v2`。它的 `Marshal` / `Unmarshal` 都多一个可变的 `Options` 参数。

字段仍然必须导出。`json:"-"` 表示忽略。`omitempty` 在 v1 里的意思是「空值省略」，空值包括 false、0、nil、空串。v2 另外有 `omitzero`，按 Go 零值或 `IsZero()` 判断。两者不要记混。

## 文件路径

用户给的路径可能是 `../` 或其他目录。Go 1.24 起的 `os.Root` 把打开操作限制在一个目录树里：

```go
root, err := os.OpenRoot(dir)
if err != nil { /* 返回错误 */ }
defer root.Close()
f, err := root.Open(name)
```

名字里的组件若跑出这棵树，`Open` 返回错误。这是给文件服务用的，比自己拼接字符串再判断前缀可靠。教程的词频工具打开的是调用方明确传入的一个文件，所以用了 `os.Open`。一旦你要「在某个目录下按用户给的相对路径找文件」，用 `os.Root`。

## 小文件和大文件

`os.ReadFile` 会把文件全部放进内存。日志和上传不能这样做。按块读用 `io.Copy` 或 `bufio`，并设置上限。HTTP 正文用 `http.MaxBytesReader`。

## 运行

```bash
go test ./docs/go-tutorial/examples/partb/jsondiff -count=1
```

## 面试常问

- 重复键谁说了算：先说你用的是哪个包。
- 未导出字段为什么是零值：`encoding/json` 只看导出字段。
- `json.Number` 或 `UseNumber` 用来避免把大整数解成 `float64`。v2 用选项控制数字的解法，不要凭记忆编一个不存在的字段名。

## 练习

用 `encoding/json/v2` 写一个测试：输入 `{"a":1,"a":2}` 时期望错误；再用 `encoding/json` 写一个测试，期望 `a` 为 2。仓库里已经有这两份，先自己写，再对。
