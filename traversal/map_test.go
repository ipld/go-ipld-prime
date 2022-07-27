package traversal_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	nodetests "github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/traversal"
)

func TestMapUpdate(t *testing.T) {
	nestedMapNode, _ := qp.BuildMap(basicnode.Prototype.Map, 3, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "alink", qp.Link(leafAlphaLnk))
		qp.MapEntry(ma, "nonlink", qp.String("zoo"))
	})
	t.Run("read map", func(t *testing.T) {
		m := traversal.NewMap(middleMapNode)
		foo, foundFoo := m.Get("foo")
		qt.Check(t, foundFoo, qt.IsTrue)
		qt.Check(t, foo, nodetests.NodeContentEquals, basicnode.NewBool(true))
		bar, foundBar := m.Get("bar")
		qt.Check(t, foundBar, qt.IsTrue)
		qt.Check(t, bar, nodetests.NodeContentEquals, basicnode.NewBool(false))
		nested, foundNested := m.Get("nested")
		qt.Check(t, foundNested, qt.IsTrue)
		qt.Check(t, nested, nodetests.NodeContentEquals, nestedMapNode)
	})
	t.Run("new map", func(t *testing.T) {
		m1 := traversal.NewMap(nil)
		putFoo := m1.Put("foo", true)
		qt.Check(t, putFoo, qt.IsTrue)
		putBar := m1.Put("bar", false)
		qt.Check(t, putBar, qt.IsTrue)
		m2 := traversal.NewMap(nil)
		putLink := m2.Put("alink", leafAlphaLnk)
		qt.Check(t, putLink, qt.IsTrue)
		putNonlink := m2.Put("nonlink", "zoo")
		qt.Check(t, putNonlink, qt.IsTrue)
		qt.Check(t, m2.(datamodel.Node), nodetests.NodeContentEquals, nestedMapNode)
		putNested := m1.Put("nested", m2)
		qt.Check(t, putNested, qt.IsTrue)
		qt.Check(t, m1.(datamodel.Node), nodetests.NodeContentEquals, middleMapNode)
	})
	t.Run("update map", func(t *testing.T) {
		m := traversal.NewMap(map[string]interface{}{
			"foo": true,
			"bar": false,
		})
		// Check values
		foo, foundFoo := m.Get("foo")
		qt.Check(t, foundFoo, qt.IsTrue)
		qt.Check(t, foo, nodetests.NodeContentEquals, basicnode.NewBool(true))
		bar, foundBar := m.Get("bar")
		qt.Check(t, foundBar, qt.IsTrue)
		qt.Check(t, bar, nodetests.NodeContentEquals, basicnode.NewBool(false))

		// Flip values
		putFoo := m.Put("foo", false)
		qt.Check(t, putFoo, qt.IsTrue)
		putBar := m.Put("bar", true)
		qt.Check(t, putBar, qt.IsTrue)

		// Check flipped values
		foo, foundFoo = m.Get("foo")
		qt.Check(t, foundFoo, qt.IsTrue)
		qt.Check(t, foo, nodetests.NodeContentEquals, basicnode.NewBool(false))
		bar, foundBar = m.Get("bar")
		qt.Check(t, foundBar, qt.IsTrue)
		qt.Check(t, bar, nodetests.NodeContentEquals, basicnode.NewBool(true))
	})
}
