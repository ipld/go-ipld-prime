package traversal_test

import (
	"fmt"
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
func ExploreRecursiveWithStop(limit selector.RecursionLimit, sequence builder.SelectorSpec, stopLnk datamodel.Link) datamodel.Node {
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

func TestStopAtLink(t *testing.T) {
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype__Any{})
	t.Run("test ExploreRecursive stopAt with simple node", func(t *testing.T) {
		// Selector that passes through the map
		s, err := selector.CompileSelector(ExploreRecursiveWithStop(
			selector.RecursionLimitNone(), ssb.ExploreAll(ssb.ExploreRecursiveEdge()),
			middleMapNodeLnk))
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
				LinkRevisit:                    true,
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
func mkChain() (datamodel.Node, []datamodel.Link) {
	leafAlpha, leafAlphaLnk = encode(basicnode.NewString("alpha"))
	leafBeta, leafBetaLnk = encode(basicnode.NewString("beta"))
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

	_, ch1Lnk := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("linkedList").AssignLink(middleListNodeLnk)
	}))
	_, ch2Lnk := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("linkedMap").AssignLink(middleMapNodeLnk)
		na.AssembleEntry("ch1").AssignLink(ch1Lnk)
	}))
	_, ch3Lnk := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("linkedString").AssignLink(leafAlphaLnk)
		na.AssembleEntry("ch2").AssignLink(ch2Lnk)
	}))

	headNode, headLnk := encode(fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("plain").AssignString("olde string")
		na.AssembleEntry("ch3").AssignLink(ch3Lnk)
	}))
	return headNode, []datamodel.Link{headLnk, ch3Lnk, ch2Lnk, ch1Lnk}
}

func TestStopInChain(t *testing.T) {
	chainNode, chainLnks := mkChain()
	// Stay in head
	stopAtInChainTest(t, chainNode, chainLnks[1], 2)
	// Get head and following block
	stopAtInChainTest(t, chainNode, chainLnks[2], 4)
	// One more
	stopAtInChainTest(t, chainNode, chainLnks[3], 11)
	// Get the full chain
	stopAtInChainTest(t, chainNode, nil, 17)
}

func stopAtInChainTest(t *testing.T, chainNode datamodel.Node, stopLnk datamodel.Link, numSeen int) {
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype__Any{})
	t.Run(fmt.Sprintf("test ExploreRecursive stopAt in chain with stoplink: %s", stopLnk), func(t *testing.T) {
		s, err := selector.CompileSelector(ExploreRecursiveWithStop(
			selector.RecursionLimitNone(), ssb.ExploreAll(ssb.ExploreRecursiveEdge()),
			stopLnk))
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
				LinkRevisit:                    true,
			},
		}.WalkMatching(chainNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			//fmt.Println("Order", order, prog.Path.String())
			switch order {
			case 0:
				// Root
				Wish(t, prog.Path.String(), ShouldEqual, "")
			case 1:
				Wish(t, prog.Path.String(), ShouldEqual, "plain")
				Wish(t, n, ShouldEqual, basicnode.NewString("olde string"))
			case 2:
				Wish(t, prog.Path.String(), ShouldEqual, "ch3")
			case 3:
				if numSeen > 4 {
					Wish(t, prog.Path.String(), ShouldEqual, "ch3/ch2")
				} else {
					Wish(t, prog.Path.String(), ShouldEqual, "ch3/linkedString")
				}
			case 4:
				if numSeen > 11 {
					Wish(t, prog.Path.String(), ShouldEqual, "ch3/ch2/ch1")
				} else {
					Wish(t, prog.Path.String(), ShouldEqual, "ch3/ch2/linkedMap")
				}
			case 5:
				if numSeen > 11 {
					Wish(t, prog.Path.String(), ShouldEqual, "ch3/ch2/ch1/linkedList")
				} else {
					Wish(t, prog.Path.String(), ShouldEqual, "ch3/ch2/linkedMap/bar")
				}
			case 10:
				if numSeen > 11 {
					Wish(t, prog.Path.String(), ShouldEqual, "ch3/ch2/linkedMap")
				} else {
					Wish(t, prog.Path.String(), ShouldEqual, "ch3/linkedString")
				}
			}
			order++
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, order, ShouldEqual, numSeen)
	})
}
