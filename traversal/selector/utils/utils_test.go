package utils_test

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	_ "github.com/ipld/go-ipld-prime/encoding/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	"github.com/ipld/go-ipld-prime/traversal/selector/utils"
)

var fnb = fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
var (
	middleMapNode = fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
		mb.Insert(knb.CreateString("foo"), vnb.CreateBool(true))
		mb.Insert(knb.CreateString("bar"), vnb.CreateBool(false))
		mb.Insert(knb.CreateString("nested"), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString("nonlink"), vnb.CreateString("zoo"))
		}))
	})
)

func TestCreateNestedSelector(t *testing.T) {
	t.Run("when given empty slice returns a simple positional matcher", func(t *testing.T) {
		t.Skip("Pending -- does not work in current selector context")
		sn, err := utils.CreateNestedPathSelectorNode(nil)
		Require(t, err, ShouldEqual, nil)
		s, err := selector.ParseSelector(sn)
		Require(t, err, ShouldEqual, nil)
		err = traversal.Traverse(fnb.CreateString("x"), s, func(tp traversal.TraversalProgress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, fnb.CreateString("x"))
			Wish(t, tp.Path.String(), ShouldEqual, ipld.Path{}.String())
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("can create a single element path selector", func(t *testing.T) {
		sn, err := utils.CreateNestedPathSelectorNode([]string{"foo"})
		Require(t, err, ShouldEqual, nil)
		s, err := selector.ParseSelector(sn)
		Require(t, err, ShouldEqual, nil)
		var order int
		err = traversal.Traverse(middleMapNode, s, func(tp traversal.TraversalProgress, n ipld.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, fnb.CreateBool(true))
				Wish(t, tp.Path.String(), ShouldEqual, "foo")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 1)
	})
	t.Run("traverse selecting fields recursively should work", func(t *testing.T) {
		sn, err := utils.CreateNestedPathSelectorNode([]string{"nested", "nonlink"})
		Require(t, err, ShouldEqual, nil)
		s, err := selector.ParseSelector(sn)
		Require(t, err, ShouldEqual, nil)
		var order int
		err = traversal.Traverse(middleMapNode, s, func(tp traversal.TraversalProgress, n ipld.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, fnb.CreateString("zoo"))
				Wish(t, tp.Path.String(), ShouldEqual, "nested/nonlink")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 1)
	})
}
