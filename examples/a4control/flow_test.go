package flow

import "testing"

func Classify(n int) string {
	if n%2 == 1 {
		return "odd"
	}
	return "even"
}

func SumTo(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

func RangeN(n int) []int {
	var out []int
	for i := range n {
		out = append(out, i)
	}
	return out
}

func Bucket(n int) string {
	switch n {
	case 1:
		return "one"
	case 2, 3:
		return "two-or-three"
	default:
		return "other"
	}
}

func Command(cmd string) string {
	switch {
	case cmd == "sum":
		return "cmd=sum"
	default:
		return "cmd=other"
	}
}

// Fall prints the fallthrough path. The next case expression is not evaluated.
func Fall(n int) string {
	switch n {
	case 2, 3:
		return fallTail()
	default:
		return "only-default"
	}
}

func fallTail() string {
	out := "fall="
	switch 0 {
	case 0:
		out += "head"
		fallthrough
	case 1:
		out += "-through"
	}
	return out
}

func BreakOuter() bool {
	found := false
outer:
	for a := 0; a < 3; a++ {
		for b := 0; b < 3; b++ {
			if a == 1 && b == 1 {
				found = true
				break outer
			}
		}
	}
	return found
}

func TestFlow(t *testing.T) {
	if Classify(3) != "odd" || SumTo(3) != 6 {
		t.Fatal("classify/sum")
	}
	got := RangeN(3)
	if len(got) != 3 || got[0] != 0 || got[2] != 2 {
		t.Fatal(got)
	}
	if Bucket(3) != "two-or-three" || Command("sum") != "cmd=sum" {
		t.Fatal("switch")
	}
	if Fall(3) != "fall=head-through" || !BreakOuter() {
		t.Fatalf("%s %t", Fall(3), BreakOuter())
	}
}
