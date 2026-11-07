// Package check 演示哨兵错误、包装和自定义错误类型。
package check

import (
	"errors"
	"fmt"
)

var ErrEmpty = errors.New("check: empty")

type FieldError struct {
	Field string
	Err   error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("check: %s: %v", e.Field, e.Err)
}

func (e *FieldError) Unwrap() error { return e.Err }

func ParseAge(s string) (int, error) {
	if s == "" {
		return 0, ErrEmpty
	}
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, &FieldError{Field: "age", Err: fmt.Errorf("not a number")}
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}

func Together(errs ...error) error {
	return errors.Join(errs...)
}
