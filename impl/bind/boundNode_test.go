package ipldbind

import (
	"testing"

	"github.com/polydawn/refmt/obj/atlas"
)

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
