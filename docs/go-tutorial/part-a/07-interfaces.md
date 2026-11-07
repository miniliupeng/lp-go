# 07 接口

## 这一章你要会什么

定义小接口，让已有类型自动满足它，会做类型断言。

## 隐式实现

接口是一组方法。类型实现了这些方法，就实现了接口，不需要写 `implements`。

```go
type Speaker interface {
    Speak() string
}
```

`Dog` 有 `Speak() string`，所以 `Dog` 就是 `Speaker`。接口由使用方定义，经常只有一个方法。标准库里的 `io.Reader` 只有 `Read([]byte) (int, error)`。

函数接受接口、返回具体类型，是常见的分工。`CountBytes` 接受 `io.Reader`，因此字符串、文件、网络连接都能传进来，只要它们能 `Read`。

## Read 的约定

`Read` 可以同时返回 `n > 0` 和一个错误。调用方先处理这 `n` 个字节，再看错误。`io.EOF` 表示没有更多数据。示例里的循环是：

```go
n, err := r.Read(buf)
total += n
if err == io.EOF {
    return total, nil
}
```

## 断言和 type switch

接口值里既有动态类型也有动态值。

```go
v, ok := s.(Dog) // 安全断言
switch v := s.(type) {
case Dog:
    // v 的类型是 Dog
}
```

断言失败且没有 `ok` 会 panic。业务代码用两个返回值的形式。

空接口 `any`（`interface{}` 的别名）表示「任意类型」。它不提供方法。能用具体类型或小接口时，不要先改成 `any`。

## 运行

```bash
go run ./docs/go-tutorial/examples/parta/interfaces
```

输出：

```text
ada says woof
bytes=5 err=<nil>
type=main.Dog name=bea
```

## 下一章会碰到的坑

把 `var p *Dog = nil` 放进 `Speaker` 或 `error` 之后，接口本身不是 nil。判断「有没有错误」会因此出错。完整说明和测试在 [错误一章](08-errors.md) 与 [接口底层](../part-c/06-iface-reflect.md)。

## 练习

再写一个 `Cat`，同样有 `Speak()`。把 `Dog` 和 `Cat` 放进 `[]Speaker`，循环 `Announce`。
