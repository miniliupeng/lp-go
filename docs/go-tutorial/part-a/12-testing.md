# 12 测试

## 这一章你要会什么

用 `testing` 写表驱动测试和示例测试，会看覆盖率。

## 文件放哪

测试文件以 `_test.go` 结尾，和被测代码同一个包。函数名是 `TestXxx(t *testing.T)`。

```bash
go test ./docs/go-tutorial/examples/parta/textutil
go test -cover ./docs/go-tutorial/examples/parta/textutil
```

## 表驱动

[lines_test.go](../examples/parta/textutil/lines_test.go) 把输入和期望放进切片，用 `t.Run` 分成子测试。失败时用 `t.Fatalf` 带上名字和两边的值。

`t.Errorf` 记录失败并继续。`t.Fatal` 立刻结束当前测试。辅助函数开头写 `t.Helper()`，这样报错行号会指向调用处。

比较切片用 `slices.Equal`。它把 nil 和空切片看成相等（长度都是 0）。如果你的函数必须区分 nil 和空切片，不要用它，改成自己比较。

## 示例测试

`ExampleNonEmptyLines` 的注释：

```go
// Output:
// a
// b
```

`go test` 会运行这段并比对标准输出。输出多一个空格就会失败。适合锁住面向用户的格式。

## 基准测试先认识一下

`func BenchmarkAdd(b *testing.B)` 里用 `b.N` 循环。运行时加 `-bench=.`。`-benchmem` 会多打出每次操作的分配次数。数字跟机器有关，只在同一台机器上比较两次运行。详见 [性能一章](../part-c/07-pprof-trace-bench.md)。

## 练习

给 `NonEmptyLines` 加一行用例：`"a\r\n b"`。先预测结果再跑。`strings.TrimSpace` 会去掉 `\r`。
