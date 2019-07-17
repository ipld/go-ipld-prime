package traversal_test

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	_ "github.com/ipld/go-ipld-prime/encoding/dagjson"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

/* Remember, we've got the following fixtures in scope:
var (
	leafAlpha, leafAlphaLnk         = encode(fnb.CreateString("alpha"))
	leafBeta, leafBetaLnk           = encode(fnb.CreateString("beta"))
	middleMapNode, middleMapNodeLnk = encode(fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
		mb.Insert(knb.CreateString("foo"), vnb.CreateBool(true))
		mb.Insert(knb.CreateString("bar"), vnb.CreateBool(false))
		mb.Insert(knb.CreateString("nested"), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString("alink"), vnb.CreateLink(leafAlphaLnk))
			mb.Insert(knb.CreateString("nonlink"), vnb.CreateString("zoo"))
		}))
	}))
	middleListNode, middleListNodeLnk = encode(fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
		lb.Append(vnb.CreateLink(leafAlphaLnk))
		lb.Append(vnb.CreateLink(leafAlphaLnk))
		lb.Append(vnb.CreateLink(leafBetaLnk))
		lb.Append(vnb.CreateLink(leafAlphaLnk))
	}))
	rootNode, rootNodeLnk = encode(fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
		mb.Insert(knb.CreateString("plain"), vnb.CreateString("olde string"))
		mb.Insert(knb.CreateString("linkedString"), vnb.CreateLink(leafAlphaLnk))
		mb.Insert(knb.CreateString("linkedMap"), vnb.CreateLink(middleMapNodeLnk))
		mb.Insert(knb.CreateString("linkedList"), vnb.CreateLink(middleListNodeLnk))
	}))
)
*/

// covers traverse using a variety of selectors.
// all cases here use one already-loaded Node; no link-loading exercised.
func TestTraverse(t *testing.T) {
	ssb := selector.NewSelectorSpecBuilder(ipldfree.NodeBuilder())
	t.Run("traverse selecting true should visit the root", func(t *testing.T) {
		err := traversal.Traverse(fnb.CreateString("x"), selector.Matcher{}, func(tp traversal.TraversalProgress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, fnb.CreateString("x"))
			Wish(t, tp.Path.String(), ShouldEqual, ipld.Path{}.String())
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("traverse selecting true should visit only the root and no deeper", func(t *testing.T) {
		err := traversal.Traverse(middleMapNode, selector.Matcher{}, func(tp traversal.TraversalProgress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, middleMapNode)
			Wish(t, tp.Path.String(), ShouldEqual, ipld.Path{}.String())
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("traverse selecting fields should work", func(t *testing.T) {
		ss := ssb.ExploreFields(func(efsb selector.ExploreFieldsSpecBuilder) {
			efsb.Insert("foo", ssb.Matcher())
			efsb.Insert("bar", ssb.Matcher())
		})
		s, err := ss.Selector()
		Require(t, err, ShouldEqual, nil)
		var order int
		err = traversal.Traverse(middleMapNode, s, func(tp traversal.TraversalProgress, n ipld.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, fnb.CreateBool(true))
				Wish(t, tp.Path.String(), ShouldEqual, "foo")
			case 1:
				Wish(t, n, ShouldEqual, fnb.CreateBool(false))
				Wish(t, tp.Path.String(), ShouldEqual, "bar")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 2)
	})
	t.Run("traverse selecting fields recursively should work", func(t *testing.T) {
		ss := ssb.ExploreFields(func(efsb selector.ExploreFieldsSpecBuilder) {
			efsb.Insert("foo", ssb.Matcher())
			efsb.Insert("nested", ssb.ExploreFields(func(efsb selector.ExploreFieldsSpecBuilder) {
				efsb.Insert("nonlink", ssb.Matcher())
			}))
		})
		s, err := ss.Selector()
		Require(t, err, ShouldEqual, nil)
		var order int
		err = traversal.Traverse(middleMapNode, s, func(tp traversal.TraversalProgress, n ipld.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, fnb.CreateBool(true))
				Wish(t, tp.Path.String(), ShouldEqual, "foo")
			case 1:
				Wish(t, n, ShouldEqual, fnb.CreateString("zoo"))
				Wish(t, tp.Path.String(), ShouldEqual, "nested/nonlink")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 2)
	})
}
