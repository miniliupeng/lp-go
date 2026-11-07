# 08 错误处理

## 这一章你要会什么

把失败写成 `error` 返回值，用 `errors.Is` 和 `errors.As` 识别原因，知道 `panic` 不负责普通业务分支。

## error 是一个接口

```go
type error interface {
    Error() string
}
```

函数失败时返回非 nil 的 `error`，成功时返回 `nil`。调用方写成：

```go
n, err := check.ParseAge(s)
if err != nil {
    return err
}
```

不要忽略 `err`。

## 三种常见的错误

源码在 [examples/parta/check/check.go](../examples/parta/check/check.go)。

1. **哨兵错误**：包级 `var ErrEmpty = errors.New("empty")`。调用方用 `errors.Is(err, check.ErrEmpty)`。`fmt.Errorf("...: %w", ErrEmpty)` 会把哨兵包进链里，`Is` 沿着链找。
2. **带字段的类型**：`*FieldError` 实现 `Error()`。调用方用 `errors.As(err, &field)` 取出具体类型。`As` 的第二个参数必须是指向实现了 `error` 的类型的指针。
3. **再包一层上下文**：`fmt.Errorf("parse age %q: %w", s, err)`。`%w` 才能被 `Is` / `As` 解开。用 `%v` 只会得到一段文本，链就断了。

`errors.Join(err1, err2)`（Go 1.20 起）把多个错误收成一个，`Is` 和 `As` 会检查其中每一个。

## 测试

```bash
go test ./docs/go-tutorial/examples/parta/check -count=1
```

`ParseAge` 的约定：

| 输入 | 结果 |
| --- | --- |
| `"20"` | `20, nil` |
| 空白 | 链上有 `ErrEmpty` |
| `"abc"` | 链上有 `*strconv.NumError` |
| `"200"` | `*FieldError`，字段 `age`，原因 `out of range` |

## panic

`panic` 适合程序自己的不变量被破坏，例如 `limit.NewBucket` 收到非正的速率。解析用户输入、打开文件失败、HTTP 4xx，都返回 `error`。

`recover` 只能接住**同一个 goroutine** 里 defer 函数中的 panic。用它包住一个请求的边界是可以的；用它代替 `if err != nil` 会让调用方无法处理失败。

## 和 nil 接口有关的坑

`var err error` 是 nil。`var p *FieldError; var err error = p` 不是 nil 接口：动态类型是 `*FieldError`，动态值是 nil。`if err != nil` 会进去，而 `err.Error()` 会 panic。构造错误时返回 `nil` 接口，不要返回带类型的 nil 指针。测试在 [examples/partc/iface](../examples/partc/iface/iface_test.go)。

## 练习

给 `ParseAge` 增加「必须是整数」以外的一条规则，并用 `errors.As` 写一个测试。不要用字符串比较错误文本。
