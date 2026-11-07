package profile

import "testing"

var sink int

func BenchmarkAdd(b *testing.B) {
	n := 0
	for i := 0; i < b.N; i++ {
		n += i
	}
	sink = n
}
