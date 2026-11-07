package escape

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLeakMovedToHeap(t *testing.T) {
	cmd := exec.Command("go", "build", "-gcflags=-m", "-o", os.DevNull, ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	text := string(out)
	if !strings.Contains(text, "Leak") {
		t.Fatalf("missing Leak:\n%s", text)
	}
	if !strings.Contains(text, "moved to heap") && !strings.Contains(text, "escapes to heap") {
		t.Fatalf("unexpected escape wording:\n%s", text)
	}
	if Leak() == nil || *Leak() != 1 || Stay(2, 3) != 5 {
		t.Fatal("runtime values")
	}
}
