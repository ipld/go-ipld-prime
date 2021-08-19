package traversal_test

import (
	"testing"

	. "github.com/warpfork/go-wish"

	_ "github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	"github.com/ipld/go-ipld-prime/traversal/selector/builder"
)

/* Remember, we've got the following fixtures in scope:
var (
	leafAlpha, leafAlphaLnk         = encode(basicnode.NewString("alpha"))
	leafBeta, leafBetaLnk           = encode(basicnode.NewString("beta"))
	middleMapNode, middleMapNodeLnk = encode(fluent.MustBuildMap(basicnode.Prototype.Map, 3, func(na fluent.MapAssembler) {
		na.AssembleEntry("foo").AssignBool(true)
		na.AssembleEntry("bar").AssignBool(false)
		na.AssembleEntry("nested").CreateMap(2, func(na fluent.MapAssembler) {
			na.AssembleEntry("alink").AssignLink(leafAlphaLnk)
			na.AssembleEntry("nonlink").AssignString("zoo")
		})
	}))
	middleListNode, middleListNodeLnk = encode(fluent.MustBuildList(basicnode.Prototype.List, 4, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignLink(leafAlphaLnk)
		na.AssembleValue().AssignLink(leafAlphaLnk)
		na.AssembleValue().AssignLink(leafBetaLnk)
		na.AssembleValue().AssignLink(leafAlphaLnk)
	}))
	rootNode, rootNodeLnk = encode(fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("plain").AssignString("olde string")
		na.AssembleEntry("linkedString").AssignLink(leafAlphaLnk)
		na.AssembleEntry("linkedMap").AssignLink(middleMapNodeLnk)
		na.AssembleEntry("linkedList").AssignLink(middleListNodeLnk)
	}))
)
*/

// covers traverse using a variety of selectors.
// all cases here use one already-loaded Node; no link-loading exercised.

func TestWalkMatching(t *testing.T) {
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
	t.Run("traverse selecting true should visit the root", func(t *testing.T) {
		err := traversal.WalkMatching(basicnode.NewString("x"), selector.Matcher{}, func(prog traversal.Progress, n datamodel.Node) error {
			Wish(t, n, ShouldEqual, basicnode.NewString("x"))
			Wish(t, prog.Path.String(), ShouldEqual, datamodel.Path{}.String())
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("traverse selecting true should visit only the root and no deeper", func(t *testing.T) {
		err := traversal.WalkMatching(middleMapNode, selector.Matcher{}, func(prog traversal.Progress, n datamodel.Node) error {
			Wish(t, n, ShouldEqual, middleMapNode)
			Wish(t, prog.Path.String(), ShouldEqual, datamodel.Path{}.String())
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("traverse selecting fields should work", func(t *testing.T) {
		ss := ssb.ExploreFields(func(efsb builder.ExploreFieldsSpecBuilder) {
			efsb.Insert("foo", ssb.Matcher())
			efsb.Insert("bar", ssb.Matcher())
		})
		s, err := ss.Selector()
		Require(t, err, ShouldEqual, nil)
		var order int
		err = traversal.WalkMatching(middleMapNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, basicnode.NewBool(true))
				Wish(t, prog.Path.String(), ShouldEqual, "foo")
			case 1:
				Wish(t, n, ShouldEqual, basicnode.NewBool(false))
				Wish(t, prog.Path.String(), ShouldEqual, "bar")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 2)
	})
	t.Run("traverse selecting fields recursively should work", func(t *testing.T) {
		ss := ssb.ExploreFields(func(efsb builder.ExploreFieldsSpecBuilder) {
			efsb.Insert("foo", ssb.Matcher())
			efsb.Insert("nested", ssb.ExploreFields(func(efsb builder.ExploreFieldsSpecBuilder) {
				efsb.Insert("nonlink", ssb.Matcher())
			}))
		})
		s, err := ss.Selector()
		Require(t, err, ShouldEqual, nil)
		var order int
		err = traversal.WalkMatching(middleMapNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, basicnode.NewBool(true))
				Wish(t, prog.Path.String(), ShouldEqual, "foo")
			case 1:
				Wish(t, n, ShouldEqual, basicnode.NewString("zoo"))
				Wish(t, prog.Path.String(), ShouldEqual, "nested/nonlink")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 2)
	})
	t.Run("traversing across nodes should work", func(t *testing.T) {
		ss := ssb.ExploreRecursive(selector.RecursionLimitDepth(3), ssb.ExploreUnion(
			ssb.Matcher(),
			ssb.ExploreAll(ssb.ExploreRecursiveEdge()),
		))
		s, err := ss.Selector()
		var order int
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.WalkMatching(middleMapNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, middleMapNode)
				Wish(t, prog.Path.String(), ShouldEqual, "")
			case 1:
				Wish(t, n, ShouldEqual, basicnode.NewBool(true))
				Wish(t, prog.Path.String(), ShouldEqual, "foo")
			case 2:
				Wish(t, n, ShouldEqual, basicnode.NewBool(false))
				Wish(t, prog.Path.String(), ShouldEqual, "bar")
			case 3:
				Wish(t, n, ShouldEqual, fluent.MustBuildMap(basicnode.Prototype.Map, 2, func(na fluent.MapAssembler) {
					na.AssembleEntry("alink").AssignLink(leafAlphaLnk)
					na.AssembleEntry("nonlink").AssignString("zoo")
				}))
				Wish(t, prog.Path.String(), ShouldEqual, "nested")
			case 4:
				Wish(t, n, ShouldEqual, basicnode.NewString("alpha"))
				Wish(t, prog.Path.String(), ShouldEqual, "nested/alink")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "nested/alink")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafAlphaLnk.String())

			case 5:
				Wish(t, n, ShouldEqual, basicnode.NewString("zoo"))
				Wish(t, prog.Path.String(), ShouldEqual, "nested/nonlink")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 6)
	})
	t.Run("traversing lists should work", func(t *testing.T) {
		ss := ssb.ExploreRange(0, 3, ssb.Matcher())
		s, err := ss.Selector()
		var order int
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.WalkMatching(middleListNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, basicnode.NewString("alpha"))
				Wish(t, prog.Path.String(), ShouldEqual, "0")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "0")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafAlphaLnk.String())
			case 1:
				Wish(t, n, ShouldEqual, basicnode.NewString("alpha"))
				Wish(t, prog.Path.String(), ShouldEqual, "1")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "1")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafAlphaLnk.String())
			case 2:
				Wish(t, n, ShouldEqual, basicnode.NewString("beta"))
				Wish(t, prog.Path.String(), ShouldEqual, "2")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "2")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafBetaLnk.String())
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 3)
	})
	t.Run("multiple layers of link traversal should work", func(t *testing.T) {
		ss := ssb.ExploreFields(func(efsb builder.ExploreFieldsSpecBuilder) {
			efsb.Insert("linkedList", ssb.ExploreAll(ssb.Matcher()))
			efsb.Insert("linkedMap", ssb.ExploreRecursive(selector.RecursionLimitDepth(3), ssb.ExploreFields(func(efsb builder.ExploreFieldsSpecBuilder) {
				efsb.Insert("foo", ssb.Matcher())
				efsb.Insert("nonlink", ssb.Matcher())
				efsb.Insert("alink", ssb.Matcher())
				efsb.Insert("nested", ssb.ExploreRecursiveEdge())
			})))
		})
		s, err := ss.Selector()
		var order int
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.WalkMatching(rootNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			switch order {
			case 0:
				Wish(t, n, ShouldEqual, basicnode.NewString("alpha"))
				Wish(t, prog.Path.String(), ShouldEqual, "linkedList/0")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "linkedList/0")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafAlphaLnk.String())
			case 1:
				Wish(t, n, ShouldEqual, basicnode.NewString("alpha"))
				Wish(t, prog.Path.String(), ShouldEqual, "linkedList/1")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "linkedList/1")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafAlphaLnk.String())
			case 2:
				Wish(t, n, ShouldEqual, basicnode.NewString("beta"))
				Wish(t, prog.Path.String(), ShouldEqual, "linkedList/2")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "linkedList/2")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafBetaLnk.String())
			case 3:
				Wish(t, n, ShouldEqual, basicnode.NewString("alpha"))
				Wish(t, prog.Path.String(), ShouldEqual, "linkedList/3")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "linkedList/3")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafAlphaLnk.String())
			case 4:
				Wish(t, n, ShouldEqual, basicnode.NewBool(true))
				Wish(t, prog.Path.String(), ShouldEqual, "linkedMap/foo")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "linkedMap")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, middleMapNodeLnk.String())
			case 5:
				Wish(t, n, ShouldEqual, basicnode.NewString("zoo"))
				Wish(t, prog.Path.String(), ShouldEqual, "linkedMap/nested/nonlink")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "linkedMap")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, middleMapNodeLnk.String())
			case 6:
				Wish(t, n, ShouldEqual, basicnode.NewString("alpha"))
				Wish(t, prog.Path.String(), ShouldEqual, "linkedMap/nested/alink")
				Wish(t, prog.LastBlock.Path.String(), ShouldEqual, "linkedMap/nested/alink")
				Wish(t, prog.LastBlock.Link.String(), ShouldEqual, leafAlphaLnk.String())
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 7)
	})
}
