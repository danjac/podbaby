package models

import "testing"

func TestPageIfEmpty(t *testing.T) {
	p := NewPaginator(1, 0)
	if p.NumPages != 0 {
		t.Fail()
	}
	if p.Offset != 0 {
		t.Fail()
	}
}

func TestPageIfNotEmpty(t *testing.T) {
	p := NewPaginator(1, 100)
	if p.NumPages != 10 {
		t.Fail()
	}
	if p.Offset != 0 {
		t.Fail()
	}
}

func TestPageOffset(t *testing.T) {
	p := NewPaginator(5, 100)
	if p.NumPages != 10 {
		t.Fail()
	}
	if p.Offset != 40 {
		t.Fail()
	}
}
