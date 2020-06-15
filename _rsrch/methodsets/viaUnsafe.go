package methodsets

import (
	"unsafe"
)

type Thing2ViaUnsafe struct {
	Alpha string
	Beta  string
}

func FlipUnsafe(x *Thing) *Thing2ViaUnsafe {
	return (*Thing2ViaUnsafe)(unsafe.Pointer(x))
}

func UnflipUnsafe(x *Thing2ViaUnsafe) *Thing {
	return (*Thing)(unsafe.Pointer(x))
}

func (x *Thing2ViaUnsafe) Pow() {
	x.Alpha = "unsafe"
}
