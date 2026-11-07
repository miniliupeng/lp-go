// Package gen 演示类型参数、Go 1.23 迭代器和 Go 1.27 泛型方法。
package gen

import "iter"

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: map[T]struct{}{}}
}

func (s *Set[T]) Add(v T) { s.m[v] = struct{}{} }

func (s *Set[T]) Has(v T) bool {
	_, ok := s.m[v]
	return ok
}

func MapSlice[S ~[]E, E any, R any](s S, f func(E) R) []R {
	out := make([]R, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

type Box struct{}

// Identity 是泛型方法。它不能用来实现接口方法：接口方法不能带类型参数。
func (Box) Identity[T any](v T) T { return v }

type inner struct{ X int }

type outer struct{ inner }

func Promoted() outer { return outer{X: 7} }

// CountDown 从 n 倒数到 1。yield 返回 false 后必须停止。
func CountDown(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := n; i > 0; i-- {
			if !yield(i) {
				return
			}
		}
	}
}
