# 10 常用标准库

## 这一章你要会什么

用 `flag`、`encoding/json`、`path/filepath`、`strings`、`time` 完成一个小命令，并记住每个包里最容易写错的那一条语义。

## 示例在做什么

[examples/parta/stdlib/main.go](../examples/parta/stdlib/main.go) 做五件事：解析旗标、把一段 JSON 解成结构体、按旗标决定标题要不要变成大写、取路径的最后一段、按固定时区格式化一个时间。

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

## flag

`flag.String` 和 `flag.Bool` 返回的是指针。声明时的第二个参数是默认值，读的时候用 `*name`、`*upper`。`flag.Parse` 之后才有命令行里的值。旗标要出现在普通参数前面。布尔旗标写成 `-upper` 就是 true，不必写 `-upper=true`。

示例不带 `-upper` 时标题保持 JSON 里的小写。默认名字已经是 `Ada`，所以 `-name Ada` 和省略它的结果一样。

## json

`Task` 的字段必须大写，`encoding/json` 才能看见它们。`json:"title"` 决定对象里的名字，和 Go 字段名可以不同。`Unmarshal` 失败时示例直接 `panic`：这段 JSON 是写在源码里的常量，失败说明程序本身写错了。解析用户输入时不要 `panic`，把 `error` 返回去，见 [错误](08-errors.md)。

`encoding/json` 的 `Unmarshal` 在 Go 1.27 里由 v2 实现承接，v1 的可见行为还在：对象里重复的键，后面的值覆盖前面的值。更严的 `encoding/json/v2` 会拒绝重复键，对照见 [Part B 的 JSON 一节](../part-b/01-cli-files-json.md)。

## 路径

`path/filepath` 处理本机文件系统路径，`path` 处理 URL 路径。示例在任何操作系统上都用 `filepath.Join("docs", "readme.md")` 再取 `filepath.Base`，得到 `readme.md`。不要用字符串拼接 `/` 来拼磁盘路径。

## 时间

`time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)` 构造的是一个绝对时刻，位置参数是年、月、日、时、分、秒、纳秒、时区。`Format` 的参数是布局，不是 `strftime` 的 `%Y`。布局参照的是 Go 规定的那一个时刻：2006 年 1 月 2 日 15:04:05，时区 MST（数值上是 UTC-7）。`time.RFC3339` 就是 `2006-01-02T15:04:05Z07:00`。时间本身在 UTC，所以结果以 `Z` 结尾。

测试和示例里用 `time.UTC`。要显示运行机器上的本地时间再用 `time.Local`，那一行输出会随机器变化。

## 字符串

`strings.ToUpper` 按 Unicode 做大小写转换，`"learn go"` 变成 `"LEARN GO"`。按空白切开用 `strings.Fields`，词频工具用的是它。`strings.TrimSpace` 去掉两端的 Unicode 空白，包括 `\r`，所以 `"a\r"` 修剪后是 `"a"`。这一点被 [textutil 的测试](../examples/parta/textutil/lines_test.go) 锁住。

## 文件

读整个小文件用 `os.ReadFile`。按行读用 `bufio.Scanner`。`Scanner` 默认单行上限是 64 KiB，超长行会让 `Scan` 返回 false，原因用 `sc.Err()` 取。写文件用 `os.WriteFile`，权限用 `0o644` 这种八进制写法。

```go
f, err := os.Open(path)
if err != nil {
    return err
}
defer f.Close()
```

先判断 `Open` 的错误，成功后再 `defer Close`。打开失败时 `f` 是 nil，对它调用 `Close` 会 panic。

## 练习

给示例加一个 `-repeat` 整数旗标，默认值 1，把任务列表打印这么多遍。旗标仍然在 `flag.Parse` 之前声明。
