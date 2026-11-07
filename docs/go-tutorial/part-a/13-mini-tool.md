# 13 一个小工具：词频

这一章把前面的旗标、文件、切片、map、错误和测试收成一个命令 `wordstat`。

## 行为

- `-file` 指向文本文件。空字符串表示从标准输入读。
- `-top` 是输出的词数，默认 10，小于 0 则退出码 2。
- 按空白切词。去掉词两端「不是字母也不是数字」的字符，再转成小写。
- 次数高的在前。次数相同按字典序，结果稳定。
- 打开文件或读取失败时，把原因写到标准错误，退出码 1。

中文句子如果中间没有空格，整句会算一个词。分词需要别的库。这个工具明确不处理那件事。

`bufio.Scanner` 的默认行缓冲上限是 64 KiB。更长的一行会让 `Count` 返回错误。

## 为什么统计函数接受 io.Reader

[stats.Count](../examples/parta/wordstat/stats/stats.go) 不接受文件名。测试可以传入 `strings.NewReader`，命令行再把 `*os.File` 传进去。文件只在 `main` 里打开，并 `defer Close`。

## 运行

```bash
go test ./docs/go-tutorial/examples/parta/wordstat/...
go run ./docs/go-tutorial/examples/parta/wordstat \
  -file docs/go-tutorial/examples/parta/wordstat/testdata/sample.txt -top 3
```

样例文件是：

```text
Go is simple.
Go is practical.
Practical tools matter.
```

输出：

```text
go           2
is           2
practical    2
```

`Go,` 和 `go` 会落成同一个词。`simple.` 会变成 `simple`。三个次数为 2 的词按字典序是 `go`、`is`、`practical`。

## 你可以接着改的地方

- 停用词：忽略 `is`。
- 从标准输入读：`echo 'Go go' | go run ./docs/go-tutorial/examples/parta/wordstat -top 1`。
- 并发按文件统计再合并。合并时每个文件一个 map，最后在一个 goroutine 里相加，避免并发写同一个 map。

到这里，Part A 结束。你应该能独立写一个带参数、能失败、带测试的小程序。Part B 把同样的习惯带到 HTTP、泛型和日志上。
