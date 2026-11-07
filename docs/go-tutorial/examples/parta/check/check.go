// Package check 演示把错误当成值：哨兵错误、类型错误、包装。
package check

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrEmpty 是调用方可以用 errors.Is 识别的哨兵错误。
var ErrEmpty = errors.New("empty")

// FieldError 描述一个字段为什么不合法。用指针接收者，调用方传 *FieldError。
type FieldError struct {
	Field  string
	Reason string
}

func (e *FieldError) Error() string {
	return e.Field + ": " + e.Reason
}

// ParseAge 接受十进制年龄。空字符串包装 ErrEmpty；格式错误包装 strconv 的错误；
// 超出 0 到 150 返回 *FieldError。
func ParseAge(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("parse age: %w", ErrEmpty)
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse age %q: %w", s, err)
	}
	if n < 0 || n > 150 {
		return 0, &FieldError{Field: "age", Reason: "out of range"}
	}
	return n, nil
}
