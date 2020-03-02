package traversal_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"unicode"

	. "github.com/warpfork/go-wish"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"

	_ "github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/traversal"
)

// Do some fixture fabrication.
// We assume all the builders and serialization must Just Work here.

var storage = make(map[ipld.Link][]byte)
var (
	leafAlpha, leafAlphaLnk         = encode(basicnode.NewString("alpha"))
	leafBeta, leafBetaLnk           = encode(basicnode.NewString("beta"))
	middleMapNode, middleMapNodeLnk = encode(fluent.MustBuildMap(basicnode.Style__Map{}, 3, func(na fluent.MapNodeAssembler) {
		na.AssembleDirectly("foo").AssignBool(true)
		na.AssembleDirectly("bar").AssignBool(false)
		na.AssembleDirectly("nested").CreateMap(2, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly("alink").AssignLink(leafAlphaLnk)
			na.AssembleDirectly("nonlink").AssignString("zoo")
		})
	}))
	middleListNode, middleListNodeLnk = encode(fluent.MustBuildList(basicnode.Style__List{}, 4, func(na fluent.ListNodeAssembler) {
		na.AssembleValue().AssignLink(leafAlphaLnk)
		na.AssembleValue().AssignLink(leafAlphaLnk)
		na.AssembleValue().AssignLink(leafBetaLnk)
		na.AssembleValue().AssignLink(leafAlphaLnk)
	}))
	rootNode, rootNodeLnk = encode(fluent.MustBuildMap(basicnode.Style__Map{}, 4, func(na fluent.MapNodeAssembler) {
		na.AssembleDirectly("plain").AssignString("olde string")
		na.AssembleDirectly("linkedString").AssignLink(leafAlphaLnk)
		na.AssembleDirectly("linkedMap").AssignLink(middleMapNodeLnk)
		na.AssembleDirectly("linkedList").AssignLink(middleListNodeLnk)
	}))
)

// encode hardcodes some encoding choices for ease of use in fixture generation;
// just gimme a link and stuff the bytes in a map.
// (also return the node again for convenient assignment.)
func encode(n ipld.Node) (ipld.Node, ipld.Link) {
	lb := cidlink.LinkBuilder{cid.Prefix{
		Version:  1,
		Codec:    0x0129,
		MhType:   0x17,
		MhLength: 4,
	}}
	lnk, err := lb.Build(context.Background(), ipld.LinkContext{}, n,
		func(ipld.LinkContext) (io.Writer, ipld.StoreCommitter, error) {
			buf := bytes.Buffer{}
			return &buf, func(lnk ipld.Link) error {
				storage[lnk] = buf.Bytes()
				return nil
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}
	return n, lnk
}

// Print a quick little table of our fixtures for sanity check purposes.
func init() {
	withoutWhitespace := func(s string) string {
		return strings.Map(func(r rune) rune {
			if !unicode.IsPrint(r) {
				return -1
			} else {
				return r
			}
		}, s)
	}
	fmt.Printf("fixtures:\n"+strings.Repeat("\t%v\t%v\n", 5),
		leafAlphaLnk, withoutWhitespace(string(storage[leafAlphaLnk])),
		leafBetaLnk, withoutWhitespace(string(storage[leafBetaLnk])),
		middleMapNodeLnk, withoutWhitespace(string(storage[middleMapNodeLnk])),
		middleListNodeLnk, withoutWhitespace(string(storage[middleListNodeLnk])),
		rootNodeLnk, withoutWhitespace(string(storage[rootNodeLnk])),
	)
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

func TestFocusWithLinkLoading(t *testing.T) {
	t.Run("link traversal with no configured loader should fail", func(t *testing.T) {
		t.Run("terminal link should fail", func(t *testing.T) {
			err := traversal.Focus(middleMapNode, ipld.ParsePath("nested/alink"), func(prog traversal.Progress, n ipld.Node) error {
				t.Errorf("should not be reached; no way to load this path")
				return nil
			})
			Wish(t, err.Error(), ShouldEqual, `error traversing node at "nested/alink": could not load link "`+leafAlphaLnk.String()+`": no LinkTargetNodeStyleChooser configured`)
		})
		t.Run("mid-path link should fail", func(t *testing.T) {
			err := traversal.Focus(rootNode, ipld.ParsePath("linkedMap/nested/nonlink"), func(prog traversal.Progress, n ipld.Node) error {
				t.Errorf("should not be reached; no way to load this path")
				return nil
			})
			Wish(t, err.Error(), ShouldEqual, `error traversing node at "linkedMap": could not load link "`+middleMapNodeLnk.String()+`": no LinkTargetNodeStyleChooser configured`)
		})
	})
	t.Run("link traversal with loader should work", func(t *testing.T) {
		err := traversal.Progress{
			Cfg: &traversal.Config{
				LinkLoader: func(lnk ipld.Link, _ ipld.LinkContext) (io.Reader, error) {
					return bytes.NewBuffer(storage[lnk]), nil
				},
				LinkTargetNodeStyleChooser: func(_ ipld.Link, _ ipld.LinkContext) (ipld.NodeStyle, error) {
					return basicnode.Style__Any{}, nil
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
