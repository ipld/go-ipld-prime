package traversal_test

import (
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp"
	"github.com/ipfs/go-cid"

	_ "github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	nodetests "github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/ipld/go-ipld-prime/traversal"
)

// Do some fixture fabrication.
// We assume all the builders and serialization must Just Work here.

var deepEqualsAllowAllUnexported = qt.CmpEquals(cmp.Exporter(func(reflect.Type) bool { return true }))

var store = memstore.Store{}
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
	// note that using `rootNode` directly will have a different field ordering than
	// the encoded form if you were to load `rootNodeLnk` due to dag-json field
	// reordering on encode, beware the difference for traversal order between
	// created, in-memory nodes and those that have passed through a codec with
	// field ordering rules
	// baguqeeyeie4ajfy
	rootNode, rootNodeLnk = encode(fluent.MustBuildMap(basicnode.Prototype.Map, 4, func(na fluent.MapAssembler) {
		na.AssembleEntry("plain").AssignString("olde string")
		na.AssembleEntry("linkedString").AssignLink(leafAlphaLnk)
		na.AssembleEntry("linkedMap").AssignLink(middleMapNodeLnk)
		na.AssembleEntry("linkedList").AssignLink(middleListNodeLnk)
	}))
)

// encode hardcodes some encoding choices for ease of use in fixture generation;
// just gimme a link and stuff the bytes in a map.
// (also return the node again for convenient assignment.)
func encode(n datamodel.Node) (datamodel.Node, datamodel.Link) {
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x13,
		MhLength: 4,
	}}
	lsys := cidlink.DefaultLinkSystem()
	lsys.SetWriteStorage(&store)

	lnk, err := lsys.Store(linking.LinkContext{}, lp, n)
	if err != nil {
		panic(err)
	}
	return n, lnk
}

// covers Focus used on one already-loaded Node; no link-loading exercised.
func TestFocusSingleTree(t *testing.T) {
	t.Run("empty path on scalar node returns start node", func(t *testing.T) {
		err := traversal.Focus(basicnode.NewString("x"), datamodel.Path{}, func(prog traversal.Progress, n datamodel.Node) error {
			qt.Check(t, n, nodetests.NodeContentEquals, basicnode.NewString("x"))
			qt.Check(t, prog.Path.String(), qt.Equals, datamodel.Path{}.String())
			return nil
		})
		qt.Check(t, err, qt.IsNil)
	})
	t.Run("one step path on map node works", func(t *testing.T) {
		err := traversal.Focus(middleMapNode, datamodel.ParsePath("foo"), func(prog traversal.Progress, n datamodel.Node) error {
			qt.Check(t, n, nodetests.NodeContentEquals, basicnode.NewBool(true))
			qt.Check(t, prog.Path, deepEqualsAllowAllUnexported, datamodel.ParsePath("foo"))
			return nil
		})
		qt.Check(t, err, qt.IsNil)
	})
	t.Run("two step path on map node works", func(t *testing.T) {
		err := traversal.Focus(middleMapNode, datamodel.ParsePath("nested/nonlink"), func(prog traversal.Progress, n datamodel.Node) error {
			qt.Check(t, n, nodetests.NodeContentEquals, basicnode.NewString("zoo"))
			qt.Check(t, prog.Path, deepEqualsAllowAllUnexported, datamodel.ParsePath("nested/nonlink"))
			return nil
		})
		qt.Check(t, err, qt.IsNil)
	})
}

// covers Get used on one already-loaded Node; no link-loading exercised.
// same fixtures as the test for Focus; just has fewer assertions, since Get does no progress tracking.
func TestGetSingleTree(t *testing.T) {
	t.Run("empty path on scalar node returns start node", func(t *testing.T) {
		n, err := traversal.Get(basicnode.NewString("x"), datamodel.Path{})
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n, nodetests.NodeContentEquals, basicnode.NewString("x"))
	})
	t.Run("one step path on map node works", func(t *testing.T) {
		n, err := traversal.Get(middleMapNode, datamodel.ParsePath("foo"))
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n, nodetests.NodeContentEquals, basicnode.NewBool(true))
	})
	t.Run("two step path on map node works", func(t *testing.T) {
		n, err := traversal.Get(middleMapNode, datamodel.ParsePath("nested/nonlink"))
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n, nodetests.NodeContentEquals, basicnode.NewString("zoo"))
	})
}

func TestFocusWithLinkLoading(t *testing.T) {
	t.Run("link traversal with no configured loader should fail", func(t *testing.T) {
		t.Run("terminal link should fail", func(t *testing.T) {
			err := traversal.Focus(middleMapNode, datamodel.ParsePath("nested/alink"), func(prog traversal.Progress, n datamodel.Node) error {
				t.Errorf("should not be reached; no way to load this path")
				return nil
			})
			qt.Check(t, err.Error(), qt.Equals, `error traversing node at "nested/alink": could not load link "`+leafAlphaLnk.String()+`": no LinkTargetNodePrototypeChooser configured`)
		})
		t.Run("mid-path link should fail", func(t *testing.T) {
			err := traversal.Focus(rootNode, datamodel.ParsePath("linkedMap/nested/nonlink"), func(prog traversal.Progress, n datamodel.Node) error {
				t.Errorf("should not be reached; no way to load this path")
				return nil
			})
			qt.Check(t, err.Error(), qt.Equals, `error traversing node at "linkedMap": could not load link "`+middleMapNodeLnk.String()+`": no LinkTargetNodePrototypeChooser configured`)
		})
	})
	t.Run("link traversal with loader should work", func(t *testing.T) {
		lsys := cidlink.DefaultLinkSystem()
		lsys.SetReadStorage(&store)
		err := traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.Focus(rootNode, datamodel.ParsePath("linkedMap/nested/nonlink"), func(prog traversal.Progress, n datamodel.Node) error {
			qt.Check(t, n, nodetests.NodeContentEquals, basicnode.NewString("zoo"))
			qt.Check(t, prog.Path, deepEqualsAllowAllUnexported, datamodel.ParsePath("linkedMap/nested/nonlink"))
			qt.Check(t, prog.LastBlock.Link, deepEqualsAllowAllUnexported, middleMapNodeLnk)
			qt.Check(t, prog.LastBlock.Path, deepEqualsAllowAllUnexported, datamodel.ParsePath("linkedMap"))
			return nil
		})
		qt.Check(t, err, qt.IsNil)
	})
}

func TestGetWithLinkLoading(t *testing.T) {
	t.Run("link traversal with no configured loader should fail", func(t *testing.T) {
		t.Run("terminal link should fail", func(t *testing.T) {
			_, err := traversal.Get(middleMapNode, datamodel.ParsePath("nested/alink"))
			qt.Check(t, err.Error(), qt.Equals, `error traversing node at "nested/alink": could not load link "`+leafAlphaLnk.String()+`": no LinkTargetNodePrototypeChooser configured`)
		})
		t.Run("mid-path link should fail", func(t *testing.T) {
			_, err := traversal.Get(rootNode, datamodel.ParsePath("linkedMap/nested/nonlink"))
			qt.Check(t, err.Error(), qt.Equals, `error traversing node at "linkedMap": could not load link "`+middleMapNodeLnk.String()+`": no LinkTargetNodePrototypeChooser configured`)
		})
	})
	t.Run("link traversal with loader should work", func(t *testing.T) {
		lsys := cidlink.DefaultLinkSystem()
		lsys.SetReadStorage(&store)
		n, err := traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem:                     lsys,
				LinkTargetNodePrototypeChooser: basicnode.Chooser,
			},
		}.Get(rootNode, datamodel.ParsePath("linkedMap/nested/nonlink"))
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n, nodetests.NodeContentEquals, basicnode.NewString("zoo"))
	})
}

func TestFocusedTransform(t *testing.T) {
	t.Run("UpdateMapEntry", func(t *testing.T) {
		n, err := traversal.FocusedTransform(rootNode, datamodel.ParsePath("plain"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "plain")
			qt.Check(t, must.String(prev), qt.Equals, "olde string")
			nb := prev.Prototype().NewBuilder()
			nb.AssignString("new string!")
			return nb.Build(), nil
		}, false)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		// updated value should be there
		qt.Check(t, must.Node(n.LookupByString("plain")), nodetests.NodeContentEquals, basicnode.NewString("new string!"))
		// everything else should be there
		qt.Check(t, must.Node(n.LookupByString("linkedString")), qt.Equals, must.Node(rootNode.LookupByString("linkedString")))
		qt.Check(t, must.Node(n.LookupByString("linkedMap")), qt.Equals, must.Node(rootNode.LookupByString("linkedMap")))
		qt.Check(t, must.Node(n.LookupByString("linkedList")), qt.Equals, must.Node(rootNode.LookupByString("linkedList")))
		// everything should still be in the same order
		qt.Check(t, keys(n), qt.DeepEquals, []string{"plain", "linkedString", "linkedMap", "linkedList"})
	})
	t.Run("UpdateDeeperMap", func(t *testing.T) {
		n, err := traversal.FocusedTransform(middleMapNode, datamodel.ParsePath("nested/alink"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "nested/alink")
			qt.Check(t, prev, nodetests.NodeContentEquals, basicnode.NewLink(leafAlphaLnk))
			return basicnode.NewString("new string!"), nil
		}, false)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		// updated value should be there
		qt.Check(t, must.Node(must.Node(n.LookupByString("nested")).LookupByString("alink")), nodetests.NodeContentEquals, basicnode.NewString("new string!"))
		// everything else in the parent map should should be there!
		qt.Check(t, must.Node(n.LookupByString("foo")), qt.Equals, must.Node(middleMapNode.LookupByString("foo")))
		qt.Check(t, must.Node(n.LookupByString("bar")), qt.Equals, must.Node(middleMapNode.LookupByString("bar")))
		// everything should still be in the same order
		qt.Check(t, keys(n), qt.DeepEquals, []string{"foo", "bar", "nested"})
	})
	t.Run("AppendIfNotExists", func(t *testing.T) {
		n, err := traversal.FocusedTransform(rootNode, datamodel.ParsePath("newpart"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "newpart")
			qt.Check(t, prev, qt.IsNil) // REVIEW: should datamodel.Absent be used here?  I lean towards "no" but am unsure what's least surprising here.
			// An interesting thing to note about inserting a value this way is that you have no `prev.Prototype().NewBuilder()` to use if you wanted to.
			//  But if that's an issue, then what you do is a focus or walk (transforming or not) to the parent node, get its child prototypes, and go from there.
			return basicnode.NewString("new string!"), nil
		}, false)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		// updated value should be there
		qt.Check(t, must.Node(n.LookupByString("newpart")), nodetests.NodeContentEquals, basicnode.NewString("new string!"))
		// everything should still be in the same order... with the new entry at the end.
		qt.Check(t, keys(n), qt.DeepEquals, []string{"plain", "linkedString", "linkedMap", "linkedList", "newpart"})
	})
	t.Run("CreateParents", func(t *testing.T) {
		n, err := traversal.FocusedTransform(rootNode, datamodel.ParsePath("newsection/newpart"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "newsection/newpart")
			qt.Check(t, prev, qt.IsNil) // REVIEW: should datamodel.Absent be used here?  I lean towards "no" but am unsure what's least surprising here.
			return basicnode.NewString("new string!"), nil
		}, true)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		// a new map node in the middle should've been created
		n2 := must.Node(n.LookupByString("newsection"))
		qt.Check(t, n2.Kind(), qt.Equals, datamodel.Kind_Map)
		// updated value should in there
		qt.Check(t, must.Node(n2.LookupByString("newpart")), nodetests.NodeContentEquals, basicnode.NewString("new string!"))
		// everything in the root map should still be in the same order... with the new entry at the end.
		qt.Check(t, keys(n), qt.DeepEquals, []string{"plain", "linkedString", "linkedMap", "linkedList", "newsection"})
		// and the created intermediate map of course has just one entry.
		qt.Check(t, keys(n2), qt.DeepEquals, []string{"newpart"})
	})
	t.Run("CreateParentsRequiresPermission", func(t *testing.T) {
		_, err := traversal.FocusedTransform(rootNode, datamodel.ParsePath("newsection/newpart"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, true, qt.IsFalse) // ought not be reached
			return nil, nil
		}, false)
		qt.Check(t, err.Error(), qt.Equals, "transform: parent position at \"newsection\" did not exist (and createParents was false)")
	})
	t.Run("UpdateListEntry", func(t *testing.T) {
		n, err := traversal.FocusedTransform(middleListNode, datamodel.ParsePath("2"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "2")
			qt.Check(t, prev, nodetests.NodeContentEquals, basicnode.NewLink(leafBetaLnk))
			return basicnode.NewString("new string!"), nil
		}, false)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_List)
		// updated value should be there
		qt.Check(t, must.Node(n.LookupByIndex(2)), nodetests.NodeContentEquals, basicnode.NewString("new string!"))
		// everything else should be there
		qt.Check(t, n.Length(), qt.Equals, int64(4))
		qt.Check(t, must.Node(n.LookupByIndex(0)), nodetests.NodeContentEquals, basicnode.NewLink(leafAlphaLnk))
		qt.Check(t, must.Node(n.LookupByIndex(1)), nodetests.NodeContentEquals, basicnode.NewLink(leafAlphaLnk))
		qt.Check(t, must.Node(n.LookupByIndex(3)), nodetests.NodeContentEquals, basicnode.NewLink(leafAlphaLnk))
	})
	t.Run("AppendToList", func(t *testing.T) {
		n, err := traversal.FocusedTransform(middleListNode, datamodel.ParsePath("-"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "4")
			qt.Check(t, prev, qt.IsNil)
			return basicnode.NewString("new string!"), nil
		}, false)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_List)
		// updated value should be there
		qt.Check(t, must.Node(n.LookupByIndex(4)), nodetests.NodeContentEquals, basicnode.NewString("new string!"))
		// everything else should be there
		qt.Check(t, n.Length(), qt.Equals, int64(5))
	})
	t.Run("ListBounds", func(t *testing.T) {
		_, err := traversal.FocusedTransform(middleListNode, datamodel.ParsePath("4"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, true, qt.IsFalse) // ought not be reached
			return nil, nil
		}, false)
		qt.Check(t, err, qt.ErrorMatches, "transform: cannot navigate path segment \"4\" at \"\" because it is beyond the list bounds")
	})
	t.Run("ReplaceRoot", func(t *testing.T) { // a fairly degenerate case and no reason to do this, but should work.
		n, err := traversal.FocusedTransform(middleListNode, datamodel.ParsePath(""), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "")
			qt.Check(t, prev, nodetests.NodeContentEquals, middleListNode)
			nb := basicnode.Prototype.Any.NewBuilder()
			la, _ := nb.BeginList(0)
			la.Finish()
			return nb.Build(), nil
		}, false)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_List)
		qt.Check(t, n.Length(), qt.Equals, int64(0))
	})
}

func TestFocusedTransformWithLinks(t *testing.T) {
	var store2 = memstore.Store{}
	lsys := cidlink.DefaultLinkSystem()
	lsys.SetReadStorage(&store)
	lsys.SetWriteStorage(&store2)
	cfg := traversal.Config{
		LinkSystem:                     lsys,
		LinkTargetNodePrototypeChooser: basicnode.Chooser,
	}
	t.Run("UpdateMapBeyondLink", func(t *testing.T) {
		n, err := traversal.Progress{
			Cfg: &cfg,
		}.FocusedTransform(rootNode, datamodel.ParsePath("linkedMap/nested/nonlink"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "linkedMap/nested/nonlink")
			qt.Check(t, must.String(prev), qt.Equals, "zoo")
			qt.Check(t, progress.LastBlock.Path.String(), qt.Equals, "linkedMap")
			qt.Check(t, progress.LastBlock.Link.String(), qt.Equals, "baguqeeyezhlahvq")
			nb := prev.Prototype().NewBuilder()
			nb.AssignString("new string!")
			return nb.Build(), nil
		}, false)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		// there should be a new object in our new storage!
		qt.Check(t, store2.Bag, qt.HasLen, 1)
		// cleanup for next test
		store2 = memstore.Store{}
	})
	t.Run("UpdateNotBeyondLink", func(t *testing.T) {
		// This is replacing a link with a non-link.  Doing so shouldn't hit storage.
		n, err := traversal.Progress{
			Cfg: &cfg,
		}.FocusedTransform(rootNode, datamodel.ParsePath("linkedMap"), func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			qt.Check(t, progress.Path.String(), qt.Equals, "linkedMap")
			nb := prev.Prototype().NewBuilder()
			nb.AssignString("new string!")
			return nb.Build(), nil
		}, false)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		// there should be no new objects in our new storage!
		qt.Check(t, store2.Bag, qt.HasLen, 0)
		// cleanup for next test
		store2 = memstore.Store{}
	})

	// link traverse to scalar // this is unspecifiable using the current path syntax!  you'll just end up replacing the link with the scalar!
}

func keys(n datamodel.Node) []string {
	v := make([]string, 0, n.Length())
	for itr := n.MapIterator(); !itr.Done(); {
		k, _, _ := itr.Next()
		v = append(v, must.String(k))
	}
	return v
}
