package speak

import (
	"io"
	"testing"
)

func TestAnnounce(t *testing.T) {
	if Announce(Dog{Name: "ada"}) != "ada says woof" {
		t.Fatal(Announce(Dog{Name: "ada"}))
	}
	var s Speaker = Dog{Name: "bea"}
	d, ok := s.(Dog)
	if !ok || d.Name != "bea" {
		t.Fatalf("%v %t", d, ok)
	}
}

func TestPointerMethodNeedsPointer(t *testing.T) {
	n := locked{n: 'A'}
	b, err := ReadOne(&n)
	if err != nil || b != 'A' {
		t.Fatalf("%q %v", b, err)
	}
	var r io.Reader = &n
	if _, ok := r.(*locked); !ok {
		t.Fatal("type assert")
	}
}

func TestTypedNilIsNotNilInterface(t *testing.T) {
	var p *locked
	var r io.Reader = p
	if r == nil {
		t.Fatal("typed nil stored in an interface is not a nil interface")
	}
	_, err := ReadOne(r)
	if err == nil || err.Error() != "nil reader" {
		t.Fatal(err)
	}
}
