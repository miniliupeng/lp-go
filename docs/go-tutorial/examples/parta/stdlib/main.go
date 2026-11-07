// 标准库里最常用的几块。运行：
//
//	go run ./docs/go-tutorial/examples/parta/stdlib -name Ada -upper
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type Task struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	name := flag.String("name", "Ada", "显示名")
	upper := flag.Bool("upper", false, "标题改成大写")
	flag.Parse()

	raw := []byte(`[{"title":"learn go","done":true},{"title":"write tool","done":false}]`)
	var tasks []Task
	if err := json.Unmarshal(raw, &tasks); err != nil {
		panic(err)
	}

	fmt.Printf("name=%s\n", *name)
	for _, task := range tasks {
		title := task.Title
		if *upper {
			title = strings.ToUpper(title)
		}
		state := "todo"
		if task.Done {
			state = "done"
		}
		fmt.Printf("- %s (%s)\n", title, state)
	}

	fmt.Printf("base=%s\n", filepath.Base(filepath.Join("docs", "readme.md")))

	when := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	fmt.Printf("when=%s\n", when.Format(time.RFC3339))
}
