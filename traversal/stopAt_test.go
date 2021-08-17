package traversal_test

import (
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	_ "github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	"github.com/ipld/go-ipld-prime/traversal/selector/builder"
)

/* Remember, we've got the following fixtures in scope:
var (
	leafAlpha, leafAlphaLnk         = encode(basicnode.NewString("alpha"))
	leafBeta, leafBetaLnk           = encode(basicnode.NewString("beta"))
	middleMapNode, middleMapNodeLnk = encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
		na.AssembleEntry("foo").AssignBool(true)
		na.AssembleEntry("bar").AssignBool(false)
		na.AssembleEntry("nested").CreateMap(2, func(na fluent.MapAssembler) {
			na.AssembleEntry("alink").AssignLink(leafAlphaLnk)
			na.AssembleEntry("nonlink").AssignString("zoo")
		})
	}))
	middleListNode, middleListNodeLnk = encode(fluent.MustBuildList(basicnode.Prototype__List{}, 4, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignLink(leafAlphaLnk)
		na.AssembleValue().AssignLink(leafAlphaLnk)
		na.AssembleValue().AssignLink(leafBetaLnk)
		na.AssembleValue().AssignLink(leafAlphaLnk)
	}))
	rootNode, rootNodeLnk = encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("plain").AssignString("olde string")
		na.AssembleEntry("linkedString").AssignLink(leafAlphaLnk)
		na.AssembleEntry("linkedMap").AssignLink(middleMapNodeLnk)
		na.AssembleEntry("linkedList").AssignLink(middleListNodeLnk)
	}))
)
*/

// covers traverse using a variety of selectors.
// all cases here use one already-loaded Node; no link-loading exercised.

func ExploreRecursiveWithStop(limit selector.RecursionLimit, sequence builder.SelectorSpec, stopLnk ipld.Link) ipld.Node {
	np := basicnode.Prototype__Map{}
	return fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
		// RecursionLimit
		na.AssembleEntry(selector.SelectorKey_ExploreRecursive).CreateMap(3, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
				switch limit.Mode() {
				case selector.RecursionLimit_Depth:
					na.AssembleEntry(selector.SelectorKey_LimitDepth).AssignInt(limit.Depth())
				case selector.RecursionLimit_None:
					na.AssembleEntry(selector.SelectorKey_LimitNone).CreateMap(0, func(na fluent.MapAssembler) {})
				default:
					panic("Unsupported recursion limit type")
				}
			})
			// Sequence
			na.AssembleEntry(selector.SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(selector.SelectorKey_ExploreUnion).CreateList(2, func(na fluent.ListAssembler) {
					na.AssembleValue().AssignNode(fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
						na.AssembleEntry(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
					}))
					na.AssembleValue().AssignNode(sequence.Node())
				})
			})

			// Stop condition
			if stopLnk != nil {
				cond := fluent.MustBuildMap(basicnode.Prototype__Map{}, 1, func(na fluent.MapAssembler) {
					na.AssembleEntry(string(selector.ConditionMode_Link)).AssignLink(stopLnk)
				})
				na.AssembleEntry(selector.SelectorKey_StopAt).AssignNode(cond)
			}
		})
	})

}

func TestStopAt(t *testing.T) {
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype__Any{})
	t.Run("test stop at", func(t *testing.T) {
		s, err := selector.CompileSelector(ExploreRecursiveWithStop(selector.RecursionLimitNone(), ssb.ExploreAll(ssb.ExploreRecursiveEdge()), leafAlphaLnk))

		/* NOTE: ExploreRecursive that can be used for testing purposes
		s, err := ssb.ExploreRecursive(selector.RecursionLimitNone(), ssb.ExploreUnion(
			ssb.Matcher(),
			ssb.ExploreAll(ssb.ExploreRecursiveEdge())),
		).Selector()
		*/
		var order int
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.WalkMatching(rootNode, s, func(prog traversal.Progress, n ipld.Node) error {
			fmt.Println("Order", order, prog.Path.String())
			switch order {
			case 0:
				// Root
				Wish(t, prog.Path.String(), ShouldEqual, "")
			case 1:
				Wish(t, prog.Path.String(), ShouldEqual, "plain")
				Wish(t, n, ShouldEqual, basicnode.NewString("olde string"))
			case 2:
				Wish(t, prog.Path.String(), ShouldEqual, "linkedString")
				/*
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
				*/
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 7)
	})
}
