// 词频小工具。运行：
//
//	go run ./docs/go-tutorial/examples/parta/wordstat -file docs/go-tutorial/examples/parta/wordstat/testdata/sample.txt -top 3
//
// 省略 -file 时从标准输入读。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"lp-go/docs/go-tutorial/examples/parta/wordstat/stats"
)

func main() {
	file := flag.String("file", "", "要统计的文件；空字符串表示标准输入")
	topN := flag.Int("top", 10, "输出前多少个词")
	flag.Parse()

	if *topN < 0 {
		fmt.Fprintln(os.Stderr, "top 必须 >= 0")
		os.Exit(2)
	}

	var in io.Reader = os.Stdin
	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "打开文件: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		in = f
	}

	freq, err := stats.Count(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "统计: %v\n", err)
		os.Exit(1)
	}
	for _, item := range stats.Top(freq, *topN) {
		fmt.Printf("%-12s %d\n", item.Word, item.Count)
	}
}
