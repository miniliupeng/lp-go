package escape

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLeakEscapes(t *testing.T) {
	cmd := exec.Command("go", "build", "-gcflags=-m", "-o", os.DevNull, ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	text := string(out)
	// Go 1.27.1 对返回参数地址打印 “moved to heap”。更早的编译器常写 “escapes to heap”。
	if !strings.Contains(text, "Leak") || (!strings.Contains(text, "moved to heap") && !strings.Contains(text, "escapes to heap")) {
		t.Fatalf("compiler output:\n%s", text)
	}
	if Stay(4) != 6 {
		t.Fatal(Stay(4))
	}
	p := Leak(3)
	if p == nil || *p != 3 {
		t.Fatal(p)
	}
}
