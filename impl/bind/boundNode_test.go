package ipldbind

import (
	"testing"

	"github.com/polydawn/refmt/obj/atlas"
)

// n.b. we can't use any of the standard spec tests yet because we
//  need a MutableNode implementation.
// ... and even then, we're gonna have a fun ride; we need *more* tests
//  for everything that's a MutableNode bound over *anything but*
//   an interface wildcard.  Woofdah.

func Test(t *testing.T) {
	type tFoo struct {
		Bar string
		Baz string
	}
	atl := atlas.MustBuild(
		atlas.BuildEntry(tFoo{}).StructMap().Autogenerate().Complete(),
	)
	n := Bind(tFoo{"one", "two"}, atl)
	_ = n
}
