package main

import (
	"reflect"
	"testing"
)

type Greeter interface {
	SayHello() int
}

type RealGreeter struct{}

func (r RealGreeter) SayHello() int {
	return 1
}

// TestNilInterfaceTrap 严格单测证实 typed-nil 赋给接口后不为 nil
func TestNilInterfaceTrap(t *testing.T) {
	err := returnTypedNilError()
	if err == nil {
		t.Fatal("预期 typed-nil error 赋给接口后应当 != nil")
	}
}

// BenchmarkDirectMethodCall 测试原生静态方法调用
func BenchmarkDirectMethodCall(b *testing.B) {
	g := RealGreeter{}
	b.ReportAllocs()

	for b.Loop() {
		_ = g.SayHello()
	}
}

// BenchmarkInterfaceMethodCall 测试接口 itab 动态分发调用
func BenchmarkInterfaceMethodCall(b *testing.B) {
	var g Greeter = RealGreeter{}
	b.ReportAllocs()

	for b.Loop() {
		_ = g.SayHello()
	}
}

// BenchmarkReflectionMethodCall 测试反射动态寻找并调用方法开销
func BenchmarkReflectionMethodCall(b *testing.B) {
	g := RealGreeter{}
	val := reflect.ValueOf(g)
	method := val.MethodByName("SayHello")
	b.ReportAllocs()

	for b.Loop() {
		_ = method.Call(nil)
	}
}
