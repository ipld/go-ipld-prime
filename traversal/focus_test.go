package traversal_test

import (
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	_ "github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/must"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/storage"
	"github.com/ipld/go-ipld-prime/traversal"
)

// Do some fixture fabrication.
// We assume all the builders and serialization must Just Work here.

var store = storage.Memory{}
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

// encode hardcodes some encoding choices for ease of use in fixture generation;
// just gimme a link and stuff the bytes in a map.
// (also return the node again for convenient assignment.)
func encode(n ipld.Node) (ipld.Node, ipld.Link) {
	lp := cidlink.LinkPrototype{cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x13,
		MhLength: 4,
	}}
	lsys := cidlink.DefaultLinkSystem()
	lsys.StorageWriteOpener = (&store).OpenWrite

	lnk, err := lsys.Store(ipld.LinkContext{}, lp, n)
	if err != nil {
		panic(err)
	}
	return n, lnk
}

// covers Focus used on one already-loaded Node; no link-loading exercised.
func TestFocusSingleTree(t *testing.T) {
	t.Run("empty path on scalar node returns start node", func(t *testing.T) {
		err := traversal.Focus(basicnode.NewString("x"), ipld.Path{}, func(prog traversal.Progress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, basicnode.NewString("x"))
			Wish(t, prog.Path.String(), ShouldEqual, ipld.Path{}.String())
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("one step path on map node works", func(t *testing.T) {
		err := traversal.Focus(middleMapNode, ipld.ParsePath("foo"), func(prog traversal.Progress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, basicnode.NewBool(true))
			Wish(t, prog.Path, ShouldEqual, ipld.ParsePath("foo"))
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("two step path on map node works", func(t *testing.T) {
		err := traversal.Focus(middleMapNode, ipld.ParsePath("nested/nonlink"), func(prog traversal.Progress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, basicnode.NewString("zoo"))
			Wish(t, prog.Path, ShouldEqual, ipld.ParsePath("nested/nonlink"))
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
}

// covers Get used on one already-loaded Node; no link-loading exercised.
// same fixtures as the test for Focus; just has fewer assertions, since Get does no progress tracking.
func TestGetSingleTree(t *testing.T) {
	t.Run("empty path on scalar node returns start node", func(t *testing.T) {
		n, err := traversal.Get(basicnode.NewString("x"), ipld.Path{})
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n, ShouldEqual, basicnode.NewString("x"))
	})
	t.Run("one step path on map node works", func(t *testing.T) {
		n, err := traversal.Get(middleMapNode, ipld.ParsePath("foo"))
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n, ShouldEqual, basicnode.NewBool(true))
	})
	t.Run("two step path on map node works", func(t *testing.T) {
		n, err := traversal.Get(middleMapNode, ipld.ParsePath("nested/nonlink"))
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n, ShouldEqual, basicnode.NewString("zoo"))
	})
}

func TestFocusWithLinkLoading(t *testing.T) {
	t.Run("link traversal with no configured loader should fail", func(t *testing.T) {
		t.Run("terminal link should fail", func(t *testing.T) {
			err := traversal.Focus(middleMapNode, ipld.ParsePath("nested/alink"), func(prog traversal.Progress, n ipld.Node) error {
				t.Errorf("should not be reached; no way to load this path")
				return nil
			})
			Wish(t, err.Error(), ShouldEqual, `error traversing node at "nested/alink": could not load link "`+leafAlphaLnk.String()+`": no LinkTargetNodePrototypeChooser configured`)
		})
		t.Run("mid-path link should fail", func(t *testing.T) {
			err := traversal.Focus(rootNode, ipld.ParsePath("linkedMap/nested/nonlink"), func(prog traversal.Progress, n ipld.Node) error {
				t.Errorf("should not be reached; no way to load this path")
				return nil
			})
			Wish(t, err.Error(), ShouldEqual, `error traversing node at "linkedMap": could not load link "`+middleMapNodeLnk.String()+`": no LinkTargetNodePrototypeChooser configured`)
		})
	})
	t.Run("link traversal with loader should work", func(t *testing.T) {
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		err := traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem: lsys,
				LinkTargetNodePrototypeChooser: func(_ ipld.Link, _ ipld.LinkContext) (ipld.NodePrototype, error) {
					return basicnode.Prototype__Any{}, nil
				},
			},
		}.Focus(rootNode, ipld.ParsePath("linkedMap/nested/nonlink"), func(prog traversal.Progress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, basicnode.NewString("zoo"))
			Wish(t, prog.Path, ShouldEqual, ipld.ParsePath("linkedMap/nested/nonlink"))
			Wish(t, prog.LastBlock.Link, ShouldEqual, middleMapNodeLnk)
			Wish(t, prog.LastBlock.Path, ShouldEqual, ipld.ParsePath("linkedMap"))
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
}

func TestGetWithLinkLoading(t *testing.T) {
	t.Run("link traversal with no configured loader should fail", func(t *testing.T) {
		t.Run("terminal link should fail", func(t *testing.T) {
			_, err := traversal.Get(middleMapNode, ipld.ParsePath("nested/alink"))
			Wish(t, err.Error(), ShouldEqual, `error traversing node at "nested/alink": could not load link "`+leafAlphaLnk.String()+`": no LinkTargetNodePrototypeChooser configured`)
		})
		t.Run("mid-path link should fail", func(t *testing.T) {
			_, err := traversal.Get(rootNode, ipld.ParsePath("linkedMap/nested/nonlink"))
			Wish(t, err.Error(), ShouldEqual, `error traversing node at "linkedMap": could not load link "`+middleMapNodeLnk.String()+`": no LinkTargetNodePrototypeChooser configured`)
		})
	})
	t.Run("link traversal with loader should work", func(t *testing.T) {
		lsys := cidlink.DefaultLinkSystem()
		lsys.StorageReadOpener = (&store).OpenRead
		n, err := traversal.Progress{
			Cfg: &traversal.Config{
				LinkSystem: lsys,
				LinkTargetNodePrototypeChooser: func(_ ipld.Link, _ ipld.LinkContext) (ipld.NodePrototype, error) {
					return basicnode.Prototype__Any{}, nil
				},
			},
		}.Get(rootNode, ipld.ParsePath("linkedMap/nested/nonlink"))
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n, ShouldEqual, basicnode.NewString("zoo"))
	})
}

func TestFocusedTransform(t *testing.T) {
	t.Run("UpdateMapEntry", func(t *testing.T) {
		n, err := traversal.FocusedTransform(rootNode, ipld.ParsePath("plain"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "plain")
			Wish(t, must.String(prev), ShouldEqual, "olde string")
			nb := prev.Prototype().NewBuilder()
			nb.AssignString("new string!")
			return nb.Build(), nil
		}, false)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
		// updated value should be there
		Wish(t, must.Node(n.LookupByString("plain")), ShouldEqual, basicnode.NewString("new string!"))
		// everything else should be there
		Wish(t, must.Node(n.LookupByString("linkedString")), ShouldEqual, must.Node(rootNode.LookupByString("linkedString")))
		Wish(t, must.Node(n.LookupByString("linkedMap")), ShouldEqual, must.Node(rootNode.LookupByString("linkedMap")))
		Wish(t, must.Node(n.LookupByString("linkedList")), ShouldEqual, must.Node(rootNode.LookupByString("linkedList")))
		// everything should still be in the same order
		Wish(t, keys(n), ShouldEqual, []string{"plain", "linkedString", "linkedMap", "linkedList"})
	})
	t.Run("UpdateDeeperMap", func(t *testing.T) {
		n, err := traversal.FocusedTransform(middleMapNode, ipld.ParsePath("nested/alink"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "nested/alink")
			Wish(t, prev, ShouldEqual, basicnode.NewLink(leafAlphaLnk))
			return basicnode.NewString("new string!"), nil
		}, false)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
		// updated value should be there
		Wish(t, must.Node(must.Node(n.LookupByString("nested")).LookupByString("alink")), ShouldEqual, basicnode.NewString("new string!"))
		// everything else in the parent map should should be there!
		Wish(t, must.Node(n.LookupByString("foo")), ShouldEqual, must.Node(middleMapNode.LookupByString("foo")))
		Wish(t, must.Node(n.LookupByString("bar")), ShouldEqual, must.Node(middleMapNode.LookupByString("bar")))
		// everything should still be in the same order
		Wish(t, keys(n), ShouldEqual, []string{"foo", "bar", "nested"})
	})
	t.Run("AppendIfNotExists", func(t *testing.T) {
		n, err := traversal.FocusedTransform(rootNode, ipld.ParsePath("newpart"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "newpart")
			Wish(t, prev, ShouldEqual, nil) // REVIEW: should ipld.Absent be used here?  I lean towards "no" but am unsure what's least surprising here.
			// An interesting thing to note about inserting a value this way is that you have no `prev.Prototype().NewBuilder()` to use if you wanted to.
			//  But if that's an issue, then what you do is a focus or walk (transforming or not) to the parent node, get its child prototypes, and go from there.
			return basicnode.NewString("new string!"), nil
		}, false)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
		// updated value should be there
		Wish(t, must.Node(n.LookupByString("newpart")), ShouldEqual, basicnode.NewString("new string!"))
		// everything should still be in the same order... with the new entry at the end.
		Wish(t, keys(n), ShouldEqual, []string{"plain", "linkedString", "linkedMap", "linkedList", "newpart"})
	})
	t.Run("CreateParents", func(t *testing.T) {
		n, err := traversal.FocusedTransform(rootNode, ipld.ParsePath("newsection/newpart"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "newsection/newpart")
			Wish(t, prev, ShouldEqual, nil) // REVIEW: should ipld.Absent be used here?  I lean towards "no" but am unsure what's least surprising here.
			return basicnode.NewString("new string!"), nil
		}, true)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
		// a new map node in the middle should've been created
		n2 := must.Node(n.LookupByString("newsection"))
		Wish(t, n2.Kind(), ShouldEqual, ipld.Kind_Map)
		// updated value should in there
		Wish(t, must.Node(n2.LookupByString("newpart")), ShouldEqual, basicnode.NewString("new string!"))
		// everything in the root map should still be in the same order... with the new entry at the end.
		Wish(t, keys(n), ShouldEqual, []string{"plain", "linkedString", "linkedMap", "linkedList", "newsection"})
		// and the created intermediate map of course has just one entry.
		Wish(t, keys(n2), ShouldEqual, []string{"newpart"})
	})
	t.Run("CreateParentsRequiresPermission", func(t *testing.T) {
		_, err := traversal.FocusedTransform(rootNode, ipld.ParsePath("newsection/newpart"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, true, ShouldEqual, false) // ought not be reached
			return nil, nil
		}, false)
		Wish(t, err, ShouldEqual, fmt.Errorf("transform: parent position at \"newsection\" did not exist (and createParents was false)"))
	})
	t.Run("UpdateListEntry", func(t *testing.T) {
		n, err := traversal.FocusedTransform(middleListNode, ipld.ParsePath("2"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "2")
			Wish(t, prev, ShouldEqual, basicnode.NewLink(leafBetaLnk))
			return basicnode.NewString("new string!"), nil
		}, false)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_List)
		// updated value should be there
		Wish(t, must.Node(n.LookupByIndex(2)), ShouldEqual, basicnode.NewString("new string!"))
		// everything else should be there
		Wish(t, n.Length(), ShouldEqual, int64(4))
		Wish(t, must.Node(n.LookupByIndex(0)), ShouldEqual, basicnode.NewLink(leafAlphaLnk))
		Wish(t, must.Node(n.LookupByIndex(1)), ShouldEqual, basicnode.NewLink(leafAlphaLnk))
		Wish(t, must.Node(n.LookupByIndex(3)), ShouldEqual, basicnode.NewLink(leafAlphaLnk))
	})
	t.Run("AppendToList", func(t *testing.T) {
		n, err := traversal.FocusedTransform(middleListNode, ipld.ParsePath("-"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "4")
			Wish(t, prev, ShouldEqual, nil)
			return basicnode.NewString("new string!"), nil
		}, false)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_List)
		// updated value should be there
		Wish(t, must.Node(n.LookupByIndex(4)), ShouldEqual, basicnode.NewString("new string!"))
		// everything else should be there
		Wish(t, n.Length(), ShouldEqual, int64(5))
	})
	t.Run("ListBounds", func(t *testing.T) {
		_, err := traversal.FocusedTransform(middleListNode, ipld.ParsePath("4"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, true, ShouldEqual, false) // ought not be reached
			return nil, nil
		}, false)
		Wish(t, err, ShouldEqual, fmt.Errorf("transform: cannot navigate path segment \"4\" at \"\" because it is beyond the list bounds"))
	})
	t.Run("ReplaceRoot", func(t *testing.T) { // a fairly degenerate case and no reason to do this, but should work.
		n, err := traversal.FocusedTransform(middleListNode, ipld.ParsePath(""), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "")
			Wish(t, prev, ShouldEqual, middleListNode)
			nb := basicnode.Prototype.Any.NewBuilder()
			la, _ := nb.BeginList(0)
			la.Finish()
			return nb.Build(), nil
		}, false)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_List)
		Wish(t, n.Length(), ShouldEqual, int64(0))
	})
}

func TestFocusedTransformWithLinks(t *testing.T) {
	var store2 = storage.Memory{}
	lsys := cidlink.DefaultLinkSystem()
	lsys.StorageReadOpener = (&store).OpenRead
	lsys.StorageWriteOpener = (&store2).OpenWrite
	cfg := traversal.Config{
		LinkSystem: lsys,
		LinkTargetNodePrototypeChooser: func(_ ipld.Link, _ ipld.LinkContext) (ipld.NodePrototype, error) {
			return basicnode.Prototype.Any, nil
		},
	}
	t.Run("UpdateMapBeyondLink", func(t *testing.T) {
		n, err := traversal.Progress{
			Cfg: &cfg,
		}.FocusedTransform(rootNode, ipld.ParsePath("linkedMap/nested/nonlink"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "linkedMap/nested/nonlink")
			Wish(t, must.String(prev), ShouldEqual, "zoo")
			Wish(t, progress.LastBlock.Path.String(), ShouldEqual, "linkedMap")
			Wish(t, progress.LastBlock.Link.String(), ShouldEqual, "baguqeeyevmbz3ga")
			nb := prev.Prototype().NewBuilder()
			nb.AssignString("new string!")
			return nb.Build(), nil
		}, false)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
		// there should be a new object in our new storage!
		Wish(t, len(store2.Bag), ShouldEqual, 1)
		// cleanup for next test
		store2 = storage.Memory{}
	})
	t.Run("UpdateNotBeyondLink", func(t *testing.T) {
		// This is replacing a link with a non-link.  Doing so shouldn't hit storage.
		n, err := traversal.Progress{
			Cfg: &cfg,
		}.FocusedTransform(rootNode, ipld.ParsePath("linkedMap"), func(progress traversal.Progress, prev ipld.Node) (ipld.Node, error) {
			Wish(t, progress.Path.String(), ShouldEqual, "linkedMap")
			nb := prev.Prototype().NewBuilder()
			nb.AssignString("new string!")
			return nb.Build(), nil
		}, false)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.Kind_Map)
		// there should be no new objects in our new storage!
		Wish(t, len(store2.Bag), ShouldEqual, 0)
		// cleanup for next test
		store2 = storage.Memory{}
	})

	// link traverse to scalar // this is unspecifiable using the current path syntax!  you'll just end up replacing the link with the scalar!
}

func keys(n ipld.Node) []string {
	v := make([]string, 0, n.Length())
	for itr := n.MapIterator(); !itr.Done(); {
		k, _, _ := itr.Next()
		v = append(v, must.String(k))
	}
	return v
}
