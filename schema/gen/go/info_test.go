package gengo

import (
	"testing"
	"unsafe"
)

// This file contains some short informative tests that were useful to design.
// It's stuff that probably would be just as suited to life in the go playground.

func TestUnderstandingStructMemoryLayout(t *testing.T) {
	t.Skip("This is for human information and not usually necessary to run.")

	// This test informs why the additional bools for optional+nullable fields
	// in structs are grouped together, even though it makes codegen more fiddly.

	// https://go101.org/article/memory-layout.html also has a nice writeup.

	t.Logf("%d\n", unsafe.Sizeof(struct {
		x int32
		y int32
	}{})) // 8
	t.Logf("%d\n", unsafe.Sizeof(struct {
		x int32
		a bool
		y int32
		b bool
	}{})) // 16 ... !  word alignment.
	t.Logf("%d\n", unsafe.Sizeof(struct {
		a bool
		b bool
		x int32
		y int32
	}{})) // 12.  consecutive bools get packed.
	t.Logf("%d\n", unsafe.Sizeof(struct {
		x int32
		y int32
		a bool
		b bool
	}{})) // 12.  consecutive bool packing works anywhere.
	t.Logf("%d\n", unsafe.Sizeof(struct {
		x          int32
		y          int32
		a, b, c, d bool
	}{})) // 12.  bool packing works up to four.
	t.Logf("%d\n", unsafe.Sizeof(struct {
		x             int32
		y             int32
		a, b, c, d, e bool
	}{})) // 16 ... !  bools take a byte; the fifth triggers a new word.
}
