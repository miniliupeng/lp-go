# 12 测试

## 这一章你要会什么

用 `testing` 写表驱动测试和示例测试，分清 `t.Fatal`、`t.Helper`、`t.Cleanup` 和 `t.Context`，会看覆盖率。

## 文件放哪

测试文件以 `_test.go` 结尾，和被测代码放在同一个目录。函数名是 `TestXxx(t *testing.T)`。`go test` 只编译和运行这些文件，它们不会进 `go build` 生成的程序。

```bash
go test ./docs/go-tutorial/examples/parta/textutil
go test -cover ./docs/go-tutorial/examples/parta/textutil
```

`-cover` 报告的是语句有没有被执行到，不是分支覆盖率。100% 只说明每条语句都跑过，不说明每种输入都对。

## 表驱动

[lines_test.go](../examples/parta/textutil/lines_test.go) 把输入和期望放进切片，用 `t.Run` 分成子测试。失败时 `t.Fatalf` 带上两边的值。`t.Errorf` 记录失败并继续当前测试函数；`t.Fatal` 立刻结束当前测试函数。辅助函数的第一行写 `t.Helper()`，报错行号会指向调用它的那一行。

比较切片用 `slices.Equal`。它把 nil 和空切片看成相等，因为长度都是 0。`NonEmptyLines("")` 返回的是 nil，测试的期望也写成 nil，所以能通过。如果函数必须区分「没分配」和「空切片」，不要用 `slices.Equal`，改成同时比较 `len` 和是否为 nil。

`"a\r\n b"` 这组用例的结果是 `["a", "b"]`：`strings.Split` 按 `\n` 切开后，`TrimSpace` 去掉 `"a\r"` 末尾的回车。

Go 1.22 起每次循环迭代有自己的循环变量，`t.Run` 里再调用 `t.Parallel()` 时，子测试看到的是当次的 `tt`。更早的版本必须在循环体里再写一次 `tt := tt`。子测试并行之后，不要让它们写同一个没有同步的变量。

## 示例测试

`ExampleNonEmptyLines` 没有 `*testing.T`。注释：

```go
// Output:
// a
// b
```

`go test` 运行这段函数，把标准输出和注释逐字比较。多一个空格就失败。用它锁住给用户看的格式，不要用它代替上面的表驱动测试。

## 清理和取消

函数里的 `defer` 在这个函数返回时执行。`t.Cleanup(fn)` 在该测试和它的子测试都结束之后执行，后注册的先执行。测试打开了临时目录、又在里面 `t.Run` 了几个子测试时，用 `t.Cleanup` 删除目录，避免子测试还没跑完目录就被收掉。

Go 1.24 起 `t.Context()` 返回一个 context。文档写明：它在 `Cleanup` 注册的函数被调用之前取消。把这个 context 传给会监听取消的 API，测试结束时那些调用能停下来。`Cleanup` 函数自己还可以等这些调用退出。

## 基准和模糊

`func BenchmarkAdd(b *testing.B)` 里的循环次数是 `b.N`，由测试框架决定。运行时加 `-bench=.`。`-benchmem` 会多打出每次操作的分配次数。数字跟机器有关，只在同一台机器上比较两次运行。详见 [性能一章](../part-c/07-pprof-trace-bench.md)。

`func FuzzXxx(f *testing.F)` 先用 `f.Add` 放种子，再在 `f.Fuzz` 里写断言。`go test -fuzz=FuzzXxx` 会不断变异输入。本教程没有附带模糊测试；写解析器时再加，失败的输入会被存成语料，下一次 `go test` 会重放。

## 练习

给 `NonEmptyLines` 再加一组用例：`"a\n \n"`。先写出期望，再跑 `go test`。空行修剪后会被丢掉。
