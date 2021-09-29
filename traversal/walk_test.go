package traversal_test

import (
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	. "github.com/warpfork/go-wish"

	_ "github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	"github.com/ipld/go-ipld-prime/traversal/selector/builder"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

/* Remember, we've got the following fixtures in scope:
var (
	// baguqeeyexkjwnfy
	leafAlpha, leafAlphaLnk = encode(basicnode.NewString("alpha"))
	// baguqeeyeqvc7t3a
	leafBeta, leafBetaLnk = encode(basicnode.NewString("beta"))
	// baguqeeyezhlahvq
	middleMapNode, middleMapNodeLnk = encode(fluent.MustBuildMap(basicnode.Prototype.Map, 3, func(na fluent.MapAssembler) {
		na.AssembleEntry("foo").AssignBool(true)
		na.AssembleEntry("bar").AssignBool(false)
		na.AssembleEntry("nested").CreateMap(2, func(na fluent.MapAssembler) {
			na.AssembleEntry("alink").AssignLink(leafAlphaLnk)
			na.AssembleEntry("nonlink").AssignString("zoo")
		})
	}))
	// baguqeeyehfkkfwa
	middleListNode, middleListNodeLnk = encode(fluent.MustBuildList(basicnode.Prototype.List, 4, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignLink(leafAlphaLnk)
		na.AssembleValue().AssignLink(leafAlphaLnk)
		na.AssembleValue().AssignLink(leafBetaLnk)
		na.AssembleValue().AssignLink(leafAlphaLnk)
	}))
	// baguqeeyeie4ajfy
	rootNode, rootNodeLnk = encode(fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("plain").AssignString("olde string")
		na.AssembleEntry("linkedString").AssignLink(leafAlphaLnk)
		na.AssembleEntry("linkedMap").AssignLink(middleMapNodeLnk)
		na.AssembleEntry("linkedList").AssignLink(middleListNodeLnk)
	})))
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

func TestWalkBudgets(t *testing.T) {
	ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
	t.Run("node-budget-halts", func(t *testing.T) {
		ss := ssb.ExploreFields(func(efsb builder.ExploreFieldsSpecBuilder) {
			efsb.Insert("foo", ssb.Matcher())
			efsb.Insert("bar", ssb.Matcher())
		})
		s, err := ss.Selector()
		qt.Assert(t, err, qt.Equals, nil)
		var order int
		prog := traversal.Progress{}
		prog.Budget = &traversal.Budget{
			NodeBudget: 2, // should reach root, then "foo", then stop.
		}
		err = prog.WalkMatching(middleMapNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			switch order {
			case 0:
				qt.Assert(t, n, qt.CmpEquals(), basicnode.NewBool(true))
				qt.Assert(t, prog.Path.String(), qt.Equals, "foo")
			}
			order++
			return nil
		})
		qt.Check(t, order, qt.Equals, 1) // because it should've stopped early
		qt.Assert(t, err, qt.Not(qt.Equals), nil)
		qt.Check(t, err.Error(), qt.Equals, `traversal budget exceeded: budget for nodes reached zero while on path "bar"`)
	})
	t.Run("link-budget-halts", func(t *testing.T) {
		ss := ssb.ExploreAll(ssb.Matcher())
		s, err := ss.Selector()
		qt.Assert(t, err, qt.Equals, nil)
		var order int
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
			Budget: &traversal.Budget{
				NodeBudget: 9000,
				LinkBudget: 3,
			},
		}.WalkMatching(middleListNode, s, func(prog traversal.Progress, n datamodel.Node) error {
			switch order {
			case 0:
				qt.Assert(t, n, qt.CmpEquals(), basicnode.NewString("alpha"))
				qt.Assert(t, prog.Path.String(), qt.Equals, "0")
			case 1:
				qt.Assert(t, n, qt.CmpEquals(), basicnode.NewString("alpha"))
				qt.Assert(t, prog.Path.String(), qt.Equals, "1")
			case 2:
				qt.Assert(t, n, qt.CmpEquals(), basicnode.NewString("beta"))
				qt.Assert(t, prog.Path.String(), qt.Equals, "2")
			}
			order++
			return nil
		})
		qt.Check(t, order, qt.Equals, 3)
		qt.Assert(t, err, qt.Not(qt.Equals), nil)
		qt.Check(t, err.Error(), qt.Equals, `traversal budget exceeded: budget for links reached zero while on path "3" (link: "baguqeeyexkjwnfy")`)
	})
}

func TestWalkBlockLoadOrder(t *testing.T) {
	// a more nested root that we can use to test SkipMe as well
	// note that in using `rootNodeLnk` here rather than `rootNode` we're using the
	// dag-json round-trip version which will have different field ordering
	newRootNode, newRootLink := encode(fluent.MustBuildList(basicnode.Prototype.List, 6, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignLink(rootNodeLnk)
		na.AssembleValue().AssignLink(middleListNodeLnk)
		na.AssembleValue().AssignLink(rootNodeLnk)
		na.AssembleValue().AssignLink(middleListNodeLnk)
		na.AssembleValue().AssignLink(rootNodeLnk)
		na.AssembleValue().AssignLink(middleListNodeLnk)
	}))

	linkNames := make(map[datamodel.Link]string)
	linkNames[newRootLink] = "newRootLink"
	linkNames[rootNodeLnk] = "rootNodeLnk"
	linkNames[leafAlphaLnk] = "leafAlphaLnk"
	linkNames[middleMapNodeLnk] = "middleMapNodeLnk"
	linkNames[leafAlphaLnk] = "leafAlphaLnk"
	linkNames[middleListNodeLnk] = "middleListNodeLnk"
	linkNames[leafAlphaLnk] = "leafAlphaLnk"
	linkNames[leafBetaLnk] = "leafBetaLnk"
	for v, n := range linkNames {
		fmt.Printf("n:%v:%v\n", n, v)
	}
	// the links that we expect from the root node, starting _at_ the root node itself
	rootNodeExpectedLinks := []datamodel.Link{
		rootNodeLnk,
		middleListNodeLnk,
		leafAlphaLnk,
		leafAlphaLnk,
		leafBetaLnk,
		leafAlphaLnk,
		middleMapNodeLnk,
		leafAlphaLnk,
		leafAlphaLnk,
	}
	// same thing but for middleListNode
	middleListNodeLinks := []datamodel.Link{
		middleListNodeLnk,
		leafAlphaLnk,
		leafAlphaLnk,
		leafBetaLnk,
		leafAlphaLnk,
	}
	// our newRootNode is a list that contains 3 consecutive links to rootNode
	expectedAllBlocks := make([]datamodel.Link, 3*(len(rootNodeExpectedLinks)+len(middleListNodeLinks)))
	for i := 0; i < 3; i++ {
		copy(expectedAllBlocks[i*len(rootNodeExpectedLinks)+i*len(middleListNodeLinks):], rootNodeExpectedLinks[:])
		copy(expectedAllBlocks[(i+1)*len(rootNodeExpectedLinks)+i*len(middleListNodeLinks):], middleListNodeLinks[:])
	}

	verifySelectorLoads := func(t *testing.T, expected []datamodel.Link, s datamodel.Node, readFn func(lc linking.LinkContext, l datamodel.Link) (io.Reader, error)) {
		var count int
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = func(lc linking.LinkContext, l datamodel.Link) (io.Reader, error) {
			Wish(t, l.String(), ShouldEqual, expected[count].String())
			log.Printf("%v (%v) %s<> %v (%v)\n", l, linkNames[l], strings.Repeat(" ", 17-len(linkNames[l])), expected[count], linkNames[expected[count]])
			count++
			return readFn(lc, l)
		}
		sel, err := selector.CompileSelector(s)
		Wish(t, err, ShouldEqual, nil)
		err = traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.WalkMatching(newRootNode, sel, func(prog traversal.Progress, n datamodel.Node) error {
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, count, ShouldEqual, len(expected))
	}

	t.Run("CommonSelector_MatchAllRecursively", func(t *testing.T) {
		s := selectorparse.CommonSelector_MatchAllRecursively
		verifySelectorLoads(t, expectedAllBlocks, s, (&store).OpenRead)
	})

	t.Run("CommonSelector_ExploreAllRecursively", func(t *testing.T) {
		s := selectorparse.CommonSelector_ExploreAllRecursively
		verifySelectorLoads(t, expectedAllBlocks, s, (&store).OpenRead)
	})

	t.Run("constructed explore-all selector", func(t *testing.T) {
		// used commonly in Filecoin and other places to "visit all blocks in stable order"
		ssb := builder.NewSelectorSpecBuilder(basicnode.Prototype.Any)
		s := ssb.ExploreRecursive(selector.RecursionLimitNone(),
			ssb.ExploreAll(ssb.ExploreRecursiveEdge())).
			Node()
		verifySelectorLoads(t, expectedAllBlocks, s, (&store).OpenRead)
	})

	t.Run("explore-all with duplicate load skips", func(t *testing.T) {
		// when we use SkipMe to skip loading of already visited blocks we expect
		// to see the links show up in Loads but the lack of the links inside rootNode
		// and middleListNode in this list beyond the first set of loads show that
		// the block is not traversed when the SkipMe is received
		expectedSkipMeBlocks := []datamodel.Link{
			rootNodeLnk,
			middleListNodeLnk,
			leafAlphaLnk,
			leafAlphaLnk,
			leafBetaLnk,
			leafAlphaLnk,
			middleMapNodeLnk,
			leafAlphaLnk,
			leafAlphaLnk,
			middleListNodeLnk,
			rootNodeLnk,
			middleListNodeLnk,
			rootNodeLnk,
			middleListNodeLnk,
		}
		s := selectorparse.CommonSelector_ExploreAllRecursively
		visited := make(map[datamodel.Link]bool)
		verifySelectorLoads(t, expectedSkipMeBlocks, s, func(lc linking.LinkContext, l datamodel.Link) (io.Reader, error) {
			log.Printf("load %v [%v]\n", l, visited[l])
			if visited[l] {
				return nil, traversal.SkipMe{}
			}
			visited[l] = true
			return (&store).OpenRead(lc, l)
		})
	})
}
