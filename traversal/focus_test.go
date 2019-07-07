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
	_ "github.com/ipld/go-ipld-prime/encoding/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/traversal"
)

// Do some fixture fabrication.
// We assume all the builders and serialization must Just Work here.

var storage = make(map[ipld.Link][]byte)
var fnb = fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
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

// encode hardcodes some encoding choices for ease of use in fixture generation;
// just gimme a link and stuff the bytes in a map.
// (also return the node again for convenient assignment.)
func encode(n ipld.Node) (ipld.Node, ipld.Link) {
	lb := cidlink.LinkBuilder{Prefix: cid.Prefix{
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
		err := traversal.Focus(fnb.CreateString("x"), ipld.Path{}, func(tp traversal.TraversalProgress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, fnb.CreateString("x"))
			Wish(t, tp.Path.String(), ShouldEqual, ipld.Path{}.String())
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("one step path on map node works", func(t *testing.T) {
		err := traversal.Focus(middleMapNode, ipld.ParsePath("foo"), func(tp traversal.TraversalProgress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, fnb.CreateBool(true))
			Wish(t, tp.Path, ShouldEqual, ipld.ParsePath("foo"))
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("two step path on map node works", func(t *testing.T) {
		err := traversal.Focus(middleMapNode, ipld.ParsePath("nested/nonlink"), func(tp traversal.TraversalProgress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, fnb.CreateString("zoo"))
			Wish(t, tp.Path, ShouldEqual, ipld.ParsePath("nested/nonlink"))
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
}

func TestFocusWithLinkLoading(t *testing.T) {
	t.Run("link traversal with no configured loader should fail", func(t *testing.T) {
		t.Run("terminal link should fail", func(t *testing.T) {
			err := traversal.Focus(middleMapNode, ipld.ParsePath("nested/alink"), func(tp traversal.TraversalProgress, n ipld.Node) error {
				t.Errorf("should not be reached; no way to load this path")
				return nil
			})
			Wish(t, err.Error(), ShouldEqual, `error traversing node at "nested/alink": could not load link "`+leafAlphaLnk.String()+`": no link loader configured`)
		})
		t.Run("mid-path link should fail", func(t *testing.T) {
			err := traversal.Focus(rootNode, ipld.ParsePath("linkedMap/nested/nonlink"), func(tp traversal.TraversalProgress, n ipld.Node) error {
				t.Errorf("should not be reached; no way to load this path")
				return nil
			})
			Wish(t, err.Error(), ShouldEqual, `error traversing node at "linkedMap": could not load link "`+middleMapNodeLnk.String()+`": no link loader configured`)
		})
	})
	t.Run("link traversal with loader should work", func(t *testing.T) {
		err := traversal.TraversalProgress{
			Cfg: &traversal.TraversalConfig{
				LinkLoader: func(lnk ipld.Link, _ ipld.LinkContext) (io.Reader, error) {
					return bytes.NewBuffer(storage[lnk]), nil
				},
			},
		}.Focus(rootNode, ipld.ParsePath("linkedMap/nested/nonlink"), func(tp traversal.TraversalProgress, n ipld.Node) error {
			Wish(t, n, ShouldEqual, fnb.CreateString("zoo"))
			Wish(t, tp.Path, ShouldEqual, ipld.ParsePath("linkedMap/nested/nonlink"))
			Wish(t, tp.LastBlock.Link, ShouldEqual, middleMapNodeLnk)
			Wish(t, tp.LastBlock.Path, ShouldEqual, ipld.ParsePath("linkedMap"))
			return nil
		})
		Wish(t, err, ShouldEqual, nil)
	})
}
