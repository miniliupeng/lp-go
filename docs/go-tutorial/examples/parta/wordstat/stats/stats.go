// Package stats 按空白统计词频。中文没有空格时会整段算一个词；这是分词库的工作，这里不引入第三方依赖。
package stats

import (
	"bufio"
	"io"
	"sort"
	"strings"
	"unicode"
)

// Item 是一个词和它的次数。
type Item struct {
	Word  string
	Count int
}

// Count 读取 r，按空白切词。词两端不是字母也不是数字的字符会被去掉，其余部分转成小写。
func Count(r io.Reader) (map[string]int, error) {
	freq := make(map[string]int)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		for _, field := range strings.Fields(sc.Text()) {
			word := normalize(field)
			if word == "" {
				continue
			}
			freq[word]++
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return freq, nil
}

func normalize(word string) string {
	word = strings.ToLower(word)
	return strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// Top 按次数从高到低取前 n 个。次数相同按词的字典序，结果稳定。
// n < 0 当作 0。n 大于词的种类数时返回全部。
func Top(freq map[string]int, n int) []Item {
	if n < 0 {
		n = 0
	}
	items := make([]Item, 0, len(freq))
	for word, count := range freq {
		items = append(items, Item{Word: word, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Word < items[j].Word
	})
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}
