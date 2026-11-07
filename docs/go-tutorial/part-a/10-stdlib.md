# 10 常用标准库

## 这一章你要会什么

用 `flag`、`encoding/json`、`path/filepath`、`strings`、`time` 完成一个小命令。这些包覆盖了日常脚本和 HTTP 服务的一半样板。

## 示例在做什么

[examples/parta/stdlib/main.go](../examples/parta/stdlib/main.go) 解析两个旗标，把一段 JSON 解成结构体，按旗标决定要不要把标题转成大写，再打印路径的最后一段和固定时区的时间。

```bash
go run ./docs/go-tutorial/examples/parta/stdlib -name Ada -upper
```

输出：

```text
name=Ada
- LEARN GO (done)
- WRITE TOOL (todo)
base=readme.md
when=2026-01-02T15:04:05Z
```

## 这几处语义要记准

**flag。** `flag.String` 返回的是指针，读的时候用 `*name`。`flag.Parse` 之后才有值。旗标必须出现在普通参数前面。

**json。** 结构体字段的 `json:"title"` 决定对象里的名字。字段必须大写，编码包才能看见它。`encoding/json` 的 `Unmarshal` 在 Go 1.27 里由 v2 实现承接，但 v1 的行为还在：重复键时后面的值覆盖前面的值。更严的 `encoding/json/v2` 会拒绝重复键，见 [Part B 的 JSON 一节](../part-b/01-cli-files-json.md)。

**路径。** 用 `path/filepath` 拼接本机路径，用 `path` 处理 URL 路径。`filepath.Base(filepath.Join("docs", "readme.md"))` 得到 `readme.md`。

**时间。** `time.Date(..., time.UTC)` 构造的是绝对时刻。`Format(time.RFC3339)` 用固定的参考时间布局 `2006-01-02T15:04:05Z07:00`，不是 `strftime` 的 `%Y`。想要本地时区就用 `time.Local`，测试里用 `UTC` 才能稳定。

**字符串。** `strings.ToUpper` 按 Unicode 处理大小写。`strings.Fields` 按空白切开，下一章的词频工具用它。

## 文件

读整个小文件用 `os.ReadFile`。按行读用 `bufio.Scanner`。`Scanner` 默认单行上限是 64 KiB，超长行会让 `Scan` 失败，错误用 `sc.Err()` 取出。写文件用 `os.WriteFile`，权限比如 `0o644`。

打开后要关：`defer f.Close()`。先处理 `Open` 的错误，再 `defer`，避免对 nil 文件调用 `Close`。

## 练习

给示例加一个 `-repeat` 整数旗标，把任务列表打印多遍。缺省值用 1。
