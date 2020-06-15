package crosspkg

import (
	"github.com/ipld/go-ipld-prime/_rsrch/methodsets"
)

type What struct {
	Alpha string
	Beta  string
}

func FlipWhat(x *methodsets.Thing) *What {
	return (*What)(x)
}

type WhatPrivate struct {
	a string
	b string
}

//func FlipWhatPrivate(x *methodsets.ThingPrivate) *What {
//	return (*WhatPrivate)(x)
//}
// NOPE!
// (thank HEAVENS.)
// ./isomorphicCast.go:22:23: cannot convert x (type *methodsets.ThingPrivate) to type *WhatPrivate
