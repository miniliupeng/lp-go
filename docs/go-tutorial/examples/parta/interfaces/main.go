// 接口。运行：
//
//	go run ./docs/go-tutorial/examples/parta/interfaces
package main

import (
	"fmt"
	"io"
	"strings"
)

// Speaker 只要求一个方法。Dog 没有写 “implements”。
type Speaker interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name + " says woof"
}

func Announce(s Speaker) {
	fmt.Println(s.Speak())
}

func CountBytes(r io.Reader) (int, error) {
	buf := make([]byte, 8)
	total := 0
	for {
		n, err := r.Read(buf)
		total += n
		if err == io.EOF {
			return total, nil
		}
		if err != nil {
			return total, err
		}
	}
}

func main() {
	Announce(Dog{Name: "ada"})

	n, err := CountBytes(strings.NewReader("hello"))
	fmt.Printf("bytes=%d err=%v\n", n, err)

	var s Speaker = Dog{Name: "bea"}
	switch v := s.(type) {
	case Dog:
		fmt.Printf("type=%T name=%s\n", v, v.Name)
	default:
		fmt.Printf("type=%T\n", v)
	}
}
