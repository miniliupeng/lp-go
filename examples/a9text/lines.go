// Package text 提供给测试一章的纯函数。
package text

import "strings"

// NonEmptyLines 按换行切开，去掉每行首尾空白，丢掉空行。
func NonEmptyLines(s string) []string {
	if s == "" {
		return nil
	}
	raw := strings.Split(s, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}
