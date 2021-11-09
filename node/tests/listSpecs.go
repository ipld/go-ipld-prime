package tests

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
)

func SpecTestListString(t *testing.T, np datamodel.NodePrototype) {
	t.Run("list<string>, 3 entries", func(t *testing.T) {
		n := fluent.MustBuildList(np, 3, func(la fluent.ListAssembler) {
			la.AssembleValue().AssignString("one")
			la.AssembleValue().AssignString("two")
			la.AssembleValue().AssignString("three")
		})
		t.Run("reads back out", func(t *testing.T) {
			qt.Check(t, n.Length(), qt.Equals, int64(3))

			v, err := n.LookupByIndex(0)
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, must.String(v), qt.Equals, "one")

			v, err = n.LookupByIndex(1)
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, must.String(v), qt.Equals, "two")

			v, err = n.LookupByIndex(2)
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, must.String(v), qt.Equals, "three")
		})
		t.Run("reads via iteration", func(t *testing.T) {
			itr := n.ListIterator()

			qt.Check(t, itr.Done(), qt.IsFalse)
			idx, v, err := itr.Next()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, idx, qt.Equals, int64(0))
			qt.Check(t, must.String(v), qt.Equals, "one")

			qt.Check(t, itr.Done(), qt.IsFalse)
			idx, v, err = itr.Next()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, idx, qt.Equals, int64(1))
			qt.Check(t, must.String(v), qt.Equals, "two")

			qt.Check(t, itr.Done(), qt.IsFalse)
			idx, v, err = itr.Next()
			qt.Check(t, err, qt.IsNil)
			qt.Check(t, idx, qt.Equals, int64(2))
			qt.Check(t, must.String(v), qt.Equals, "three")

			qt.Check(t, itr.Done(), qt.IsTrue)
			idx, v, err = itr.Next()
			qt.Check(t, err, qt.Equals, datamodel.ErrIteratorOverread{})
			qt.Check(t, idx, qt.Equals, int64(-1))
			qt.Check(t, v, qt.IsNil)
		})
	})
}
