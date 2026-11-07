# 03 泛型

示例在 Go 1.27.1 上编译。类型参数从 Go 1.18 就有。**方法自己的类型参数是 Go 1.27 才有的。**

## 约束

[Set](../examples/partb/generics/set.go) 的元素必须能用 `==` 比较，所以约束是 `comparable`。map 的键也是这个要求。

```go
func MapSlice[S ~[]E, E any, R any](s S, f func(E) R) []R
```

`~[]E` 表示底层类型是切片。因此自定义的 `type myInts []int` 也能传进去。只写 `[]E` 时，具名切片类型会被拒绝。`any` 表示没有额外约束。

编译器按 GC shape 做实例化（同一形状的类型可以共享代码，用字典传递类型信息）。不要假设「每种类型各生成一份完全独立的机器码」，也不要假设「完全没有运行时成本」。正确性不依赖这个选择。

## 泛型方法

```go
type Box struct{}

func (Box) Identity[T any](v T) T { return v }
```

`b.Identity(3)` 和 `b.Identity("go")` 都成立，类型从实参推断。

两条硬限制，来自 Go 1.27 的语言变更说明：

- 接口方法不能声明类型参数。
- 泛型方法不能用来实现接口方法。

需要多态时，优先用接口或「类型参数在函数上、不在方法上」。

## 提升字段的键

同文件的测试写了 `outer{X: 7}`，其中 `X` 来自嵌入的 `inner`。这是 Go 1.27 的字面量规则，见 [结构体一章](../part-a/06-structs-and-methods.md)。

## 迭代器（Go 1.23）

`iter.Seq[V]` 是 `func(yield func(V) bool)`。`yield` 返回 false 时，迭代器必须停，并且不能再调用 `yield`，否则会 panic。

[CountDown](../examples/partb/generics/set.go) 从 n 倒数到 1。`for n := range CountDown(5)` 在 Go 1.23 及以后可以直接写。`break` 会让 `yield` 返回 false，测试确认只拿到 `5, 4, 3`。

`maps.Keys`、`slices.Values` 这些标准库函数返回的就是迭代器。自己写的时候，不要在 `yield` 返回 false 之后继续发值。

## 什么时候不要用泛型

一组类型没有共享的算法，只是看起来都有 `Name()`，用接口更直接。泛型适合容器、`Map` / `Filter` 这种算法，以及必须在编译期保住具体类型的地方。

## 运行

```bash
go test ./docs/go-tutorial/examples/partb/generics -count=1
```

## 面试常问

- `comparable` 和 `any` 差在哪：前者可以用 `==`，能做 map 键。
- `~` 的作用：放行底层类型相同的具名类型。
- 泛型方法能不能塞进接口：不能。
- `yield` 返回 false 之后还能不能再调用：不能，会 panic。调用方 `break` 就是在要求停止。

## 练习

给 `Set` 写 `Delete`。再写一个编译失败的例子（放在注释里即可）：在接口里声明 `Identity[T any](T) T`。
