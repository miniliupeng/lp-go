package speak

import (
	"fmt"
	"io"
)

type Speaker interface {
	Speak() string
}

type Dog struct{ Name string }

func (d Dog) Speak() string { return d.Name + " says woof" }

func Announce(s Speaker) string { return s.Speak() }

type locked struct{ n int }

func (l *locked) Read(p []byte) (int, error) {
	if l == nil {
		return 0, fmt.Errorf("nil reader")
	}
	if len(p) == 0 {
		return 0, nil
	}
	p[0] = byte(l.n)
	return 1, io.EOF
}

func ReadOne(r io.Reader) (byte, error) {
	var buf [1]byte
	n, err := r.Read(buf[:])
	if n == 1 {
		return buf[0], nil
	}
	return 0, err
}
