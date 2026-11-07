# 06 接口与反射

## 接口值

一个接口值有两样东西：动态类型，和该类型的动态值。两者都空，才是 nil 接口。

```go
var p *myErr          // nil 指针
var err error = p     // 动态类型是 *myErr，动态值是 nil
err == nil            // false
```

测试在 [iface_test.go](../examples/partc/iface/iface_test.go)。`errors.As` 和类型断言都会看到这个具体类型。函数要表示「没有错误」时，返回未赋类型的 `nil`，不要 `return (*myErr)(nil)`。

空接口（`any`）和带方法的接口在运行时是两种表示：前者只有类型和数据，后者还要有方法表（itab）。面试可以说出这个差别。具体字段布局会随版本变，不要背偏移。

## 方法集

规范：

- 类型 `T` 的方法集只包含接收者是 `T` 的方法。
- 类型 `*T` 的方法集包含接收者是 `*T` 和 `T` 的方法。

值可以取地址时，编译器允许 `t.PointerMethod()`，等价于 `(&t).PointerMethod()`。这个值放进接口之后，自动取地址消失。所以 `var c counter; var api interface{ Inc() } = c` 不能编译，`= &c` 可以。测试只覆盖能编译的那一支，编译失败的那支写在这里，避免仓库里出现编不过的文件。

嵌入会提升方法，提升进方法集。外层同名方法覆盖内层。

## 反射

`reflect` 的三条实用规则：

1. 从接口值得到 `reflect.Value`：`reflect.ValueOf(x)`。
2. 从 `Value` 回到接口值：`Interface()`。
3. 要修改，`Value` 必须可设置。`CanSet` 对不可寻址的值、以及未导出字段为 false。

```go
n := 2
v := reflect.ValueOf(&n).Elem()
v.SetInt(3)                  // n == 3
reflect.ValueOf(n).CanSet()  // false
```

`Kind` 是底层种类（`Int`、`Struct`），`Type` 是具体类型。读未导出字段可以用 `CanInterface` 判断；强行改未导出字段不是普通反射能做的事。

反射会让值逃逸，调用也比直接调用慢。编解码、依赖注入、测试工具可以用。热循环里的业务分支不要走反射。

Go 1.27 的泛型方法不能实现接口方法，接口方法也不能自带类型参数。需要按类型派发时，用接口或函数上的类型参数。见 [泛型一章](../part-b/03-generics.md)。

## 面试常问

- nil 接口和接口里的 nil 指针。
- 值接收者和指针接收者进不进接口。
- `CanSet` 为什么对 `ValueOf(n)` 为 false：拷进去的值不可寻址。

## 对比陷阱

- 用 `==` 比较两个接口，动态类型不可比较时会 panic。不确定就用类型断言或 `reflect.DeepEqual`，并知道 `DeepEqual` 对 nil 切片和空切片的区分。
- 在热路径对每个元素 `ValueOf`。

## 运行

```bash
go test ./docs/go-tutorial/examples/partc/iface -count=1
```
