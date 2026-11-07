// Package generics 演示类型参数、约束、Go 1.23 的迭代器，以及 Go 1.27 的泛型方法。
package generics

import "iter"

// Set 是一个集合。元素必须能用 == 比较，所以约束是 comparable。
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet 返回空集合。
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

// Add 放入元素。
func (s *Set[T]) Add(v T) {
	s.m[v] = struct{}{}
}

// Has 报告元素是否存在。
func (s *Set[T]) Has(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Len 返回元素个数。
func (s *Set[T]) Len() int {
	return len(s.m)
}

// MapSlice 对切片每个元素调用 f。S 的底层类型必须是切片，所以写成 ~[]E。
func MapSlice[S ~[]E, E any, R any](s S, f func(E) R) []R {
	out := make([]R, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

// Box 本身没有类型参数。Go 1.27 起方法可以自己声明类型参数。
// 这种方法和接口方法不是一回事：接口方法不能带类型参数，泛型方法也不能用来实现接口方法。
type Box struct{}

// Identity 原样返回 v。
func (Box) Identity[T any](v T) T {
	return v
}

// CountDown 从 n 倒数到 1。Go 1.23 起可以 for range 一个 iter.Seq。
// yield 返回 false 时必须停，并且不能再调用 yield。
func CountDown(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := n; i > 0; i-- {
			if !yield(i) {
				return
			}
		}
	}
}
