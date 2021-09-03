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

// ExploreRecursiveWithStop builds a recursive selector node with a stop condition
func ExploreRecursiveWithBypass(limit selector.RecursionLimit, sequence builder.SelectorSpec, stopField string) datamodel.Node {
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

			// Bypass condition
			if stopField != "" {
				cond := fluent.MustBuildMap(basicnode.Prototype__Map{}, 1, func(na fluent.MapAssembler) {
					na.AssembleEntry(string(selector.ConditionMode_HasField)).AssignString(stopField)
				})
				na.AssembleEntry(selector.SelectorKey_Bypass).AssignNode(cond)
			}
		})
	})

}

func TestBypass(t *testing.T) {
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype__Any{})
	t.Run("test ExploreRecursive stopAt with simple node", func(t *testing.T) {
		// Selector that passes through the map
		s, err := selector.CompileSelector(ExploreRecursiveWithBypass(
			selector.RecursionLimitNone(), ssb.ExploreAll(ssb.ExploreRecursiveEdge()),
			"linkedMap"))
		if err != nil {
			t.Fatal(err)
		}
		var order int
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.WalkMatching(rootNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			// fmt.Println("Order", order, prog.Path.String())
			switch order {
			case 0:
				// Root
				Wish(t, prog.Path.String(), ShouldEqual, "")
			case 1:
				Wish(t, prog.Path.String(), ShouldEqual, "plain")
				Wish(t, n, ShouldEqual, basicnode.NewString("olde string"))
			case 2:
				Wish(t, prog.Path.String(), ShouldEqual, "linkedString")
			case 3:
				Wish(t, prog.Path.String(), ShouldEqual, "linkedList")
			// We are starting to traverse the linkedList, we passed through the map already
			case 4:
				Wish(t, prog.Path.String(), ShouldEqual, "linkedList/0")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 8)
	})
}

// mkChain creates a DAG that represent a chain of subDAGs.
// The stopAt condition is extremely appealing for these use cases, as we can
// partially sync a chain using ExploreRecursive without having to sync the
// chain from scratch if we are already partially synced
func mkGitLike() (datamodel.Node, []datamodel.Link) {
	leafAlpha, leafAlphaLnk = encode(basicnode.NewString("alpha"))
	leafBeta, leafBetaLnk = encode(basicnode.NewString("beta"))
	_, mapLnk1 := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
		na.AssembleEntry("foo").AssignBool(true)
		na.AssembleEntry("bar").AssignBool(false)
		na.AssembleEntry("nested").CreateMap(2, func(na fluent.MapAssembler) {
			na.AssembleEntry("alink").AssignLink(leafAlphaLnk)
			na.AssembleEntry("nonlink").AssignString("zoo")
		})
	}))
	_, mapLnk2 := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
		na.AssembleEntry("foo").AssignBool(true)
		na.AssembleEntry("bar").AssignBool(false)
		na.AssembleEntry("nested").CreateMap(2, func(na fluent.MapAssembler) {
			na.AssembleEntry("alink").AssignLink(leafAlphaLnk)
			na.AssembleEntry("nonlink").AssignString("zoo")
		})
	}))
	_, mapLnk3 := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
		na.AssembleEntry("foo").AssignBool(true)
		na.AssembleEntry("bar").AssignBool(false)
		na.AssembleEntry("nested").CreateMap(2, func(na fluent.MapAssembler) {
			na.AssembleEntry("alink").AssignLink(leafAlphaLnk)
			na.AssembleEntry("nonlink").AssignString("zoo")
		})
	}))

	_, ch1Lnk := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("genesis").AssignLink(mapLnk1)
	}))
	_, ch2Lnk := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("ch2").AssignLink(mapLnk2)
		na.AssembleEntry("ch1").AssignLink(ch1Lnk)
	}))
	headNode, headLnk := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("ch3").AssignLink(mapLnk3)
		na.AssembleEntry("ch2").AssignLink(ch2Lnk)
	}))

	return headNode, []datamodel.Link{headLnk, ch2Lnk, ch1Lnk}
}
func TestBypassGitLike(t *testing.T) {
	chainNode, _ := mkGitLike()
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype__Any{})
	t.Run("test ExploreRecursive stopAt with simple node", func(t *testing.T) {
		// Selector that passes through the map
		s, err := selector.CompileSelector(ExploreRecursiveWithBypass(
			selector.RecursionLimitNone(), ssb.ExploreAll(ssb.ExploreRecursiveEdge()),
			"nested"))
		if err != nil {
			t.Fatal(err)
		}
		var order int
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.WalkMatching(chainNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			if prog.Path.String() == "alink" || prog.Path.String() == "nonlink" {
				t.Fatal("we shouldn't have reecursed when seeing nested field")
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, 12)
	})
}
