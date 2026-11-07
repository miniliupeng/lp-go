package iface

import (
	"reflect"
	"testing"
)

type myErr struct{}

func (myErr) Error() string { return "boom" }

func TestTypedNil(t *testing.T) {
	var p *myErr
	var err error = p
	if err == nil {
		t.Fatal("interface holding (*myErr)(nil) is not a nil interface")
	}
	var bare error
	if bare != nil {
		t.Fatal("nil interface")
	}
}

type counter struct{ n int }

func (c *counter) Inc() { c.n++ }

func TestPointerMethodSet(t *testing.T) {
	var c counter
	var api interface{ Inc() } = &c
	api.Inc()
	if c.n != 1 {
		t.Fatal(c.n)
	}
}

func TestReflectCanSet(t *testing.T) {
	n := 2
	v := reflect.ValueOf(&n).Elem()
	if !v.CanSet() {
		t.Fatal("pointer elem should be settable")
	}
	v.SetInt(3)
	if n != 3 {
		t.Fatal(n)
	}
	if reflect.ValueOf(n).CanSet() {
		t.Fatal("non-addressable value must not be settable")
	}
}
