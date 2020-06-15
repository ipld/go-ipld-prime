package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
)

func SpecTestListString(t *testing.T, ns ipld.NodeStyle) {
	t.Run("list<string>, 3 entries", func(t *testing.T) {
		n := fluent.MustBuildList(ns, 3, func(la fluent.ListAssembler) {
			la.AssembleValue().AssignString("one")
			la.AssembleValue().AssignString("two")
			la.AssembleValue().AssignString("three")
		})
		t.Run("reads back out", func(t *testing.T) {
			Wish(t, n.Length(), ShouldEqual, 3)

			v, err := n.LookupIndex(0)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, must.String(v), ShouldEqual, "one")

			v, err = n.LookupIndex(1)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, must.String(v), ShouldEqual, "two")

			v, err = n.LookupIndex(2)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, must.String(v), ShouldEqual, "three")
		})
		t.Run("reads via iteration", func(t *testing.T) {
			itr := n.ListIterator()

			Wish(t, itr.Done(), ShouldEqual, false)
			idx, v, err := itr.Next()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, idx, ShouldEqual, 0)
			Wish(t, must.String(v), ShouldEqual, "one")

			Wish(t, itr.Done(), ShouldEqual, false)
			idx, v, err = itr.Next()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, idx, ShouldEqual, 1)
			Wish(t, must.String(v), ShouldEqual, "two")

			Wish(t, itr.Done(), ShouldEqual, false)
			idx, v, err = itr.Next()
			Wish(t, err, ShouldEqual, nil)
			Wish(t, idx, ShouldEqual, 2)
			Wish(t, must.String(v), ShouldEqual, "three")

			Wish(t, itr.Done(), ShouldEqual, true)
			idx, v, err = itr.Next()
			Wish(t, err, ShouldEqual, ipld.ErrIteratorOverread{})
			Wish(t, idx, ShouldEqual, -1)
			Wish(t, v, ShouldEqual, nil)
		})
	})
}
