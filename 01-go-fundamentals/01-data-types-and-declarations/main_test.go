package main

import (
	"testing"
	"unsafe"
)

func TestDataTypesAndSizes(t *testing.T) {
	var b bool
	var i8 int8
	var i64 int64
	var f64 float64
	var r rune

	if unsafe.Sizeof(b) != 1 {
		t.Errorf("bool 预期为 1 字节，实际为 %d", unsafe.Sizeof(b))
	}
	if unsafe.Sizeof(i8) != 1 {
		t.Errorf("int8 预期为 1 字节，实际为 %d", unsafe.Sizeof(i8))
	}
	if unsafe.Sizeof(i64) != 8 {
		t.Errorf("int64 预期为 8 字节，实际为 %d", unsafe.Sizeof(i64))
	}
	if unsafe.Sizeof(f64) != 8 {
		t.Errorf("float64 预期为 8 字节，实际为 %d", unsafe.Sizeof(f64))
	}
	if unsafe.Sizeof(r) != 4 {
		t.Errorf("rune 预期为 4 字节，实际为 %d", unsafe.Sizeof(r))
	}
}

func TestIotaEnumeration(t *testing.T) {
	if OrderStateCreated != 0 {
		t.Errorf("OrderStateCreated 预期为 0，实际为 %d", OrderStateCreated)
	}
	if OrderStatePaid != 1 {
		t.Errorf("OrderStatePaid 预期为 1，实际为 %d", OrderStatePaid)
	}
	if OrderStateShipped != 2 {
		t.Errorf("OrderStateShipped 预期为 2，实际为 %d", OrderStateShipped)
	}
	if OrderStateFinished != 3 {
		t.Errorf("OrderStateFinished 预期为 3，实际为 %d", OrderStateFinished)
	}
}
